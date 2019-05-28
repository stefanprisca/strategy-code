package main

import (
	"bytes"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	alli "github.com/stefanprisca/strategy-code/alliance"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

// Dummy struct for hyperledger
type AllianceChaincode struct {
}

func (gc *AllianceChaincode) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (gc *AllianceChaincode) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	/*
		1) If new alliance, then
			1.1) Create new alliance contract with data.ContractID
			1.2) Give it the collection correponding to the allies (RG, RB, GB)
		2) If Invoke, then forward to corresponding ContractID
	*/

	protoArgs := APIstub.GetArgs()[1]
	trxArgs := &tfcPb.AllianceTrxArgs{}
	err := proto.Unmarshal(protoArgs, trxArgs)
	if err != nil {
		return shim.Error(
			fmt.Sprintf("could not unmarshal arguments proto message <%v>: %s", protoArgs, err))
	}

	creatorIsAlly, err := isCreatorAlly(APIstub, trxArgs)
	if err != nil {
		return shim.Error(
			fmt.Sprintf("could not determine if creator is ally: %s", err))
	}

	if !creatorIsAlly {
		return shim.Success([]byte("Creator is not ally, cannot endorse..."))
	}

	fcn := trxArgs.Type

	switch fcn {
	case tfcPb.AllianceTrxType_INIT:
		return alli.HandleInit(APIstub)
	case tfcPb.AllianceTrxType_INVOKE:
		return alli.HandleInvoke(APIstub)
	}

	return shim.Error("unkown transaction type")
}

func isCreatorAlly(APIstub shim.ChaincodeStubInterface, trxArgs *tfcPb.AllianceTrxArgs) (bool, error) {

	ccName := APIstub.GetChannelID()
	r := APIstub.InvokeChaincode(ccName, [][]byte{[]byte("query")}, ccName)
	if r.Status != shim.OK {
		return false, fmt.Errorf("could not get the game data:%s", r.Message)
	}

	gameData := &tfcPb.GameData{}
	err := proto.Unmarshal(r.Payload, gameData)
	if err != nil {
		return false, fmt.Errorf("could not unmarshal game data: %s", err)
	}

	creatorSign, err := APIstub.GetCreator()
	if err != nil {
		return false, fmt.Errorf("could not obtain creator: %s", err)
	}

	for _, a := range trxArgs.Allies {
		if bytes.Equal(gameData.IdentityMap[int32(a)], creatorSign) {
			return true, nil
		}
	}

	return false, nil
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(AllianceChaincode))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
