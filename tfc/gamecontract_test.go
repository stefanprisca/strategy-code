package tfc

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
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

	_, err := newArgsBuilder().
		withJoinArgs(tfcPb.Player_RED).
		invoke(stub)
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

	_, err := newArgsBuilder().
		withBuildSettleArgs(tfcPb.Player_RED, sID).
		invoke(stub)
	require.NoError(t, err)

	_, err = newArgsBuilder().
		withBuildRoadArgs(tfcPb.Player_RED, eID).
		invoke(stub)
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
		_, err := newArgsBuilder().
			withJoinArgs(p).
			invoke(stub)
		require.NoError(t, err)
	}
}

type argsBuilder struct {
	trxArgs *tfcPb.GameContractTrxArgs
}

func newArgsBuilder() *argsBuilder {
	return &argsBuilder{}
}

func (ab *argsBuilder) build() ([][]byte, error) {
	protoArgs, err := proto.Marshal(ab.trxArgs)
	return [][]byte{[]byte("Testing"), protoArgs}, err
}

func (ab *argsBuilder) withJoinArgs(player tfcPb.Player) *argsBuilder {
	pLoad := &tfcPb.JoinTrxPayload{Player: player}
	ab.trxArgs = &tfcPb.GameContractTrxArgs{
		Type:           tfcPb.GameTrxType_JOIN,
		JoinTrxPayload: pLoad,
	}

	return ab
}

func (ab *argsBuilder) withBuildSettleArgs(player tfcPb.Player, sID uint32) *argsBuilder {
	settlePLoad := &tfcPb.BuildSettlePayload{
		Player:   player,
		SettleID: int32(sID),
	}

	ab.trxArgs = &tfcPb.GameContractTrxArgs{
		Type: tfcPb.GameTrxType_DEV,
		BuildTrxPayload: &tfcPb.BuildTrxPayload{
			Type:               tfcPb.BuildType_SETTLE,
			BuildSettlePayload: settlePLoad,
		},
	}

	return ab
}

func (ab *argsBuilder) withBuildRoadArgs(player tfcPb.Player, eID uint32) *argsBuilder {
	roadPLoad := &tfcPb.BuildRoadPayload{
		Player: player,
		EdgeID: int32(eID),
	}

	ab.trxArgs = &tfcPb.GameContractTrxArgs{
		Type: tfcPb.GameTrxType_DEV,
		BuildTrxPayload: &tfcPb.BuildTrxPayload{
			Type:             tfcPb.BuildType_ROAD,
			BuildRoadPayload: roadPLoad,
		},
	}

	return ab
}

func (ab *argsBuilder) invoke(stub *shim.MockStub) (pb.Response, error) {
	protoArgs, err := ab.build()
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
