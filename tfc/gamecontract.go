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
	case tfcPb.GameTrxType_TRADE:
		newGameData, err = handleTrade(APIstub, creatorSign, *gameData, *trxArgs.TradeTrxPayload)
	case tfcPb.GameTrxType_DEV:
		newGameData, err = handleDev(APIstub, *gameData, *trxArgs.BuildTrxPayload)
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

func handleDev(APIstub shim.ChaincodeStubInterface, creatorSign []byte,
	gameData tfcPb.GameData, payload tfcPb.BuildTrxPayload) (tfcPb.GameData, error) {

	err := assertDevelopmentPrecond(gameData, payload)
	if err != nil {
		return gameData, fmt.Errorf(
			"development preconditions not met: %s", err)
	}

	switch payload.Type {
	case tfcPb.BuildType_SETTLE:
		return buildSettlement(gameData, *payload.BuildSettlePayload)
	case tfcPb.BuildType_ROAD:
		return buildRoad(gameData, *payload.BuildRoadPayload)
	}

	return gameData, nil
}

var devStateValidationRegexp = regexp.MustCompile(
	fmt.Sprintf("%vRDEV", tfcPb.Player_RED) +
		fmt.Sprintf("|%vBDEV", tfcPb.Player_BLUE) +
		fmt.Sprintf("|%vGDEV", tfcPb.Player_GREEN))

// TODO: implement this
func assertDevelopmentPrecond(gameData tfcPb.GameData, creatorSign []byte, payload tfcPb.BuildTrxPayload) error {

	/*
		1) correct state
		2.e) empty and an edge point is one of the player settlements
		2.s) empty and neighbouring intersections for distance rule: all surrounding intersections are free
	*/

	state := gameData.State
	creator, err := getCreator(gameData, creatorSign)
	if err != nil {
		return err
	}

	stateValidationStr := stateValidationString(creator, state)
	if !devStateValidationRegexp.MatchString(stateValidationStr) {
		return fmt.Errorf("expected state to match one of %v, got %v",
			devStateValidationRegexp, stateValidationStr)
	}

	switch payload.Type {
	case tfcPb.BuildType_ROAD:
		return assertBuildRoadPrecond(gameData, creator, *payload.BuildRoadPayload)
	case tfcPb.BuildType_SETTLE:
		return assertBuildSettlePrecond(gameData, creator, *payload.BuildSettlePayload)
	}

	return fmt.Errorf("Unkown build type %v", payload.Type)
}

var canBuildRoad = regexp.MustCompile(
	fmt.Sprintf("(%s)|(%s)|(%s)",
		fmt.Sprintf("%v%v.*(%v|%v)+.*",
			tfcPb.Player_RED, tfcPb.Road_NOROAD,
			tfcPb.Road_REDROAD, tfcPb.Settlement_REDSETTLE),
		fmt.Sprintf("%v%v.*(%v|%v)+.*",
			tfcPb.Player_GREEN, tfcPb.Road_NOROAD,
			tfcPb.Road_GREENROAD, tfcPb.Settlement_GREENSETTLE),
		fmt.Sprintf("%v%v.*(%v|%v)+.*",
			tfcPb.Player_BLUE, tfcPb.Road_NOROAD,
			tfcPb.Road_BLUEROAD, tfcPb.Settlement_BLUESETTLE)))

func buildRoadValidString(p tfcPb.Player, r tfcPb.Road, s1, s2 tfcPb.Settlement, r1, r2, r3, r4 tfcPb.Road) string {
	return fmt.Sprintf("%v%v%v%v%v%v%v%v", p, r, s1, s2, r1, r2, r3, r4)
}

func assertBuildRoadPrecond(gameData tfcPb.GameData, creator tfcPb.Player, payload tfcPb.BuildRoadPayload) error {
	if creator != payload.Player {
		return fmt.Errorf("expected creator to match trx player. expected %v, got %v", creator, payload.Player)
	}

	eID := uint32(payload.EdgeID)
	E, exists := gameData.Board.Edges[eID]
	if !exists {
		return fmt.Errorf("gameboard edge %v does not exist", eID)
	}

	r := E.Attributes.Road
	nE := gameData.Board.Edges[E.Next]

	s1 := gameData.Board.Intersections[E.Origin].Attributes.Settlement
	s2 := gameData.Board.Intersections[nE.Origin].Attributes.Settlement

	r1 := nE.Attributes.Road
	r2 := gameData.Board.Edges[E.Prev].Attributes.Road

	twin := gameData.Board.Edges[E.Twin]
	r3 := gameData.Board.Edges[twin.Prev].Attributes.Road
	r4 := gameData.Board.Edges[twin.Next].Attributes.Road

	validationS := buildRoadValidString(creator, r, s1, s2, r1, r2, r3, r4)

	if !canBuildRoad.MatchString(validationS) {
		return fmt.Errorf("could not build road for player %v, conditions not fulfilled: %s", creator,
			fmt.Sprintf("existing road %v; surrounding settlements [%v, %v]; surrounding edges [%v, %v, %v, %v]",
				r, s1, s2, r1, r2, r3, r4))
	}

	return nil
}

