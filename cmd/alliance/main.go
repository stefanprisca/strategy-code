package main

import (
	"fmt"
	"log"

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

	// Test if it can access the collection

	log.Printf("Processing alliance transaction %v", trxArgs)
	fcn := trxArgs.Type

	switch fcn {

	case tfcPb.AllianceTrxType_INIT:
		if !canReadPrivate(APIstub, trxArgs.CollectionID, trxArgs.InitPayload.ContractID) {
			return shim.Success(nil)
		}
		return alli.HandleInit(APIstub)

	case tfcPb.AllianceTrxType_INVOKE:
		if !canReadPrivate(APIstub, trxArgs.CollectionID, trxArgs.InvokePayload.ObserverID) {
			return shim.Success(nil)
		}
		return alli.HandleInvoke(APIstub)
	}

	return shim.Error("unkown transaction type")
}

func canReadPrivate(APIstub shim.ChaincodeStubInterface, collectionID string, alliUUID uint32) bool {

	alliName := fmt.Sprintf("alliance%v", alliUUID)
	_, err := APIstub.GetPrivateData(collectionID, alliName)
	log.Printf("Tested if it can read private data %v: %v", collectionID, err)
	return err != nil
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(AllianceChaincode))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
