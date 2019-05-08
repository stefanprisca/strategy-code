package tfc

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/stefanprisca/strategy-protobufs/tfc"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
	"github.com/stretchr/testify/require"
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

func initContract(t *testing.T, cUUID string) *shim.MockStub {
	stub := shim.NewMockStub("mockGameContract", new(MockContract))
	if stub == nil {
		t.Fatalf("Failed to init mock")
	}
	r := stub.MockInit(cUUID, [][]byte{})

	if r.GetStatus() != shim.OK {
		t.Fatalf("Could not init the contract. Error: %s", r.Message)
	}
	return stub
}

func TestInitGameContract(t *testing.T) {
	cUUID := "01010101"
	stub := initContract(t, cUUID)
	gameData, err := getLedgerData(stub)
	require.NoError(t, err, "could not get ledger data")

	assertGameBoard(t, *gameData.Board)

	require.Empty(t, gameData.Profiles,
		"expected profiles to be empty")

	require.Equal(t, gameData.State, tfcPb.GameState_JOINING,
		"expected game to be in state %v", tfcPb.GameState_JOINING)

	require.Contains(t, gameData.IdentityMap, ContractID)
	actualUUID := string(gameData.IdentityMap[ContractID])
	require.Equal(t, cUUID, actualUUID)

	// profilesNotNil := gameData.Profiles != nil
	// require.True(t, profilesNotNil,
	// 	"expected profiles to be initialized")
}

func TestJoinGame(t *testing.T) {
	cUUID := "01010101"
	stub := initContract(t, cUUID)

	_, err := NewArgsBuilder().
		WithJoinArgs(tfcPb.Player_RED).
		invokeMock(stub)
	require.NoError(t, err)

	gameData, err := getLedgerData(stub)
	require.NoError(t, err)
	expectedId := GetPlayerId(tfcPb.Player_RED)
	_, ok := gameData.IdentityMap[expectedId]
	require.True(t, ok,
		fmt.Sprintf("expected to find player id for %v after join operation.", tfcPb.Player_RED))

	require.Equal(t, tfcPb.GameState_JOINING, gameData.State,
		"unexpected state after one player joined")
}

func TestRGBJoinGame(t *testing.T) {
	cUUID := "01010101"
	stub := initContract(t, cUUID)
	joinRGB(t, stub)

	gameData, err := getLedgerData(stub)
	require.NoError(t, err)

	require.Equal(t, tfcPb.GameState_RROLL, gameData.State,
		"unexpected state after one player joined")
}

func TestBuildSettle(t *testing.T) {
	cUUID := "01010101"
	stub := initContract(t, cUUID)

	joinRGB(t, stub)

	sID := pointHash(tfc.Coord{X: 0, Y: 0})
	eID := edgeHash(tfc.Coord{X: 0, Y: 0}, N)

	_, err := NewArgsBuilder().
		WithBuildSettleArgs(tfcPb.Player_RED, sID).
		invokeMock(stub)
	require.NoError(t, err)

	_, err = NewArgsBuilder().
		WithBuildRoadArgs(tfcPb.Player_RED, eID).
		invokeMock(stub)
	require.NoError(t, err)

	gameData, err := getLedgerData(stub)
	require.NoError(t, err)

	I := gameData.Board.Intersections[sID]
	expectedSettle := tfcPb.Settlement_REDSETTLE
	actualSettle := I.Attributes.Settlement
	require.Equal(t, expectedSettle, actualSettle,
		"unexpected settlement found after building red settle")

	E := gameData.Board.Edges[eID]
	expectedRoad := tfcPb.Road_REDROAD
	actualRoad := E.Attributes.Road
	require.Equal(t, expectedRoad, actualRoad,
		"unexpected road found after building red road")

	profile := gameData.Profiles[GetPlayerId(tfcPb.Player_RED)]
	for _, r := range profile.Resources {
		require.NotEqual(t, 5, r,
			"expected build to consume resources")
	}
	require.EqualValues(t, 3, profile.WinningPoints,
		"expected build to increase winning points")
}

func joinRGB(t *testing.T, stub *shim.MockStub) {
	for _, p := range []tfcPb.Player{
		tfcPb.Player_RED, tfcPb.Player_BLUE, tfcPb.Player_GREEN} {

		// proposal := pb.SignedProposal{ProposalBytes: , Signature:}
		_, err := NewArgsBuilder().
			WithJoinArgs(p).
			invokeMock(stub)
		require.NoError(t, err)
	}
}

func (ab *ArgsBuilder) invokeMock(stub *shim.MockStub) (pb.Response, error) {
	protoArgs, err := ab.Build()
	if err != nil {
		return pb.Response{}, err
	}

	resp := stub.MockInvoke("0001", protoArgs)
	if shim.OK != resp.Status {
		return resp,
			fmt.Errorf("unexpected status: expected %v, got %v. message: %s",
				shim.OK, resp.Status, resp.Message)
	}

	return resp, nil
}
