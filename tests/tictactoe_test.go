package tictactoe

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/stefanprisca/strategy-code/tests/utils"
	"github.com/stefanprisca/strategy-code/tictactoe"
)

func TestInit(t *testing.T) {
	st := make(map[string][]byte)
	stub := utils.TestCCStub{State: st}

	gc := tictactoe.GameContract{}
	r := gc.Init(stub)
	if r.GetStatus() != shim.OK {
		t.Logf("Initializing was not succesfull.")
		t.FailNow()
	}

	for _, pId := range gc.Positions {
		p, err := tictactoe.GetPosition(stub, pId)
		if err != nil {
			t.Logf("Failed to retrieve position < %s > with error: %s", pId, err.Error())
			t.FailNow()
		}

		if p.Mark != tictactoe.Empty {
			t.Logf("Unexpected postion mark: < %s >. Expected < %s >.", p.Mark, tictactoe.Empty)
			t.FailNow()
		}

		pTerm := tictactoe.ToTerm(p)
		if !(pTerm(tictactoe.X) && pTerm(tictactoe.O)) {
			t.Logf("Position < %s > not open.", pId)
		}
	}
}

func TestInvokeMove(t *testing.T) {
	st := make(map[string][]byte)
	stub := utils.TestCCStub{State: st}

	gc := tictactoe.GameContract{}
	gc.Init(stub)

	pId := "11"
	m := tictactoe.X

	stub.SetFunctionAndParameters("move", pId, m)
	r := gc.Invoke(stub)
	if r.GetStatus() != shim.OK {
		t.Logf("Could not invoke move function, error: %s", r.GetMessage())
		t.FailNow()
	}

	p, err := tictactoe.GetPosition(stub, pId)
	if err != nil {
		t.Logf("Failed to retrieve position < %s > with error: %s", pId, err.Error())
		t.FailNow()
	}

	if p.Mark != m {
		t.Logf("Wrong mark.  Expected <%s>, found: <%s>", m, p.Mark)
		t.FailNow()
	}
}

func TestInvokeMoveOnOccupiedPos(t *testing.T) {
	st := make(map[string][]byte)
	stub := utils.TestCCStub{State: st}

	gc := tictactoe.GameContract{}
	gc.Init(stub)

	pId := "11"
	m := tictactoe.X

	stub.SetFunctionAndParameters("move", pId, m)
	gc.Invoke(stub)

	m2 := tictactoe.O
	stub.SetFunctionAndParameters("move", pId, m2)

	r := gc.Invoke(stub)
	if r.GetStatus() != shim.ERROR {
		t.Logf("Expected to get error after invoking move twice on the same position.")
		t.FailNow()
	}

	p, err := tictactoe.GetPosition(stub, pId)
	if err != nil {
		t.Logf("Failed to retrieve position < %s > with error: %s", pId, err.Error())
		t.FailNow()
	}

	if p.Mark != m {
		t.Logf("Wrong mark. Expected <%s>, found: <%s>", m, p.Mark)
		t.FailNow()
	}
}
