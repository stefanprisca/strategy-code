package alliance

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

var AllianceCollection = "alliances"

func getAllianceLedgerKey(cID uint32) string {
	return fmt.Sprintf("alliance%v", cID)
}

func HandleInit(APIstub shim.ChaincodeStubInterface) pb.Response {
	protoArgs := APIstub.GetArgs()[1]
	allianceData := &tfcPb.AllianceData{}
	err := proto.Unmarshal(protoArgs, allianceData)
	if err != nil {
		return shim.Error(
			fmt.Sprintf("could not unmarshal arguments proto message <%v>: %s", protoArgs, err))
	}

	allianceData.State = tfcPb.AllianceState_ACTIVE
	protoData, err := proto.Marshal(allianceData)
	if err != nil {
		return shim.Error(
			fmt.Sprintf("could not marshal the alliance data <%v>: %s", allianceData, err))
	}

	allianceContractID := allianceData.ContractID
	ledgerKey := getAllianceLedgerKey(allianceContractID)
	APIstub.PutPrivateData(AllianceCollection, ledgerKey, protoData)
	return shim.Success(protoData)
}

func HandleInvoke(APIstub shim.ChaincodeStubInterface) pb.Response {

	protoArgs := APIstub.GetArgs()[1]
	trxArgs := &tfcPb.TrxCompletedArgs{}
	err := proto.Unmarshal(protoArgs, trxArgs)
	if err != nil {
		return shim.Error(
			fmt.Sprintf("could not unmarshal arguments proto message <%v>: %s", protoArgs, err))
	}

	// TODO: Assert preconditions
	allianceContractID := trxArgs.ObserverID
	ledgerKey := getAllianceLedgerKey(allianceContractID)

	allianceData, err := getAllianceLedgerData(APIstub, ledgerKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	newAllianceData, err := reduceAllianceTerms(*allianceData, *trxArgs.CompletedTrxArgs)
	if err != nil {
		return shim.Error(err.Error())
	}

	newAllianceState, err := computeNextAllianceState(newAllianceData, trxArgs.State)
	if err != nil {
		return shim.Error(err.Error())
	}

	newAllianceData.State = newAllianceState
	protoData, err := proto.Marshal(allianceData)
	if err != nil {
		return shim.Error(
			fmt.Sprintf("could not marshal the alliance data <%v>: %s", allianceData, err))
	}

	APIstub.PutPrivateData(AllianceCollection, ledgerKey, protoData)
	return shim.Success(protoData)

}

func reduceAllianceTerms(allianceData tfcPb.AllianceData, trxArgs tfcPb.GameContractTrxArgs) (tfcPb.AllianceData, error) {
	return allianceData, nil
}

func computeNextAllianceState(allianceData tfcPb.AllianceData, gameState tfcPb.GameState) (tfcPb.AllianceState, error) {
	return tfcPb.AllianceState_ACTIVE, nil
}

func getAllianceLedgerData(APIstub shim.ChaincodeStubInterface, ledgerKey string) (*tfcPb.AllianceData, error) {

	protoData, err := APIstub.GetPrivateData(AllianceCollection, ledgerKey)
	if err != nil {
		return nil, err
	}

	allianceData := &tfcPb.AllianceData{}
	err = proto.Unmarshal(protoData, allianceData)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the proto contract. Error: %s", err.Error())
	}
	return allianceData, nil
}
