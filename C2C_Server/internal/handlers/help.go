package handlers

import (
	"assimilate/c2c_server/v2/internal/types"
	"fmt"
)

var Info = []types.CommandInfo{
	{
		Cmd:         "/help",
		Description: "Usage: /help",
	},
	{
		Cmd:         "/clear",
		Description: "Usage: /clear",
	},
	{
		Cmd:         "/exit",
		Description: "Usage: /exit",
	},
	{
		Cmd:         "/list",
		Description: "Usage: /list",
	},
	{
		Cmd:         "/floodping",
		Description: "Usage: /floodping <victimHash> <ip> <port>",
	},
	{
		Cmd:         "/supercat",
		Description: "Usage: /supercat <victimHash> <filePattern>",
	},
	{
		Cmd:         "/writefile",
		Description: "Usage: /writeFile <victimHash> <file1> [file2] ...",
	},
	{
		Cmd:         "/showsupercat",
		Description: "Usage: /showsupercat <victimHash>",
	},
}

func SpecificHelp() {
	fmt.Println("Available commands:")
	for _, opt := range Info {
		fmt.Printf("\t%s\n", opt.Cmd)
		fmt.Printf("\t\t%s\n", opt.Description)
	}
}
