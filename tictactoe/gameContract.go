package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/golang/protobuf/proto"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	tttPb "github.com/stefanprisca/strategy-protobufs/tictactoe"
)

// Dummy struct for hyperledger
type GameContract struct {
}

const CONTRACT_STATE_KEY = "contract.tictactoe"

func (gc *GameContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	args, err := APIstub.GetArgsSlice()
	if err != nil {
		errMsg := fmt.Sprintf("Could not get the input arguments. Error: %s", err.Error())
		return shim.Error(errMsg)
	}

	initTrxArgs := &tttPb.InitTrxArgs{}
	err = proto.Unmarshal(args, initTrxArgs)
	if err != nil {
		errMsg := fmt.Sprintf("Could not unmarshal the input arguments. Error: %s", err.Error())
		return shim.Error(errMsg)
	}

	positions := make([]tttPb.Mark, 9)
	for i := 0; i < 9; i++ {
		positions[i] = tttPb.Mark_E
	}

	tttContract := &tttPb.TttContract{
		Status:    tttPb.TttContract_XTURN,
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

func (gc *GameContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	_, err := APIstub.GetArgsSlice()
	if err != nil {
		errMsg := fmt.Sprintf("Could not get the input arguments. Error: %s", err.Error())
		return shim.Error(errMsg)
	}
	return shim.Success(nil)
}

type moveArgs struct {
	m   string
	pId string
}

var positionRegexp = regexp.MustCompile(`[1|2|3][1|2|3]`)
var markRegexp = regexp.MustCompile(fmt.Sprintf("[%s|%s]", tttPb.Mark_X.String(), tttPb.Mark_O.String()))

func (gc *GameContract) parseMoveArgs(args []string) (moveArgs, error) {
	if len(args) != 2 {
		errMsg := fmt.Sprintf("Wrong number of arguments!. Expected 2, got %d", len(args))
		return moveArgs{}, errors.New(errMsg)
	}

	posID := args[0]
	if !positionRegexp.MatchString(posID) {
		errMsg := fmt.Sprintf("Unkown position %s, must match regular expression <%s>", posID, positionRegexp.String())
		return moveArgs{}, errors.New(errMsg)
	}

	mark := args[1]
	if !markRegexp.MatchString(mark) {
		errMsg := fmt.Sprintf("Unkown mark %s, must match regular expression <%s>", mark, markRegexp.String())
		return moveArgs{}, errors.New(errMsg)
	}

	return moveArgs{pId: posID, m: mark}, nil
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(GameContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
