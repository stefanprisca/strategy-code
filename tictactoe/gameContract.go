package main

import (
	"fmt"
	"regexp"

	"github.com/golang/protobuf/proto"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	tttPb "github.com/stefanprisca/strategy-protobufs/tictactoe"
)

// Dummy struct for hyperledger
type GameContract struct {
}

const CONTRACT_STATE_KEY = "contract.tictactoe"

func (gc *GameContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
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

func (gc *GameContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	creator, errc := APIstub.GetCreator()
	if errc == nil {
		fmt.Println("Creator: ", string(creator))
	}

	protoTrxArgs := APIstub.GetArgs()[0]

	trxArgs := &tttPb.TrxArgs{}
	err := proto.Unmarshal(protoTrxArgs, trxArgs)
	if err != nil {
		errMsg := fmt.Sprintf("Could not parse transaction args. Error %s", err.Error())
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

	if err = validateMoveArgs(APIstub, *contract, *payload); err != nil {
		return shim.Error(err.Error())
	}

	newContract := applyMove(*contract, *payload)
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

var positionRegexp = regexp.MustCompile(fmt.Sprintf("^[1-9]%v$", tttPb.Mark_E))

func positionValidationString(position int32, mark tttPb.Mark) string {
	return fmt.Sprintf("%d%v", position, mark)
}

var turnRegexp = regexp.MustCompile(
	fmt.Sprintf("^(%s|%s)$",
		turnValidationString(tttPb.Mark_X, tttPb.TttContract_XTURN),
		turnValidationString(tttPb.Mark_O, tttPb.TttContract_OTURN)))

func turnValidationString(mark tttPb.Mark, status tttPb.TttContract_Status) string {
	return fmt.Sprintf("%v%v", mark, status)
}

func validateMoveArgs(APIstub shim.ChaincodeStubInterface, contract tttPb.TttContract, payload tttPb.MoveTrxPayload) error {
	pvs := positionValidationString(payload.Position,
		contract.Positions[payload.Position])
	if !positionRegexp.MatchString(pvs) {
		return fmt.Errorf("Invalid position or position not empty. Position < %d >, mark < %v >",
			payload.Position, contract.Positions[payload.Position])
	}

	tvs := turnValidationString(payload.Mark, contract.Status)
	if !turnRegexp.MatchString(tvs) {
		return fmt.Errorf("Invalid turn. Got mark < %v >, expected < %v >", payload.Mark, contract.Status)
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

func applyMove(contract tttPb.TttContract, payload tttPb.MoveTrxPayload) tttPb.TttContract {
	newPositions := contract.Positions
	newPositions[payload.Position] = payload.Mark
	nextStatus := computeNextStatus(newPositions, contract.Status)
	return tttPb.TttContract{
		Positions: newPositions,
		Status:    nextStatus,
		XPlayer:   contract.XPlayer,
		OPlayer:   contract.OPlayer,
	}
}

func computeNextStatus(positions []tttPb.Mark, status tttPb.TttContract_Status) tttPb.TttContract_Status {
	// TODO: apply win functions
	return tttPb.TttContract_OTURN
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(GameContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
