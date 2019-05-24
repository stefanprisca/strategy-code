package main

import (
	"fmt"
	"log"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

// Dummy struct for hyperledger
type GameContract struct {
}

func (gc *GameContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return HandleInit(APIstub)
}

func (gc *GameContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	return HandleInvoke(APIstub)
}

var AllianceCollection = "alliances"

func getAllianceLedgerKey(cID uint32) string {
	return fmt.Sprintf("alliance%v", cID)
}

func HandleInit(APIstub shim.ChaincodeStubInterface) pb.Response {
	log.Println(APIstub.GetArgs())
	protoArgs := APIstub.GetArgs()[1]
	allianceData := &tfcPb.AllianceData{}
	err := proto.Unmarshal(protoArgs, allianceData)
	if err != nil {
		return shim.Error(
			fmt.Sprintf("could not unmarshal arguments proto message <%v>: %s", protoArgs, err))
	}

	allianceData.State = tfcPb.AllianceState_ACTIVE

	// the lifespan will be reduced by one after the first next,
	// so just add one up. No harm done
	allianceData.Lifespan++
	protoData, err := proto.Marshal(allianceData)
	if err != nil {
		return shim.Error(
			fmt.Sprintf("could not marshal the alliance data <%v>: %s", allianceData, err))
	}

	allianceContractID := allianceData.ContractID
	ledgerKey := getAllianceLedgerKey(allianceContractID)
	// Cannot put private data in the init...
	log.Println("Putting state on ledger....")
	err = APIstub.PutState(ledgerKey, protoData)
	if err != nil {
		return shim.Error(
			fmt.Sprintf("could not save the state on the ledger: %s", err))
	}

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

	newAllianceData := reduceAllianceTerms(*allianceData, trxArgs.CompletedTrxArgs)
	newAllianceData = reduceLifespan(newAllianceData, trxArgs)

	newAllianceState := computeNextAllianceState(newAllianceData, trxArgs.State)

	newAllianceData.State = newAllianceState
	protoData, err := proto.Marshal(&newAllianceData)
	if err != nil {
		return shim.Error(
			fmt.Sprintf("could not marshal the alliance data <%v>: %s", newAllianceData, err))
	}

	APIstub.PutState(ledgerKey, protoData)
	return shim.Success(protoData)

}

func reduceAllianceTerms(allianceData tfcPb.AllianceData, trxArgs *tfcPb.GameContractTrxArgs) tfcPb.AllianceData {

	newTerms := []*tfcPb.GameContractTrxArgs{}
	for i := range allianceData.Terms {
		if proto.Equal(allianceData.Terms[i], trxArgs) {
			continue
		}
		newTerms = append(newTerms, allianceData.Terms[i])
	}

	allianceData.Terms = newTerms
	return allianceData
}

func reduceLifespan(allianceData tfcPb.AllianceData, args *tfcPb.TrxCompletedArgs) tfcPb.AllianceData {

	if args.State != allianceData.StartGameState {
		return allianceData
	}

	if args.CompletedTrxArgs.Type != tfcPb.GameTrxType_NEXT {
		return allianceData
	}

	allianceData.Lifespan--
	return allianceData
}

func computeNextAllianceState(allianceData tfcPb.AllianceData, gameState tfcPb.GameState) tfcPb.AllianceState {
	if len(allianceData.Terms) == 0 {
		return tfcPb.AllianceState_COMPLETED
	}

	if allianceData.Lifespan == 0 {
		return tfcPb.AllianceState_FAILED
	}

	return tfcPb.AllianceState_ACTIVE
}

func getAllianceLedgerData(APIstub shim.ChaincodeStubInterface, ledgerKey string) (*tfcPb.AllianceData, error) {

	protoData, err := APIstub.GetState(ledgerKey)
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

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(GameContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}