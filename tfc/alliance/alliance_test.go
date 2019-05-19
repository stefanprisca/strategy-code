package alliance

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golang/protobuf/proto"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/stefanprisca/strategy-code/tfc"
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

func initContract(t *testing.T, cUUID uint32, terms ...*tfcPb.GameContractTrxArgs) *shim.MockStub {
	stub := shim.NewMockStub("mockGameContract", new(MockContract))
	if stub == nil {
		t.Fatalf("Failed to init mock")
	}

	initBuilder := newInitArgsBuilder(cUUID)
	for _, t := range terms {
		initBuilder.AddTerm(t)
	}

	err := initBuilder.InitMock(stub)
	require.NoError(t, err)
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

func TestAllianceCompletes(t *testing.T) {
	cuuid := uint32(1010101)

	term := tfc.NewArgsBuilder().
		WithTradeArgs(tfcPb.Player_RED, tfcPb.Player_GREEN, tfcPb.Resource_FOREST, 3).
		Args()

	stub := initContract(t, cuuid, term)
	err := newTrxCompletedBuilder().
		WithArgs(term).
		MockInvoke(stub, cuuid)

	require.NoError(t, err)

	ledgerKey := getAllianceLedgerKey(cuuid)
	allianceData, err := getAllianceLedgerData(stub, ledgerKey)
	require.NoError(t, err)
	log.Println(allianceData)

	require.Empty(t, allianceData.Terms)
	require.Equal(t, allianceData.State, tfcPb.AllianceState_COMPLETED)
}

func TestAllianceFails(t *testing.T) {
	cuuid := uint32(1010101)

	initTerm := tfc.NewArgsBuilder().
		WithTradeArgs(tfcPb.Player_RED, tfcPb.Player_GREEN, tfcPb.Resource_FOREST, 3).
		Args()
	stub := initContract(t, cuuid, initTerm)

	// randTrx := tfc.NewArgsBuilder().
	// 	WithTradeArgs(tfcPb.Player_RED, tfcPb.Player_BLUE, tfcPb.Resource_HILL, 5).
	// 	Args()
	nextTrx := tfc.NewArgsBuilder().
		WithNextArgs().Args()

	var err error
	for l := int32(0); l < lifespan+1; l++ {
		err = newTrxCompletedBuilder().
			WithArgs(nextTrx).
			MockInvoke(stub, cuuid)
		require.NoError(t, err)
	}

	ledgerKey := getAllianceLedgerKey(cuuid)
	allianceData, err := getAllianceLedgerData(stub, ledgerKey)
	require.NoError(t, err)
	log.Println(allianceData)

	require.NotEmpty(t, allianceData.Terms)
	require.Equal(t, allianceData.State, tfcPb.AllianceState_FAILED)
}

type initArgsBuilder struct {
	allianceData *tfcPb.AllianceData
}

var lifespan = int32(3)

func newInitArgsBuilder(cUUID uint32) *initArgsBuilder {
	ad := &tfcPb.AllianceData{
		Lifespan:       lifespan,
		StartGameState: tfcPb.GameState_RTRADE,
		Terms:          []*tfcPb.GameContractTrxArgs{},
		ContractID:     cUUID,
	}
	return &initArgsBuilder{ad}
}

func (ab *initArgsBuilder) AddTerm(term *tfcPb.GameContractTrxArgs) *initArgsBuilder {
	ab.allianceData.Terms = append(ab.allianceData.Terms, term)
	return ab
}

func (ab *initArgsBuilder) InitMock(stub *shim.MockStub) error {
	protoData, err := proto.Marshal(ab.allianceData)
	if err != nil {
		return err
	}

	trxID := fmt.Sprintf("%v", ab.allianceData.ContractID)
	r := stub.MockInit(trxID, [][]byte{[]byte("test"), protoData})

	if r.GetStatus() != shim.OK {
		return fmt.Errorf("Could not init the contract. Error: %s", r.Message)
	}
	return nil
}

type trxCompletedBuilder struct {
	trxCompletedArgs *tfcPb.TrxCompletedArgs
}

func newTrxCompletedBuilder() *trxCompletedBuilder {
	return &trxCompletedBuilder{
		trxCompletedArgs: &tfcPb.TrxCompletedArgs{
			State: tfcPb.GameState_RTRADE,
		},
	}
}

func (tcb *trxCompletedBuilder) WithArgs(arg *tfcPb.GameContractTrxArgs) *trxCompletedBuilder {
	tcb.trxCompletedArgs.CompletedTrxArgs = arg
	return tcb
}

func (tcb *trxCompletedBuilder) MockInvoke(stub *shim.MockStub, observerID uint32) error {
	tcb.trxCompletedArgs.ObserverID = observerID
	protoData, err := proto.Marshal(tcb.trxCompletedArgs)
	if err != nil {
		return err
	}

	n := rand.Int63()
	uuid := strconv.FormatInt(n, 8)
	r := stub.MockInvoke(uuid, [][]byte{[]byte("test"), protoData})

	if r.GetStatus() != shim.OK {
		return fmt.Errorf("Could not init the contract. Error: %s", r.Message)
	}
	return nil
}
