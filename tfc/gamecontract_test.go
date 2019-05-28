package tfc

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"strconv"
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

var playerSignedProposals = map[tfcPb.Player]*pb.SignedProposal{
	tfcPb.Player_RED:   &pb.SignedProposal{ProposalBytes: []byte{}, Signature: []byte("red")},
	tfcPb.Player_BLUE:  &pb.SignedProposal{ProposalBytes: []byte{}, Signature: []byte("blue")},
	tfcPb.Player_GREEN: &pb.SignedProposal{ProposalBytes: []byte{}, Signature: []byte("green")},
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
	p := tfcPb.Player_RED
	pP := playerSignedProposals[tfcPb.Player_RED]
	_, err := NewArgsBuilder().
		WithJoinArgs(p).
		invokeSignedMock(stub, pP)
	require.NoError(t, err)

	gameData, err := getLedgerData(stub)
	require.NoError(t, err)
	expectedId := GetPlayerId(tfcPb.Player_RED)
	_, ok := gameData.IdentityMap[expectedId]
	require.True(t, ok,
		fmt.Sprintf("expected to find player id for %v after join operation.", tfcPb.Player_RED))

	require.Equal(t, tfcPb.GameState_JOINING, gameData.State,
		"unexpected state after one player joined")

	profileNotNil := gameData.Profiles[expectedId] != nil
	require.True(t, profileNotNil,
		"expected profile to be initialized")

}

func TestRGBJoinGame(t *testing.T) {
	cUUID := "01010101"
	stub := initContract(t, cUUID)
	newSI(stub).joinRGB(playerSignedProposals)

	gameData, err := getLedgerData(stub)
	require.NoError(t, err)

	require.Equal(t, tfcPb.GameState_RROLL, gameData.State,
		"unexpected state after one player joined")

	redID := GetPlayerId(tfcPb.Player_RED)
	sign := playerSignedProposals[tfcPb.Player_RED].Signature
	signCs := crc32.ChecksumIEEE(sign)
	expectedSign := []byte(fmt.Sprintf("%d", signCs))
	actualSign := gameData.IdentityMap[redID]
	require.Equal(t, expectedSign, actualSign)
}

func TestTrade(t *testing.T) {
	cUUID := "01010101"
	stub := initContract(t, cUUID)

	err := newSI(stub).
		joinRGB(playerSignedProposals).
		roll(tfcPb.Player_RED).
		getError()
	require.NoError(t, err)

	src, dest := tfcPb.Player_RED, tfcPb.Player_BLUE
	resource := tfcPb.Resource_HILL
	amount := int32(2)
	assertCorrectTrade(t, stub, src, dest, resource, amount)
	assertCorrectTrade(t, stub, src, dest, resource, -amount)
}

func assertCorrectTrade(t *testing.T, stub *shim.MockStub,
	src, dest tfcPb.Player, resource tfcPb.Resource, amount int32) {

	gameData, err := getLedgerData(stub)
	require.NoError(t, err)

	preSrcProfile := gameData.Profiles[GetPlayerId(src)]
	preDestProfile := gameData.Profiles[GetPlayerId(dest)]

	rID := GetResourceId(resource)

	expectedSrcA := preSrcProfile.Resources[rID] - amount
	expectedDestA := preDestProfile.Resources[rID] + amount

	_, err = NewArgsBuilder().
		WithTradeArgs(src, dest,
			resource, amount).
		invokeSignedMock(stub, playerSignedProposals[src])
	require.NoError(t, err)

	gameData, err = getLedgerData(stub)
	require.NoError(t, err)

	postSrcProfile := gameData.Profiles[GetPlayerId(src)]
	postDestProfile := gameData.Profiles[GetPlayerId(dest)]

	actualSrcA := postSrcProfile.Resources[rID]
	actualDestA := postDestProfile.Resources[rID]

	require.Equal(t, expectedSrcA, actualSrcA,
		"source amount invalid after trade")
	require.Equal(t, expectedDestA, actualDestA,
		"source amount invalid after trade")
}

