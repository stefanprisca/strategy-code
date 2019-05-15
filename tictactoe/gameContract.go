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

package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	tttPb "github.com/stefanprisca/strategy-protobufs/tictactoe"
)

// Dummy struct for hyperledger
type GameContract struct {
}

const CONTRACT_STATE_KEY = "contract.tictactoe"
const BOARD_SIZE = 9

func (gc *GameContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	positions := make([]tttPb.Mark, BOARD_SIZE)
	for i := 0; i < BOARD_SIZE; i++ {
		positions[i] = tttPb.Mark_E
	}

	tttContract := &tttPb.TttContract{
		State:     tttPb.TttContract_XTURN,
		Positions: positions,
		XPlayer:   "player1",
		OPlayer:   "player2",
	}

	tttState, err := proto.Marshal(tttContract)
	if err != nil {
		errMsg := fmt.Sprintf("Could not marshal the contract. Error: %s", err.Error())
		return shim.Error(errMsg)
	}
	APIstub.PutState(CONTRACT_STATE_KEY, tttState)
	return shim.Success(tttState)
}

func (gc *GameContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	// creator, errc := APIstub.GetCreator()
	// if errc == nil {
	// 	fmt.Println("Creator: ", string(creator))
	// }

	// The first argument is the function name!
	// Second will be our protobuf payload.

	invokeST := time.Now()
	defer func() {
		invokeDuration := time.Since(invokeST).Seconds()
		log.Printf("#############\n\t FINISHED INVOKE FUNCTION IN < %v > seconds", invokeDuration)
	}()

	protoTrxArgs := APIstub.GetArgs()[1]

	trxArgs := &tttPb.TrxArgs{}
	err := proto.Unmarshal(protoTrxArgs, trxArgs)
	if err != nil {
		errMsg := fmt.Sprintf("Could not parse transaction args %v. Error %s", trxArgs, err.Error())
		return shim.Error(errMsg)
	}

	switch trxArgs.Type {
	case tttPb.TrxType_MOVE:
		return move(APIstub, trxArgs.MovePayload)
	}

	return shim.Error(fmt.Sprintf("Unkown transaction type < %v >", trxArgs.Type))
}

func move(APIstub shim.ChaincodeStubInterface, payload *tttPb.MoveTrxPayload) pb.Response {
	if payload == nil {
		return shim.Error("Unexpected empty payload. Failed to do move.")
	}

	contract, err := getLedgerContract(APIstub)
	if err != nil {
		return shim.Error(err.Error())
	}

	if contractTerminated(*contract) {
		return shim.Error(fmt.Sprintf("Contract already terminated with state %v", contract.State))
	}

	if err = validateMoveArgs(APIstub, *contract, *payload); err != nil {
		return shim.Error(err.Error())
	}

	newContract, err := applyMove(*contract, *payload)
	if err != nil {
		return shim.Error(err.Error())
	}
	newProtoContract, err := proto.Marshal(&newContract)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = APIstub.PutState(CONTRACT_STATE_KEY, newProtoContract)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(newProtoContract)
}

var terminatedRegexp = regexp.MustCompile(
	fmt.Sprintf("^%v|%v|%v$",
		tttPb.TttContract_XWON, tttPb.TttContract_OWON, tttPb.TttContract_TIE))

func contractTerminated(contract tttPb.TttContract) bool {
	return terminatedRegexp.MatchString(
		fmt.Sprintf("%v", contract.State))
}

var positionRegexp = regexp.MustCompile(fmt.Sprintf("^[0-8]%v$", tttPb.Mark_E))

func positionValidationString(position int32, mark tttPb.Mark) string {
	return fmt.Sprintf("%d%v", position, mark)
}

var turnRegexp = regexp.MustCompile(
	fmt.Sprintf("^(%s|%s)$",
		turnValidationString(tttPb.Mark_X, tttPb.TttContract_XTURN),
		turnValidationString(tttPb.Mark_O, tttPb.TttContract_OTURN)))

func turnValidationString(mark tttPb.Mark, state tttPb.TttContract_State) string {
	return fmt.Sprintf("%v%v", mark, state)
}

