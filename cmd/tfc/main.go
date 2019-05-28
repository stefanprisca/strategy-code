package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/stefanprisca/strategy-code/tfc"
)

// Dummy struct for hyperledger
type GameContract struct {
}

func (gc *GameContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return tfc.HandleInit(APIstub)
}

func (gc *GameContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	result := tfc.HandleInvoke(APIstub)

	log.Printf("Saved state on the ledger. \n\n\t ###### State: #####\n %v\n\n", result.Payload)
	return result
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(GameContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
