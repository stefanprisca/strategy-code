package tfc

import (
	"github.com/gogo/protobuf/proto"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

type ArgsBuilder struct {
	trxArgs *tfcPb.GameContractTrxArgs
}

func NewArgsBuilder() *ArgsBuilder {
	return &ArgsBuilder{}
}

func (ab *ArgsBuilder) Build() ([][]byte, error) {
	protoArgs, err := proto.Marshal(ab.trxArgs)
	return [][]byte{protoArgs}, err
}

func (ab *ArgsBuilder) Args() *tfcPb.GameContractTrxArgs {
	return ab.trxArgs
}

func (ab *ArgsBuilder) WithJoinArgs(player tfcPb.Player) *ArgsBuilder {
	pLoad := &tfcPb.JoinTrxPayload{Player: player}
	ab.trxArgs = &tfcPb.GameContractTrxArgs{
		Type:           tfcPb.GameTrxType_JOIN,
		JoinTrxPayload: pLoad,
	}

	return ab
}

func (ab *ArgsBuilder) WithRollArgs() *ArgsBuilder {
	ab.trxArgs = &tfcPb.GameContractTrxArgs{
		Type: tfcPb.GameTrxType_ROLL,
	}

	return ab
}

func (ab *ArgsBuilder) WithNextArgs() *ArgsBuilder {
	ab.trxArgs = &tfcPb.GameContractTrxArgs{
		Type: tfcPb.GameTrxType_NEXT,
	}

	return ab
}

func (ab *ArgsBuilder) WithTradeArgs(src, dest tfcPb.Player, r tfcPb.Resource, a int32) *ArgsBuilder {
	pLoad := &tfcPb.TradeTrxPayload{
		Source:   src,
		Dest:     dest,
		Resource: r,
		Amount:   a,
	}
	ab.trxArgs = &tfcPb.GameContractTrxArgs{
		Type:            tfcPb.GameTrxType_TRADE,
		TradeTrxPayload: pLoad,
	}

	return ab
}

func (ab *ArgsBuilder) WithBuildSettleArgs(player tfcPb.Player, sID uint32) *ArgsBuilder {
	settlePLoad := &tfcPb.BuildSettlePayload{
		Player:   player,
		SettleID: sID,
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

func (ab *ArgsBuilder) WithBuildRoadArgs(player tfcPb.Player, eID uint32) *ArgsBuilder {
	roadPLoad := &tfcPb.BuildRoadPayload{
		Player: player,
		EdgeID: eID,
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
