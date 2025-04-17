package menu

import (
	"assimilate/c2c_server/v2/internal/model"
	"assimilate/c2c_server/v2/internal/options"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Handler(state *model.ServerState) {
	displayMenu()

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	parts := strings.Split(input, " ")
	if len(parts) == 0 {
		return
	}
	command := parts[0]
	args := parts[1:]

	for _, opt := range options.Options(state) {
		if opt.Cmd == command {
			opt.Action.Execute(args)
			return
		}
	}

	fmt.Println("Invalid Command, use the /help command for details of all the commands")
}
