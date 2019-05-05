package tfc

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
	"github.com/stretchr/testify/require"
)

func initContract(t *testing.T, cUUID string) *shim.MockStub {
	stub := shim.NewMockStub("mockGameContract", new(GameContract))
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

	for _, p := range []tfcPb.Player{
		tfcPb.Player_RED, tfcPb.Player_BLUE, tfcPb.Player_GREEN} {

		// proposal := pb.SignedProposal{ProposalBytes: , Signature:}
		_, err := newArgsBuilder().
			withJoinArgs(p).
			invoke(stub)
		require.NoError(t, err)

	}

	gameData, err := getLedgerData(stub)
	require.NoError(t, err)

	require.Equal(t, tfcPb.GameState_RROLL, gameData.State,
		"unexpected state after one player joined")
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

func (ab *argsBuilder) invoke(stub *shim.MockStub) (pb.Response, error) {
	protoArgs, err := ab.build()
	if err != nil {
		return pb.Response{}, err
	}

	resp := stub.MockInvoke("0001", protoArgs)
	if shim.OK != resp.Status {
		return resp,
			fmt.Errorf("unexpected status: expected %v, got %v", shim.OK, resp.Status)
	}

	return resp, nil
}
