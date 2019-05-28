package main

import (
	"fmt"
	"regexp"
	"strings"

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
	collection, err := getCollection(trxArgs.OrgIDs)
	if err != nil {
		return shim.Error(err.Error())
	}

	switch fcn {
	case tfcPb.AllianceTrxType_INIT:
		return alli.HandleInit(APIstub, collection)
	case tfcPb.AllianceTrxType_INVOKE:
		return alli.HandleInvoke(APIstub, collection)
	}

	return shim.Error("unkown transaction type")
}

func colSelectionRegexp(orgIDs []string) *regexp.Regexp {
	org1ID := strings.ToLower(orgIDs[0])
	org2ID := strings.ToLower(orgIDs[1])
	return regexp.MustCompile(fmt.Sprintf("al(%v%v|%v%v)",
		org1ID, org2ID, org2ID, org1ID))
}

func getCollection(orgIDs []string) (string, error) {
	selRegexp := colSelectionRegexp(orgIDs)
	switch {
	case selRegexp.MatchString("alplayer1player2"):
		return "alplayer1player2", nil
	case selRegexp.MatchString("alplayer1player3"):
		return "alplayer1player2", nil
	case selRegexp.MatchString("alplayer1player4"):
		return "alplayer1player2", nil
	case selRegexp.MatchString("alplayer1player5"):
		return "alplayer1player2", nil
	case selRegexp.MatchString("alplayer2player3"):
		return "alplayer1player2", nil
	case selRegexp.MatchString("alplayer2player4"):
		return "alplayer1player2", nil
	case selRegexp.MatchString("alplayer2player5"):
		return "alplayer1player2", nil
	case selRegexp.MatchString("alplayer3player4"):
		return "alplayer1player2", nil
	case selRegexp.MatchString("alplayer3player5"):
		return "alplayer1player2", nil
	case selRegexp.MatchString("alplayer4player5"):
		return "alplayer1player2", nil
	}
	return "", fmt.Errorf("unkown collection for orgs %v", orgIDs)

}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(AllianceChaincode))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
