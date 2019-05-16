// Copyright 2019 Stefan Prisca

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tfc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"regexp"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

const CONTRACT_STATE_KEY = "contract.tfc.com"

var ContractID = int32(binary.LittleEndian.Uint16([]byte(CONTRACT_STATE_KEY)))

func HandleInit(APIstub shim.ChaincodeStubInterface) pb.Response {
	// The first argument is the function name!
	// Second will be our protobuf payload.

	gameBoard, err := NewGameBoard()
	if err != nil {
		errStr := fmt.Sprintf("could not create game board: %s", err)
		return shim.Error(errStr)
	}

	identityMap := make(map[int32][]byte)

	var contractUUID = []byte(APIstub.GetTxID())
	identityMap[ContractID] = contractUUID

	gameData := &tfcPb.GameData{
		Board:       gameBoard,
		State:       tfcPb.GameState_JOINING,
		IdentityMap: identityMap,
	}

	protoData, err := proto.Marshal(gameData)
	if err != nil {
		return shim.Error(fmt.Sprintf("could not marshal game data: %s", err))
	}

	APIstub.PutState(CONTRACT_STATE_KEY, protoData)

	return shim.Success(protoData)
}

func HandleInvoke(APIstub shim.ChaincodeStubInterface) pb.Response {

	protoArgs := APIstub.GetArgs()[1]
	trxArgs := &tfcPb.GameContractTrxArgs{}
	err := proto.Unmarshal(protoArgs, trxArgs)
	if err != nil {
		return shim.Error(
			fmt.Sprintf("could not unmarshal arguments proto message <%v>: %s", protoArgs, err))
	}

	gameData, err := getLedgerData(APIstub)
	if err != nil {
		return shim.Error(err.Error())
	}

	creatorSign, err := APIstub.GetCreator()
	if err != nil {
		return shim.Error(fmt.Sprintf(
			"could not retrieve transaction creator: %s", err))
	}

	log.Printf("Handling transaction from state %s", gameData.State)

	// Handle transaction logic
	var newGameData tfcPb.GameData
	switch trxArgs.Type {
	case tfcPb.GameTrxType_JOIN:
		newGameData, err = handleJoin(APIstub, creatorSign, *gameData, *trxArgs.JoinTrxPayload)
	case tfcPb.GameTrxType_ROLL:
		// TODO
		log.Println("ROLL is not yet implemented.")
		newGameData = *gameData
	case tfcPb.GameTrxType_NEXT:
		log.Println("NEXT trx. Nothing to do here")
		newGameData = *gameData
	case tfcPb.GameTrxType_TRADE:
		newGameData, err = handleTrade(APIstub, creatorSign, *gameData, *trxArgs.TradeTrxPayload)
	case tfcPb.GameTrxType_DEV:
		newGameData, err = handleDev(APIstub, creatorSign, *gameData, *trxArgs.BuildTrxPayload)
	default:
		return shim.Error(fmt.Sprint("Unkown transaction type <>"))
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Compute the next game state
	newGameState, err := computeNextState(newGameData, trxArgs.Type)
	if err != nil {
		// Disable during testing until all components are working
		return shim.Error(err.Error())
	}

	newGameData.State = newGameState
	log.Printf("Finished processing transaction, with state %v", newGameData.State)
	// Put the state back on the ledger and return a result
	protoData, err := proto.Marshal(&newGameData)
	if err != nil {
		return shim.Error(fmt.Sprintf("could not marshal game data: %s", err))
	}
	APIstub.PutState(CONTRACT_STATE_KEY, protoData)
	return shim.Success(protoData)

}

func assertJoinPrecond(gameData tfcPb.GameData, payload tfcPb.JoinTrxPayload) error {
	if gameData.State != tfcPb.GameState_JOINING {
		return fmt.Errorf("unexpected game state. expected %v, got %v",
			tfcPb.GameState_JOINING, gameData.State)
	}

	playerID := GetPlayerId(payload.Player)
	if _, ok := gameData.Profiles[playerID]; ok {
		return fmt.Errorf("player <%v> already taken", payload.Player)
	}
	return nil
}

func handleJoin(APIstub shim.ChaincodeStubInterface, creatorSign []byte,
	gameData tfcPb.GameData, payload tfcPb.JoinTrxPayload) (tfcPb.GameData, error) {

	err := assertJoinPrecond(gameData, payload)
	if err != nil {
		return gameData, fmt.Errorf(
			"join preconditions not met: %s", err)
	}

	log.Printf("Joining player %v with sign %v", payload.Player, creatorSign)
	playerID := GetPlayerId(payload.Player)
	gameData.IdentityMap[playerID] = creatorSign

	// If there was no other player that joined until now
	// the profiles will be nil
	if gameData.Profiles == nil {
		gameData.Profiles = make(map[int32]*tfcPb.PlayerProfile)
	}
	gameData.Profiles[playerID] = InitPlayerProfile()
	return gameData, nil
}

func handleTrade(APIstub shim.ChaincodeStubInterface, creatorSign []byte,
	gameData tfcPb.GameData, payload tfcPb.TradeTrxPayload) (tfcPb.GameData, error) {

	err := assertTradePrecond(gameData, creatorSign, payload)
	if err != nil {
		return gameData, fmt.Errorf(
			"trade preconditions not met: %s", err)
	}

	srcID := GetPlayerId(payload.Source)
	destID := GetPlayerId(payload.Dest)
	resID := GetResourceId(payload.Resource)

	srcProfile := gameData.Profiles[srcID]
	destProfile := gameData.Profiles[destID]

	srcProfile.Resources[resID] -= payload.Amount
	destProfile.Resources[resID] += payload.Amount

	return gameData, assertTradePostcond(gameData, payload)
}

func stateValidationString(src tfcPb.Player, state tfcPb.GameState) string {
	return fmt.Sprintf("%v%v", src, state)
}

var tradeStateValidationRegexp = regexp.MustCompile(
	fmt.Sprintf("%vRTRADE", tfcPb.Player_RED) +
		fmt.Sprintf("|%vBTRADE", tfcPb.Player_BLUE) +
		fmt.Sprintf("|%vGTRADE", tfcPb.Player_GREEN))

// TODO: implement this
func assertTradePrecond(gameData tfcPb.GameData, creatorSign []byte, payload tfcPb.TradeTrxPayload) error {
	state := gameData.State
	creator, err := getCreator(gameData, creatorSign)
	if err != nil || creator != payload.Source {
		return fmt.Errorf("invalid trx creator, or creator not identified (<0): expected %v, got %v",
			creator, payload.Source)
	}

	stateValidationStr := stateValidationString(creator, state)
	if !tradeStateValidationRegexp.MatchString(stateValidationStr) {
		return fmt.Errorf("expected state to match one of %v, got %v",
			tradeStateValidationRegexp, stateValidationStr)
	}

	return nil
}

func assertTradePostcond(gameData tfcPb.GameData, payload tfcPb.TradeTrxPayload) error {
	res := payload.Resource
	if err := hasValidPostTradeAmount(payload.Source, gameData, res); err != nil {
		return err
	}

	if err := hasValidPostTradeAmount(payload.Dest, gameData, res); err != nil {
		return err
	}
	return nil
}

func hasValidPostTradeAmount(p tfcPb.Player, gameData tfcPb.GameData, r tfcPb.Resource) error {
	rID := GetResourceId(r)
	pP := *gameData.Profiles[GetPlayerId(p)]
	available := pP.Resources[rID]
	if available < 0 {
		return fmt.Errorf("player %v does not have required %v resources: %s",
			p, r,
			fmt.Sprintf("available: %v", available))
	}
	return nil
}

func computeNextState(gameData tfcPb.GameData, txType tfcPb.GameTrxType) (tfcPb.GameState, error) {
	st := gameData.State
	switch {
	// A player just joined, move to RROLL if all are in
	case txType == tfcPb.GameTrxType_JOIN:
		if len(gameData.Profiles) == 3 {
			return tfcPb.GameState_RROLL, nil
		}
		return tfcPb.GameState_JOINING, nil

	case txType == tfcPb.GameTrxType_ROLL:
		switch st {
		case tfcPb.GameState_RROLL:
			return tfcPb.GameState_RTRADE, nil
		case tfcPb.GameState_BROLL:
			return tfcPb.GameState_BTRADE, nil
		case tfcPb.GameState_GROLL:
			return tfcPb.GameState_GTRADE, nil
		}

	case txType == tfcPb.GameTrxType_NEXT:
		switch st {
		case tfcPb.GameState_RROLL:
			return tfcPb.GameState_RTRADE, nil
		case tfcPb.GameState_RTRADE:
			return tfcPb.GameState_RDEV, nil
		case tfcPb.GameState_RDEV:
			if won(tfcPb.Player_RED, gameData) {
				return tfcPb.GameState_RWON, nil
			}
			return tfcPb.GameState_BROLL, nil
		case tfcPb.GameState_BTRADE:
			return tfcPb.GameState_BDEV, nil
		case tfcPb.GameState_BDEV:
			if won(tfcPb.Player_BLUE, gameData) {
				return tfcPb.GameState_BWON, nil
			}
			return tfcPb.GameState_GROLL, nil
		case tfcPb.GameState_GTRADE:
			return tfcPb.GameState_GDEV, nil
		case tfcPb.GameState_GDEV:
			if won(tfcPb.Player_GREEN, gameData) {
				return tfcPb.GameState_GWON, nil
			}
			return tfcPb.GameState_RROLL, nil
		}
	case txType == tfcPb.GameTrxType_TRADE:
		return st, nil
	case txType == tfcPb.GameTrxType_DEV:
		return st, nil
	case txType == tfcPb.GameTrxType_BATTLE:
		return st, nil
	}
	return st, fmt.Errorf(
		"could not compute next state from st %v and trx type %v", st, txType)
}

func won(player tfcPb.Player, gameData tfcPb.GameData) bool {
	id := GetPlayerId(player)
	profile := gameData.Profiles[id]
	return profile.WinningPoints > 10
}

func getLedgerData(APIstub shim.ChaincodeStubInterface) (*tfcPb.GameData, error) {
	protoData, err := APIstub.GetState(CONTRACT_STATE_KEY)
	if err != nil {
		return nil, fmt.Errorf("Could not get the contract from state. Error: %s", err.Error())
	}

	gameData := &tfcPb.GameData{}
	err = proto.Unmarshal(protoData, gameData)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the proto contract. Error: %s", err.Error())
	}
	return gameData, nil
}

var playerExists = regexp.MustCompile(fmt.Sprintf("%v|%v|%v",
	tfcPb.Player_BLUE, tfcPb.Player_GREEN, tfcPb.Player_RED))

func getCreator(gameData tfcPb.GameData, creatorSign []byte) (tfcPb.Player, error) {
	srcID := int32(-1)
	for pID, sign := range gameData.IdentityMap {
		if bytes.Equal(sign, creatorSign) {
			srcID = pID
			break
		}
	}
	src := tfcPb.Player(srcID)

	if !playerExists.MatchString(src.String()) {
		return src, fmt.Errorf("unkown creator signature: %v", creatorSign)
	}

	return src, nil
}