var canBuildSettle = regexp.MustCompile(
	fmt.Sprintf("%v%v%v%v",
		tfcPb.Settlement_NOSETTLE, tfcPb.Settlement_NOSETTLE,
		tfcPb.Settlement_NOSETTLE, tfcPb.Settlement_NOSETTLE))

func buildSettleValidString(s, s1, s2, s3 tfcPb.Settlement) string {
	return fmt.Sprintf("%v%v%v%v", s, s1, s2, s3)
}

func assertBuildSettlePrecond(gameData tfcPb.GameData, creator tfcPb.Player, payload tfcPb.BuildSettlePayload) error {
	if creator != payload.Player {
		return fmt.Errorf("expected creator to match trx player. expected %v, got %v", creator, payload.Player)
	}

	sID := uint32(payload.SettleID)
	I := gameData.Board.Intersections[sID]
	s := I.Attributes.Settlement

	iE := gameData.Board.Edges[I.IncidentEdge]
	iID1 := gameData.Board.Edges[iE.Next].Origin
	s1 := gameData.Board.Intersections[iID1].Attributes.Settlement

	iID2 := gameData.Board.Edges[iE.Prev].Origin
	s2 := gameData.Board.Intersections[iID2].Attributes.Settlement

	iETwin := gameData.Board.Edges[iE.Twin]
	iETwinNext := gameData.Board.Edges[iETwin.Next]
	iID3 := gameData.Board.Edges[iETwinNext.Next].Origin
	s3 := gameData.Board.Intersections[iID3].Attributes.Settlement

	validationS := buildSettleValidString(s, s1, s2, s3)

	if !canBuildRoad.MatchString(validationS) {
		return fmt.Errorf("could not build road for player %v, conditions not fulfilled: %s", creator,
			fmt.Sprintf("existing settle %v; surrounding settlements [%v, %v, %v, %v]", s, s1, s2, s3))
	}

	return nil
}

func buildSettlement(
	gameData tfcPb.GameData, payload tfcPb.BuildSettlePayload) (tfcPb.GameData, error) {

	playerID := GetPlayerId(payload.Player)
	profile := gameData.Profiles[playerID]

	for rID := range profile.Resources {
		profile.Resources[rID]--
	}

	profile.Settlements--
	profile.WinningPoints += 2

	posID := uint32(payload.SettleID)
	settleIntersection := gameData.Board.Intersections[posID]
	switch payload.Player {
	case tfcPb.Player_RED:
		settleIntersection.Attributes.Settlement = tfcPb.Settlement_REDSETTLE
	case tfcPb.Player_GREEN:
		settleIntersection.Attributes.Settlement = tfcPb.Settlement_GREENSETTLE
	case tfcPb.Player_BLUE:
		settleIntersection.Attributes.Settlement = tfcPb.Settlement_BLUESETTLE
	}

	return gameData, nil
}

func buildRoad(
	gameData tfcPb.GameData, payload tfcPb.BuildRoadPayload) (tfcPb.GameData, error) {

	playerID := GetPlayerId(payload.Player)
	profile := gameData.Profiles[playerID]

	profile.Resources[GetResourceId(tfcPb.Resource_HILL)]--
	profile.Resources[GetResourceId(tfcPb.Resource_FOREST)]--
	profile.Roads--
	profile.WinningPoints++

	eID := uint32(payload.EdgeID)
	edge := gameData.Board.Edges[eID]
	switch payload.Player {
	case tfcPb.Player_RED:
		edge.Attributes.Road = tfcPb.Road_REDROAD
	case tfcPb.Player_GREEN:
		edge.Attributes.Road = tfcPb.Road_GREENROAD
	case tfcPb.Player_BLUE:
		edge.Attributes.Road = tfcPb.Road_BLUEROAD
	}

	return gameData, nil
}

func computeNextState(gameData tfcPb.GameData, txType tfcPb.GameTrxType) (tfcPb.GameState, error) {
	st := gameData.State
	switch {
	// A player just joined, move to RROLL if all are in
	case txType == tfcPb.GameTrxType_JOIN:
		log.Printf("Computing st after Join. %v", gameData.Profiles)
		if len(gameData.Profiles) == 3 {
			log.Printf("Moved to %v", tfcPb.GameState_RROLL)
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
