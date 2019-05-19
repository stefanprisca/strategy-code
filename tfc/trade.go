package tfc

import (
	"fmt"
	"regexp"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

func handleTrade(APIstub shim.ChaincodeStubInterface, creatorSign []byte,
	gameData tfcPb.GameData, payload tfcPb.TradeTrxPayload) (tfcPb.GameData, error) {

	err := assertTradePrecond(gameData, creatorSign, payload)
	if err != nil {
		return gameData, fmt.Errorf(
			"trade preconditions not met: %s", err)
	}

	srcID := GetPlayerId(payload.Source)
	destID := GetPlayerId(payload.Dest)
	resID := GetResourceId(payload.Resource)

	srcProfile := gameData.Profiles[srcID]
	destProfile := gameData.Profiles[destID]

	srcProfile.Resources[resID] -= payload.Amount
	destProfile.Resources[resID] += payload.Amount

	return gameData, assertTradePostcond(gameData, payload)
}

func stateValidationString(src tfcPb.Player, state tfcPb.GameState) string {
	return fmt.Sprintf("%v%v", src, state)
}

var tradeStateValidationRegexp = regexp.MustCompile(
	fmt.Sprintf("%vRTRADE", tfcPb.Player_RED) +
		fmt.Sprintf("|%vBTRADE", tfcPb.Player_BLUE) +
		fmt.Sprintf("|%vGTRADE", tfcPb.Player_GREEN))

// TODO: implement this
func assertTradePrecond(gameData tfcPb.GameData, creatorSign []byte, payload tfcPb.TradeTrxPayload) error {
	state := gameData.State
	creator, err := getCreator(gameData, creatorSign)
	if err != nil || creator != payload.Source {
		return fmt.Errorf("invalid trx creator, or creator not identified (<0): expected %v, got %v",
			creator, payload.Source)
	}

	stateValidationStr := stateValidationString(creator, state)
	if !tradeStateValidationRegexp.MatchString(stateValidationStr) {
		return fmt.Errorf("expected state to match one of %v, got %v",
			tradeStateValidationRegexp, stateValidationStr)
	}

	return nil
}

func assertTradePostcond(gameData tfcPb.GameData, payload tfcPb.TradeTrxPayload) error {
	res := payload.Resource
	if err := hasValidPostTradeAmount(payload.Source, gameData, res); err != nil {
		return err
	}

	if err := hasValidPostTradeAmount(payload.Dest, gameData, res); err != nil {
		return err
	}
	return nil
}

func hasValidPostTradeAmount(p tfcPb.Player, gameData tfcPb.GameData, r tfcPb.Resource) error {
	rID := GetResourceId(r)
	pP := *gameData.Profiles[GetPlayerId(p)]
	available := pP.Resources[rID]
	if available < 0 {
		return fmt.Errorf("player %v does not have required %v resources: %s",
			p, r,
			fmt.Sprintf("available: %v", available))
	}
	return nil
}
