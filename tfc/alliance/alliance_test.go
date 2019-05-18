package alliance

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golang/protobuf/proto"

	"github.com/stefanprisca/strategy-code/tfc"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

// Dummy struct for hyperledger
type MockContract struct {
}

func (gc *MockContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return HandleInit(APIstub)
}

func (gc *MockContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	return HandleInvoke(APIstub)
}

func initContract(t *testing.T, cUUID uint32) *shim.MockStub {
	stub := shim.NewMockStub("mockGameContract", new(MockContract))
	if stub == nil {
		t.Fatalf("Failed to init mock")
	}

	allianceData := &tfcPb.AllianceData{
		Lifespan:       3,
		StartGameState: tfcPb.GameState_RTRADE,
		Terms: []*tfcPb.GameContractTrxArgs{tfc.NewArgsBuilder().
			WithTradeArgs(tfcPb.Player_RED, tfcPb.Player_GREEN, tfcPb.Resource_FOREST, 3).
			Args()},
		ContractID: cUUID,
	}

	protoData, err := proto.Marshal(allianceData)
	require.NoError(t, err)

	trxID := fmt.Sprintf("%v", cUUID)
	r := stub.MockInit(trxID, [][]byte{[]byte("test"), protoData})

	if r.GetStatus() != shim.OK {
		t.Fatalf("Could not init the contract. Error: %s", r.Message)
	}
	return stub
}

func TestInitAlliance(t *testing.T) {
	cuuid := uint32(1010101)
	stub := initContract(t, cuuid)

	ledgerKey := getAllianceLedgerKey(cuuid)
	allianceData, err := getAllianceLedgerData(stub, ledgerKey)
	require.NoError(t, err)
	log.Println(allianceData)

}
