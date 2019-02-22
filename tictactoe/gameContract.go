package tictactoe

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

const MoveTransaction = "move"

type GameContract struct {
	Positions []string
}

const (
	X     string = "X"
	O     string = "O"
	Empty string = "Empty"
)

type Position struct {
	ID   string `json:"id"`
	Mark string `json:"mark"`
}

type posFunc func(m string) bool

func (gc *GameContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	gc.Positions = make([]string, 9)
	positions := []string{"1", "2", "3"}
	for i, p1 := range positions {
		for j, p2 := range positions {
			id := p1 + p2
			p := Position{id, Empty}
			st, err := json.Marshal(p)
			if err != nil {
				fmt.Errorf("Could not marshal position: %s", err.Error())
				return shim.Error(err.Error())
			}
			APIstub.PutState(id, st)
			gc.Positions[i*3+j] = id
		}
	}

	return shim.Success([]byte("Succesfully initialized the GameBoard for tictactoe"))
}

func (gc *GameContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	if function != MoveTransaction {
		errMsg := fmt.Sprintf("Unkown transaction name: %s", function)
		fmt.Errorf(errMsg)
		return shim.Error(errMsg)
	}

	return gc.move(APIstub, args)
}

type moveArgs struct {
	m   string
	pId string
}

func containsString(l []string, x string) bool {
	found := false
	for _, y := range l {
		found = found || x == y
	}
	return found
}

func (gc *GameContract) parseMoveArgs(args []string) (moveArgs, error) {
	if len(args) != 2 {
		errMsg := fmt.Sprintf("Wrong number of arguments!. Expected 2, got %d", len(args))
		return moveArgs{}, errors.New(errMsg)
	}

	posId := args[0]
	if !containsString(gc.Positions, posId) {
		errMsg := fmt.Sprintf("Unkown position %s, expected one of %v", posId, gc.Positions)
		return moveArgs{}, errors.New(errMsg)
	}

	mark := args[1]
	if mark != X && mark != O && mark != Empty {
		errMsg := fmt.Sprintf("Unkown mark %s, expected one of [%s, %s, %s]", posId, X, O, Empty)
		return moveArgs{}, errors.New(errMsg)
	}

	return moveArgs{pId: posId, m: mark}, nil
}

func (gc *GameContract) apply(APIstub shim.ChaincodeStubInterface, moveArgs moveArgs) (bool, error) {
	p, err := GetPosition(APIstub, moveArgs.pId)
	if err != nil {
		return false, err
	}

	posF := ToTerm(p)
	if !posF(moveArgs.m) {
		return false, fmt.Errorf("Position %s is taken", moveArgs.pId)
	}

	newPos := Position{ID: moveArgs.pId, Mark: moveArgs.m}
	st, err := json.Marshal(newPos)
	if err != nil {
		return false, fmt.Errorf("Could not marshal position: %s", err.Error())
	}
	err = APIstub.PutState(moveArgs.pId, st)

	return true, err
}

func (gc *GameContract) move(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	moveArgs, err := gc.parseMoveArgs(args)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to parse arguments: %v, Parse error: %s", args, err.Error())
		return shim.Error(errMsg)
	}

	success, err := gc.apply(APIstub, moveArgs)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to apply move arguments: %v, Movement error: %s", args, err.Error())
		return shim.Error(errMsg)
	}

	// # TODO # Implement the missing functionality
	if success {
		// Apply win function
		// Marshal the game board
		return shim.Success(nil)
	}

	return shim.Error("Something went wrong, try again.")
}

func ToTerm(p Position) posFunc {
	return func(m string) bool { return p.Mark == Empty }
}

func GetPosition(APIstub shim.ChaincodeStubInterface, pId string) (Position, error) {
	data, err := APIstub.GetState(pId)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get position from stub for pId < %s >, error: < %s >", pId, err.Error())
		return Position{}, errors.New(errMsg)
	}

	p := Position{}
	err = json.Unmarshal(data, &p)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to unmarshal position for pId < %s >, error: < %s >, raw data: >> %v <<", pId, err.Error(), data)
		return p, errors.New(errMsg)
	}
	return p, nil
}
