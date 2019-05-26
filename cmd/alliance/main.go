package main

import (
	"fmt"
	"regexp"

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

	fcn := trxArgs.Type
	switch fcn {
	case tfcPb.AllianceTrxType_INIT:
		collection, err := getCollection(trxArgs.InitPayload.Allies)
		if err != nil {
			return shim.Error(err.Error())
		}
		return alli.HandleInit(APIstub, collection)
	case tfcPb.AllianceTrxType_INVOKE:
		collection, err := getCollection(trxArgs.InitPayload.Allies)
		if err != nil {
			return shim.Error(err.Error())
		}

		return alli.HandleInvoke(APIstub, collection)
	}

	return shim.Error("unkown transaction type")
}

var RGColRegex = regexp.MustCompile(fmt.Sprintf("%v%v|%v%v",
	tfcPb.Player_RED, tfcPb.Player_GREEN, tfcPb.Player_GREEN, tfcPb.Player_RED))

var GBColRegex = regexp.MustCompile(fmt.Sprintf("%v%v|%v%v",
	tfcPb.Player_BLUE, tfcPb.Player_GREEN, tfcPb.Player_GREEN, tfcPb.Player_BLUE))

var RBColRegex = regexp.MustCompile(fmt.Sprintf("%v%v|%v%v",
	tfcPb.Player_RED, tfcPb.Player_BLUE, tfcPb.Player_BLUE, tfcPb.Player_RED))

func colSelectionStr(allies []tfcPb.Player) string {
	return fmt.Sprintf("%v%v", allies[0], allies[1])
}

func getCollection(allies []tfcPb.Player) (string, error) {
	selStr := colSelectionStr(allies)
	switch {
	case RGColRegex.MatchString(selStr):
		return "redgreen", nil
	case GBColRegex.MatchString(selStr):
		return "greenblue", nil
	case RBColRegex.MatchString(selStr):
		return "redblue", nil
	}
	return "", fmt.Errorf("unkown collection for allies %v", allies)

}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(AllianceChaincode))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
