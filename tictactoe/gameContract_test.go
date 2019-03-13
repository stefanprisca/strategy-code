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
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	tttPb "github.com/stefanprisca/strategy-protobufs/tictactoe"
)

func initContract(t *testing.T) *shim.MockStub {
	stub := shim.NewMockStub("mockGameContract", new(GameContract))
	if stub == nil {
		t.Fatalf("Failed to init mock")
	}
	r := stub.MockInit("001", [][]byte{})

	if r.GetStatus() != shim.OK {
		t.Fatalf("Could not init the contract. Error: %s", r.Message)
	}
	return stub
}

func TestInit(t *testing.T) {
	stub := initContract(t)

	expectedMarks := make([]tttPb.Mark, 9)
	for i := range expectedMarks {
		expectedMarks[i] = tttPb.Mark_E
	}

	expectedContract := tttPb.TttContract{
		Positions: expectedMarks,
		State:     tttPb.TttContract_XTURN,
		XPlayer:   "player1",
		OPlayer:   "player2",
	}

	actualContract, err := getLedgerContract(stub)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ok, msg := contractsEqual(expectedContract, *actualContract); !ok {
		t.Fatalf("Initialized contract is different from expected contract. error: %s", msg)
	}
}

func contractsEqual(c1, c2 tttPb.TttContract) (bool, string) {

	if c1.State != c2.State {
		return false, fmt.Sprintf("Contract state does not match. C1: < %v >, C2: < %v >",
			c1.State, c2.State)
	}

	if len(c1.Positions) != len(c2.Positions) {
		return false, fmt.Sprintf("Contract position lengths do not match. C1: < %v >, C2: < %v >",
			len(c1.Positions), len(c2.Positions))
	}

	for i := range c1.Positions {
		if c1.Positions[i] != c2.Positions[i] {
			return false, fmt.Sprintf("Contract position mismatched at %d. C1: < %v >, C2: < %v >",
				i, c1.Positions, c2.Positions)
		}
	}

	return true, ""
}

func TestInvokeMove(t *testing.T) {

	stub := initContract(t)
	positionID := int32(1)
	mark := tttPb.Mark_X

	r, err := newMoveArgsBuilder("001", positionID, mark).
		invoke(stub)
	if err != nil {
		t.Fatal(err.Error())
	}
	if r.GetStatus() != shim.OK {
		t.Logf("Could not invoke move function, error: %s", r.GetMessage())
		t.FailNow()
	}

	actualContract, err := getLedgerContract(stub)
	if err != nil {
		t.Fatalf(err.Error())
	}
	actualMark := actualContract.Positions[positionID]
	if actualMark != mark {
		t.Fatalf("Failed to mark position <%d>. Expected mark <%v>. Actual <%v>",
			positionID, mark, actualMark)
	}
}

func TestInvokeMoveOnOccupiedPos(t *testing.T) {
	stub := initContract(t)
	positionID := int32(1)

	markX := tttPb.Mark_X

	_, err := newMoveArgsBuilder("001", positionID, markX).
		invoke(stub)
	if err != nil {
		t.Fatal(err.Error())
	}

	r, err := newMoveArgsBuilder("002", positionID, tttPb.Mark_O).
		invoke(stub)
	if err != nil {
		t.Fatal(err.Error())
	}

	if r.GetStatus() == shim.OK {
		t.Fatalf("Did not expect invocation to be successful! position <%d> was taken!", positionID)
	}

	actualContract, err := getLedgerContract(stub)
	if err != nil {
		t.Fatalf(err.Error())
	}
	actualMark := actualContract.Positions[positionID]
	if actualMark != markX {
		t.Fatalf("Unexpected mark on position <%d>. Expected mark <%v>. Actual <%v>",
			positionID, markX, actualMark)
	}
}

func TestInvokeMoveWrongMark(t *testing.T) {
	stub := initContract(t)
	position := int32(1)

	actualContract, err := getLedgerContract(stub)
	if err != nil {
		t.Fatalf(err.Error())
	}

	mark := tttPb.Mark_O
	if actualContract.State != tttPb.TttContract_XTURN {
		mark = tttPb.Mark_X
	}

	r, err := newMoveArgsBuilder("001", position, mark).
		invoke(stub)
	if err != nil {
		t.Fatal(err.Error())
	}
	if r.GetStatus() == shim.OK {
		t.Fatalf("Did not expect invocation to be successful. Turn was < %v >  and mark < %v >", actualContract.State, mark)
	}

	actualContract, err = getLedgerContract(stub)
	if err != nil {
		t.Fatalf(err.Error())
	}
	actualMark := actualContract.Positions[position]
	if actualMark != tttPb.Mark_E {
		t.Fatalf("Unexpected mark on position <%d>. Expected mark <%v>. Actual <%v>",
			position, tttPb.Mark_E, actualMark)
	}
}

