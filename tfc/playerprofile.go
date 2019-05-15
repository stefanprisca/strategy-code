package tfc

import (
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

func GetPlayerId(player tfcPb.Player) int32 {
	return int32(player)
}

func GetResourceId(r tfcPb.Resource) int32 {
	return int32(r)
}

func InitPlayerProfile() *tfcPb.PlayerProfile {

	startingResources := make(map[int32]int32)
	for _, r := range []tfcPb.Resource{tfcPb.Resource_CAMP, tfcPb.Resource_FIELD, tfcPb.Resource_FOREST,
		tfcPb.Resource_MOUNTAIN, tfcPb.Resource_PASTURE, tfcPb.Resource_HILL} {
		id := GetResourceId(r)
		startingResources[id] = 5
	}

	return &tfcPb.PlayerProfile{
		Resources:     startingResources,
		WinningPoints: 0,
		Settlements:   2,
		Roads:         2,
	}
}