func TestBuildSettle(t *testing.T) {
	cUUID := "01010101"
	stub := initContract(t, cUUID)

	err := newSI(stub).
		joinRGB(playerSignedProposals).
		roll(tfcPb.Player_RED).
		next(tfcPb.Player_RED).
		getError()

	require.NoError(t, err)

	sID := pointHash(tfc.Coord{X: 0, Y: 0})
	eID := edgeHash(tfc.Coord{X: 0, Y: 0}, N)

	player := tfcPb.Player_RED
	proposal := playerSignedProposals[player]
	_, err = NewArgsBuilder().
		WithBuildSettleArgs(player, sID).
		invokeSignedMock(stub, proposal)
	require.NoError(t, err)

	_, err = NewArgsBuilder().
		WithBuildRoadArgs(player, eID).
		invokeSignedMock(stub, proposal)
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

type scriptStep struct {
	ab *ArgsBuilder
	p  tfcPb.Player
}

func scriptTFC() []scriptStep {

	p1C, p2C, p3C := tfcPb.Player_RED, tfcPb.Player_GREEN, tfcPb.Player_BLUE

	s := []scriptStep{
		{NewArgsBuilder().WithJoinArgs(p2C), p2C},
		{NewArgsBuilder().WithJoinArgs(p1C), p1C},
		{NewArgsBuilder().WithJoinArgs(p3C), p3C},
		{NewArgsBuilder().WithRollArgs(), p1C},
		{NewArgsBuilder().WithTradeArgs(p1C, p2C, tfcPb.Resource_CAMP, 2), p1C},
		{NewArgsBuilder().WithTradeArgs(p1C, p3C, tfcPb.Resource_HILL, -2), p1C},
		{NewArgsBuilder().WithNextArgs(), p1C},
		{NewArgsBuilder().WithNextArgs(), p1C},
		{NewArgsBuilder().WithRollArgs(), p2C},
		{NewArgsBuilder().WithTradeArgs(p2C, p1C, tfcPb.Resource_CAMP, 2), p2C},
		{NewArgsBuilder().WithTradeArgs(p2C, p3C, tfcPb.Resource_PASTURE, -2), p2C},
		{NewArgsBuilder().WithNextArgs(), p2C},
		{NewArgsBuilder().WithNextArgs(), p2C},
		{NewArgsBuilder().WithRollArgs(), p3C},
		{NewArgsBuilder().WithTradeArgs(p3C, p1C, tfcPb.Resource_HILL, -2), p3C},
		{NewArgsBuilder().WithTradeArgs(p3C, p2C, tfcPb.Resource_PASTURE, -2), p3C},
		{NewArgsBuilder().WithNextArgs(), p3C},
		{NewArgsBuilder().WithNextArgs(), p3C},
	}

	for i := 0; i < 5; i++ {
		s = append(s, s[3:]...)
	}

	return s
}

func TestTFCScript(t *testing.T) {
	cUUID := "01010101"
	stub := initContract(t, cUUID)

	script := scriptTFC()

	for _, ss := range script {
		_, err := ss.ab.invokeSignedMock(stub, playerSignedProposals[ss.p])
		require.NoError(t, err)
	}
}

type stateIterator struct {
	stub *shim.MockStub
	err  error
}

func newSI(stub *shim.MockStub) *stateIterator {
	return &stateIterator{stub, nil}
}

func (si *stateIterator) getError() error {
	return si.err
}

func (si *stateIterator) joinRGB(signedProposals map[tfcPb.Player]*pb.SignedProposal) *stateIterator {

	if si.err != nil {
		return si
	}

	for _, p := range []tfcPb.Player{
		tfcPb.Player_RED, tfcPb.Player_BLUE, tfcPb.Player_GREEN} {
		proposal := signedProposals[p]
		_, err := NewArgsBuilder().
			WithJoinArgs(p).
			invokeSignedMock(si.stub, proposal)
		if err != nil {
			si.err = err
			return si
		}
	}

	return si
}

func (si *stateIterator) roll(p tfcPb.Player) *stateIterator {
	if si.err != nil {
		return si
	}
	pP := playerSignedProposals[p]
	_, err := NewArgsBuilder().
		WithRollArgs().
		invokeSignedMock(si.stub, pP)
	si.err = err
	return si
}

func (si *stateIterator) next(p tfcPb.Player) *stateIterator {
	if si.err != nil {
		return si
	}

	pP := playerSignedProposals[p]
	_, err := NewArgsBuilder().
		WithNextArgs().
		invokeSignedMock(si.stub, pP)

	si.err = err
	return si
}

func (ab *ArgsBuilder) invokeSignedMock(stub *shim.MockStub, sp *pb.SignedProposal) (pb.Response, error) {
	protoArgs, err := ab.Build()
	if err != nil {
		return pb.Response{}, err
	}

	n := rand.Int63()
	uuid := strconv.FormatInt(n, 8)

	trxargs := [][]byte{[]byte("Test")}
	trxargs = append(trxargs, protoArgs[0])
	resp := stub.MockInvokeWithSignedProposal(uuid, trxargs, sp)
	if shim.OK != resp.Status {
		return resp,
			fmt.Errorf("unexpected status: expected %v, got %v. message: %s",
				shim.OK, resp.Status, resp.Message)
	}

	return resp, nil
}