func TestWinRegexpFunctions(t *testing.T) {

	positions := []string{
		tttPb.Mark_X.String(), tttPb.Mark_O.String(), tttPb.Mark_X.String(),
		tttPb.Mark_X.String(), tttPb.Mark_X.String(), tttPb.Mark_X.String(),
		tttPb.Mark_O.String(), tttPb.Mark_O.String(), tttPb.Mark_X.String()}
	posString := strings.Join(positions, "")

	if !matchAll(tttPb.Mark_X, posString) {
		t.Fatalf("Expected X to match on all dimensions.")
	}

	if matchAny(tttPb.Mark_O, posString) {
		t.Fatalf("Expected O to not match any dimensions.")
	}
}

func matchAll(m tttPb.Mark, posString string) bool {
	tr := threeInRowRegex(m).MatchString(posString)
	tc := threeInColRegex(m).MatchString(posString)
	td := threeInDiagRegex(m).MatchString(posString)

	return tr && tc && td
}

func matchAny(m tttPb.Mark, posString string) bool {
	tr := threeInRowRegex(m).MatchString(posString)
	tc := threeInColRegex(m).MatchString(posString)
	td := threeInDiagRegex(m).MatchString(posString)

	return tr || tc || td
}

func TestGameScriptWinX(t *testing.T) {
	script := []tttPb.MoveTrxPayload{
		{Position: 0, Mark: tttPb.Mark_X},
		{Position: 1, Mark: tttPb.Mark_O},
		{Position: 4, Mark: tttPb.Mark_X},
		{Position: 8, Mark: tttPb.Mark_O},
		{Position: 3, Mark: tttPb.Mark_X},
		{Position: 5, Mark: tttPb.Mark_O},
		{Position: 6, Mark: tttPb.Mark_X},
	}
	stub := initContract(t)
	_, err := runScriptAndCheckLastState(script, tttPb.TttContract_XWON, stub)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestGameScriptWinO(t *testing.T) {
	script := []tttPb.MoveTrxPayload{
		{Position: 0, Mark: tttPb.Mark_X},
		{Position: 1, Mark: tttPb.Mark_O},
		{Position: 4, Mark: tttPb.Mark_X},
		{Position: 8, Mark: tttPb.Mark_O},
		{Position: 3, Mark: tttPb.Mark_X},
		{Position: 5, Mark: tttPb.Mark_O},
		{Position: 7, Mark: tttPb.Mark_X},
		{Position: 2, Mark: tttPb.Mark_O},
	}
	stub := initContract(t)
	_, err := runScriptAndCheckLastState(script, tttPb.TttContract_OWON, stub)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestGameScriptTie(t *testing.T) {
	script := []tttPb.MoveTrxPayload{
		{Position: 0, Mark: tttPb.Mark_X},
		{Position: 1, Mark: tttPb.Mark_O},
		{Position: 4, Mark: tttPb.Mark_X},
		{Position: 8, Mark: tttPb.Mark_O},
		{Position: 3, Mark: tttPb.Mark_X},
		{Position: 5, Mark: tttPb.Mark_O},
		{Position: 7, Mark: tttPb.Mark_X},
		{Position: 6, Mark: tttPb.Mark_O},
		{Position: 2, Mark: tttPb.Mark_X},
	}
	stub := initContract(t)
	_, err := runScriptAndCheckLastState(script, tttPb.TttContract_TIE, stub)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestGameScriptInvalidMoveXTurn(t *testing.T) {
	script := []tttPb.MoveTrxPayload{
		{Position: 0, Mark: tttPb.Mark_X},
		{Position: 1, Mark: tttPb.Mark_O},
		{Position: 8, Mark: tttPb.Mark_O},
	}
	stub := initContract(t)
	responses, err := runScriptAndCheckLastState(script, tttPb.TttContract_XTURN, stub)
	if err == nil {
		t.Fatal("Expected script to fail")
	}

	lastResponse := responses[len(responses)-1]
	if lastResponse.GetStatus() != shim.ERROR {
		t.Fatal("Expected script to fail")
	}
	t.Log(lastResponse.Message)
}

func TestGameScriptInvalidMoveXWon(t *testing.T) {
	script := []tttPb.MoveTrxPayload{
		{Position: 0, Mark: tttPb.Mark_X},
		{Position: 1, Mark: tttPb.Mark_O},
		{Position: 4, Mark: tttPb.Mark_X},
		{Position: 8, Mark: tttPb.Mark_O},
		{Position: 3, Mark: tttPb.Mark_X},
		{Position: 5, Mark: tttPb.Mark_O},
		{Position: 6, Mark: tttPb.Mark_X},
		{Position: 2, Mark: tttPb.Mark_O},
	}
	stub := initContract(t)
	responses, err := runScriptAndCheckLastState(script, tttPb.TttContract_XWON, stub)
	if err == nil {
		t.Fatal("Expected script to fail")
	}

	lastResponse := responses[len(responses)-1]
	if lastResponse.GetStatus() != shim.ERROR {
		t.Fatal("Expected script to fail")
	}
	t.Log(lastResponse.Message)
}

func TestGameScriptInvalidMoveTie(t *testing.T) {
	script := []tttPb.MoveTrxPayload{
		{Position: 0, Mark: tttPb.Mark_X},
		{Position: 1, Mark: tttPb.Mark_O},
		{Position: 4, Mark: tttPb.Mark_X},
		{Position: 8, Mark: tttPb.Mark_O},
		{Position: 3, Mark: tttPb.Mark_X},
		{Position: 5, Mark: tttPb.Mark_O},
		{Position: 7, Mark: tttPb.Mark_X},
		{Position: 6, Mark: tttPb.Mark_O},
		{Position: 2, Mark: tttPb.Mark_X},
		{Position: 2, Mark: tttPb.Mark_X},
	}
	stub := initContract(t)
	responses, err := runScriptAndCheckLastState(script, tttPb.TttContract_XWON, stub)
	if err == nil {
		t.Fatal("Expected script to fail")
	}

	lastResponse := responses[len(responses)-1]
	if lastResponse.GetStatus() != shim.ERROR {
		t.Fatal("Expected script to fail")
	}
	t.Log(lastResponse.Message)
}

func runScriptAndCheckLastState(script []tttPb.MoveTrxPayload, expectedState tttPb.TttContract_State, stub *shim.MockStub) ([]pb.Response, error) {
	responses, err := runScript(script, stub)
	if err != nil {
		return responses, fmt.Errorf("Could not run script. Error %s", err.Error())
	}

	for _, r := range responses {
		if r.GetStatus() != shim.OK {
			return responses, fmt.Errorf("Got unexpected response %s", r.Message)
		}
	}

	lastContractPayload := responses[len(responses)-1].GetPayload()
	lastContract := &tttPb.TttContract{}
	if err = proto.Unmarshal(lastContractPayload, lastContract); err != nil {
		return responses, fmt.Errorf("Could not unmarshal last contract. Error %s", err.Error())
	}

	if lastContract.State != expectedState {
		return responses, fmt.Errorf("Expected last state to be %v. Actual %v \n Board:%v",
			expectedState, lastContract.State, lastContract.Positions)
	}

	return responses, nil
}

func TestContractRandomlyTerminates(t *testing.T) {
	script := generateRandomScript()
	stub := initContract(t)
	t.Logf("Running script %v \n", script)
	runScript(script, stub)

	contract, err := getLedgerContract(stub)
	if err != nil {
		t.Fatal(err.Error())
	}

	if contract.State == tttPb.TttContract_XTURN ||
		contract.State == tttPb.TttContract_OTURN {
		t.Fatalf("Unexpected state for the contract %v. Expected terminal state (one of [%v %v %v]) ",
			contract.State, tttPb.TttContract_TIE, tttPb.TttContract_XWON, tttPb.TttContract_OWON)
	}
}

func generateRandomScript() []tttPb.MoveTrxPayload {
	rand.Seed(time.Now().UnixNano())
	script := make([]tttPb.MoveTrxPayload, BOARD_SIZE)
	positions := make([]int, BOARD_SIZE)
	for i := 0; i < BOARD_SIZE; i++ {
		positions[i] = i
	}

	nm := tttPb.Mark_X
	for i := 0; i < BOARD_SIZE; i++ {

		np := rand.Intn(len(positions))
		script[i] = tttPb.MoveTrxPayload{Position: int32(positions[np]), Mark: nm}

		pNext := positions[:np]
		for k := np + 1; k < len(positions); k++ {
			pNext = append(pNext, positions[k])
		}
		positions = pNext

		if nm == tttPb.Mark_X {
			nm = tttPb.Mark_O
		} else {
			nm = tttPb.Mark_X
		}
	}
	return script
}

func runScript(script []tttPb.MoveTrxPayload, stub *shim.MockStub) ([]pb.Response, error) {
	responses := make([]pb.Response, len(script))
	for i := range script {
		payload := script[i]
		r, err := newMoveArgsBuilder(strconv.Itoa(i), payload.Position, payload.Mark).
			invoke(stub)

		responses[i] = r
		if err != nil {
			return responses, err
		}
	}

	return responses, nil
}

type trxArgsBuilder struct {
	trxArgs *tttPb.TrxArgs
	uuid    string
}

func newMoveArgsBuilder(uuid string, posID int32, mark tttPb.Mark) *trxArgsBuilder {

	return &trxArgsBuilder{
		uuid: uuid,
		trxArgs: &tttPb.TrxArgs{
			Type: tttPb.TrxType_MOVE,
			MovePayload: &tttPb.MoveTrxPayload{
				Position: posID,
				Mark:     mark,
			},
		}}
}

func (tArgsB *trxArgsBuilder) marshal() ([]byte, error) {
	return proto.Marshal(tArgsB.trxArgs)
}

func (tArgsB *trxArgsBuilder) invoke(stub *shim.MockStub) (pb.Response, error) {
	invokeArgs, err := tArgsB.marshal()
	if err != nil {
		return pb.Response{}, fmt.Errorf("Error creating invoke args. %s", err.Error())
	}
	r := stub.MockInvoke(tArgsB.uuid, [][]byte{[]byte("foo"), invokeArgs})
	return r, nil
}
