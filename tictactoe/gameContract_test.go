package main

import (
	"fmt"
	"testing"

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

	if r.Status != shim.OK {
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
		Status:    tttPb.TttContract_XTURN,
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

	if c1.Status != c2.Status {
		return false, fmt.Sprintf("Contract status does not match. C1: < %v >, C2: < %v >",
			c1.Status, c2.Status)
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
	r := newMoveArgsBuilder("001", positionID, mark).
		invoke(stub, t)

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
	newMoveArgsBuilder("001", positionID, markX).
		invoke(stub, t)
	r := newMoveArgsBuilder("002", positionID, tttPb.Mark_O).
		invoke(stub, t)

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
	if actualContract.Status != tttPb.TttContract_XTURN {
		mark = tttPb.Mark_X
	}

	r := newMoveArgsBuilder("001", position, mark).
		invoke(stub, t)

	if r.GetStatus() == shim.OK {
		t.Fatalf("Did not expect invocation to be successful. Turn was < %v >  and mark < %v >", actualContract.Status, mark)
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

func (tArgsB *trxArgsBuilder) invoke(stub *shim.MockStub, t *testing.T) pb.Response {
	invokeArgs, err := tArgsB.marshal()
	if err != nil {
		t.Fatalf("Error creating invoke args. %s", err.Error())
	}

	r := stub.MockInvoke(tArgsB.uuid, [][]byte{invokeArgs})
	return r
}