func validateMoveArgs(APIstub shim.ChaincodeStubInterface, contract tttPb.TttContract, payload tttPb.MoveTrxPayload) error {
	pvs := positionValidationString(payload.Position,
		contract.Positions[payload.Position])
	if !positionRegexp.MatchString(pvs) {
		return fmt.Errorf("Invalid position or position not empty. Position < %d >, mark < %v >",
			payload.Position, contract.Positions[payload.Position])
	}

	tvs := turnValidationString(payload.Mark, contract.State)
	if !turnRegexp.MatchString(tvs) {
		return fmt.Errorf("Invalid turn. Got mark < %v >, expected < %v >", payload.Mark, contract.State)
	}

	return nil
}

func getLedgerContract(APIstub shim.ChaincodeStubInterface) (*tttPb.TttContract, error) {
	contractBytes, err := APIstub.GetState(CONTRACT_STATE_KEY)
	if err != nil {
		return nil, fmt.Errorf("Could not get the contract from state. Error: %s", err.Error())
	}

	actualContract := &tttPb.TttContract{}
	err = proto.Unmarshal(contractBytes, actualContract)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the proto contract. Error: %s", err.Error())
	}
	return actualContract, nil
}

func applyMove(contract tttPb.TttContract, payload tttPb.MoveTrxPayload) (tttPb.TttContract, error) {
	newPositions := contract.Positions
	newPositions[payload.Position] = payload.Mark

	nextState, err := computeNextState(newPositions, contract.State)
	if err != nil {
		return tttPb.TttContract{}, err
	}

	return tttPb.TttContract{
		Positions: newPositions,
		State:     nextState,
		XPlayer:   contract.XPlayer,
		OPlayer:   contract.OPlayer,
	}, nil
}

func computeNextState(positions []tttPb.Mark, state tttPb.TttContract_State) (tttPb.TttContract_State, error) {
	posString, err := winValidationString(positions)
	if err != nil {
		return state, err
	}

	switch state {
	case tttPb.TttContract_XTURN:
		if won(tttPb.Mark_X, posString) {
			return tttPb.TttContract_XWON, nil
		} else if boardFull(posString) {
			return tttPb.TttContract_TIE, nil
		} else {
			return tttPb.TttContract_OTURN, nil
		}
	case tttPb.TttContract_OTURN:
		if won(tttPb.Mark_O, posString) {
			return tttPb.TttContract_OWON, nil
		} else if boardFull(posString) {
			return tttPb.TttContract_TIE, nil
		} else {
			return tttPb.TttContract_XTURN, nil
		}
	default:
		return state, fmt.Errorf("Could not determine next state")
	}
}

func winValidationString(positions []tttPb.Mark) (string, error) {
	if len(positions) != BOARD_SIZE {
		return "", fmt.Errorf(
			"Invalid number of positions detected. Expected %d, got %d",
			BOARD_SIZE, len(positions))
	}

	result := ""
	for i := 0; i < BOARD_SIZE; i++ {
		result += positions[i].String()
	}
	return result, nil
}

func won(m tttPb.Mark, positions string) bool {
	return threeInRowRegex(m).MatchString(positions) ||
		threeInColRegex(m).MatchString(positions) ||
		threeInDiagRegex(m).MatchString(positions)
}

func threeInRowRegex(m tttPb.Mark) *regexp.Regexp {
	return regexp.MustCompile(
		fmt.Sprintf("(^|...)(%s%s%s)(...|$)",
			m.String(), m.String(), m.String()))
}

func threeInColRegex(m tttPb.Mark) *regexp.Regexp {
	return regexp.MustCompile(
		fmt.Sprintf("^(((%s..){3})|((.%s.){3})|((..%s){3}))$",
			m.String(), m.String(), m.String()))
}

func threeInDiagRegex(m tttPb.Mark) *regexp.Regexp {
	return regexp.MustCompile(
		fmt.Sprintf("^(%s...%s...%s|..%s.%s.%s..)$",
			m.String(), m.String(), m.String(),
			m.String(), m.String(), m.String()))
}

var anyEmptyRegex = regexp.MustCompile(fmt.Sprintf("^.*%s.*$", tttPb.Mark_E.String()))

func boardFull(positions string) bool {
	return !anyEmptyRegex.MatchString(positions)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(GameContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
