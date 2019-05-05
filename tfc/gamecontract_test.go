package tfc

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
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
}
