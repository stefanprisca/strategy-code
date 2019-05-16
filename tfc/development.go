package tfc

import (
	"fmt"
	"regexp"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

func handleDev(APIstub shim.ChaincodeStubInterface, creatorSign []byte,
	gameData tfcPb.GameData, payload tfcPb.BuildTrxPayload) (tfcPb.GameData, error) {

	err := assertDevelopmentPrecond(gameData, creatorSign, payload)
	if err != nil {
		return gameData, fmt.Errorf(
			"development preconditions not met: %s", err)
	}

	switch payload.Type {
	case tfcPb.BuildType_SETTLE:
		return buildSettlement(gameData, *payload.BuildSettlePayload)
	case tfcPb.BuildType_ROAD:
		return buildRoad(gameData, *payload.BuildRoadPayload)
	}

	return gameData, nil
}

var devStateValidationRegexp = regexp.MustCompile(
	fmt.Sprintf("%vRDEV", tfcPb.Player_RED) +
		fmt.Sprintf("|%vBDEV", tfcPb.Player_BLUE) +
		fmt.Sprintf("|%vGDEV", tfcPb.Player_GREEN))

// TODO: implement this
func assertDevelopmentPrecond(gameData tfcPb.GameData, creatorSign []byte, payload tfcPb.BuildTrxPayload) error {

	/*
		1) correct state
		2.e) empty and an edge point is one of the player settlements
		2.s) empty and neighbouring intersections for distance rule: all surrounding intersections are free
	*/

	state := gameData.State
	creator, err := getCreator(gameData, creatorSign)
	if err != nil {
		return err
	}

	stateValidationStr := stateValidationString(creator, state)
	if !devStateValidationRegexp.MatchString(stateValidationStr) {
		return fmt.Errorf("expected state to match one of %v, got %v",
			devStateValidationRegexp, stateValidationStr)
	}

	switch payload.Type {
	case tfcPb.BuildType_ROAD:
		return assertBuildRoadPrecond(gameData, creator, *payload.BuildRoadPayload)
	case tfcPb.BuildType_SETTLE:
		return assertBuildSettlePrecond(gameData, creator, *payload.BuildSettlePayload)
	}

	return fmt.Errorf("Unkown build type %v", payload.Type)
}

var canBuildRoad = regexp.MustCompile(
	fmt.Sprintf("(%s)|(%s)|(%s)",
		fmt.Sprintf("%v%v.*(%v|%v)+.*",
			tfcPb.Player_RED, tfcPb.Road_NOROAD,
			tfcPb.Road_REDROAD, tfcPb.Settlement_REDSETTLE),
		fmt.Sprintf("%v%v.*(%v|%v)+.*",
			tfcPb.Player_GREEN, tfcPb.Road_NOROAD,
			tfcPb.Road_GREENROAD, tfcPb.Settlement_GREENSETTLE),
		fmt.Sprintf("%v%v.*(%v|%v)+.*",
			tfcPb.Player_BLUE, tfcPb.Road_NOROAD,
			tfcPb.Road_BLUEROAD, tfcPb.Settlement_BLUESETTLE)))

func buildRoadValidString(p tfcPb.Player, r tfcPb.Road, s1, s2 tfcPb.Settlement, r1, r2, r3, r4 tfcPb.Road) string {
	return fmt.Sprintf("%v%v%v%v%v%v%v%v", p, r, s1, s2, r1, r2, r3, r4)
}

func assertBuildRoadPrecond(gameData tfcPb.GameData, creator tfcPb.Player, payload tfcPb.BuildRoadPayload) error {
	if creator != payload.Player {
		return fmt.Errorf("expected creator to match trx player. expected %v, got %v", creator, payload.Player)
	}

	eID := uint32(payload.EdgeID)
	E, exists := gameData.Board.Edges[eID]
	if !exists {
		return fmt.Errorf("gameboard edge %v does not exist", eID)
	}

	r := E.Attributes.Road
	nE := gameData.Board.Edges[E.Next]

	s1 := gameData.Board.Intersections[E.Origin].Attributes.Settlement
	s2 := gameData.Board.Intersections[nE.Origin].Attributes.Settlement

	r1 := nE.Attributes.Road
	r2 := gameData.Board.Edges[E.Prev].Attributes.Road

	twin := gameData.Board.Edges[E.Twin]
	r3 := gameData.Board.Edges[twin.Prev].Attributes.Road
	r4 := gameData.Board.Edges[twin.Next].Attributes.Road

	validationS := buildRoadValidString(creator, r, s1, s2, r1, r2, r3, r4)

	if !canBuildRoad.MatchString(validationS) {
		return fmt.Errorf("could not build road for player %v, conditions not fulfilled: %s", creator,
			fmt.Sprintf("existing road %v; surrounding settlements [%v, %v]; surrounding edges [%v, %v, %v, %v]",
				r, s1, s2, r1, r2, r3, r4))
	}

	return nil
}

var canBuildSettle = regexp.MustCompile(
	fmt.Sprintf("%v%v%v%v",
		tfcPb.Settlement_NOSETTLE, tfcPb.Settlement_NOSETTLE,
		tfcPb.Settlement_NOSETTLE, tfcPb.Settlement_NOSETTLE))

func buildSettleValidString(s, s1, s2, s3 tfcPb.Settlement) string {
	return fmt.Sprintf("%v%v%v%v", s, s1, s2, s3)
}

func assertBuildSettlePrecond(gameData tfcPb.GameData, creator tfcPb.Player, payload tfcPb.BuildSettlePayload) error {
	if creator != payload.Player {
		return fmt.Errorf("expected creator to match trx player. expected %v, got %v", creator, payload.Player)
	}

	sID := uint32(payload.SettleID)
	I := gameData.Board.Intersections[sID]
	s := I.Attributes.Settlement

	iE := gameData.Board.Edges[I.IncidentEdge]
	iID1 := gameData.Board.Edges[iE.Next].Origin
	s1 := gameData.Board.Intersections[iID1].Attributes.Settlement

	iID2 := gameData.Board.Edges[iE.Prev].Origin
	s2 := gameData.Board.Intersections[iID2].Attributes.Settlement

	iETwin := gameData.Board.Edges[iE.Twin]
	iETwinNext := gameData.Board.Edges[iETwin.Next]
	iID3 := gameData.Board.Edges[iETwinNext.Next].Origin
	s3 := gameData.Board.Intersections[iID3].Attributes.Settlement

	validationS := buildSettleValidString(s, s1, s2, s3)

	if !canBuildSettle.MatchString(validationS) {
		return fmt.Errorf("could not build settlement for player %v, conditions not fulfilled: %s", creator,
			fmt.Sprintf("existing settle %v; surrounding settlements [%v, %v, %v]: %s", s, s1, s2, s3,
				fmt.Sprintf("expected to match %v, got %v", canBuildSettle, validationS)))
	}

	return nil
}

func buildRoad(
	gameData tfcPb.GameData, payload tfcPb.BuildRoadPayload) (tfcPb.GameData, error) {

	playerID := GetPlayerId(payload.Player)
	profile := gameData.Profiles[playerID]

	profile.Resources[GetResourceId(tfcPb.Resource_HILL)]--
	profile.Resources[GetResourceId(tfcPb.Resource_FOREST)]--
	profile.Roads--
	profile.WinningPoints++

	eID := uint32(payload.EdgeID)
	edge := gameData.Board.Edges[eID]
	switch payload.Player {
	case tfcPb.Player_RED:
		edge.Attributes.Road = tfcPb.Road_REDROAD
	case tfcPb.Player_GREEN:
		edge.Attributes.Road = tfcPb.Road_GREENROAD
	case tfcPb.Player_BLUE:
		edge.Attributes.Road = tfcPb.Road_BLUEROAD
	}

	return gameData, nil
}

func buildSettlement(
	gameData tfcPb.GameData, payload tfcPb.BuildSettlePayload) (tfcPb.GameData, error) {

	playerID := GetPlayerId(payload.Player)
	profile := gameData.Profiles[playerID]

	for rID := range profile.Resources {
		profile.Resources[rID]--
	}

	profile.Settlements--
	profile.WinningPoints += 2

	posID := uint32(payload.SettleID)
	settleIntersection := gameData.Board.Intersections[posID]
	switch payload.Player {
	case tfcPb.Player_RED:
		settleIntersection.Attributes.Settlement = tfcPb.Settlement_REDSETTLE
	case tfcPb.Player_GREEN:
		settleIntersection.Attributes.Settlement = tfcPb.Settlement_GREENSETTLE
	case tfcPb.Player_BLUE:
		settleIntersection.Attributes.Settlement = tfcPb.Settlement_BLUESETTLE
	}

	return gameData, nil
}
