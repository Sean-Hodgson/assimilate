package handlers

import (
	"assimilate/c2c_server/v2/internal/model"
	"assimilate/c2c_server/v2/internal/types"
	"fmt"
)

func ListVictims(state *model.ServerState) types.ListVictimsHandlerFunc {
	return func() {
		state.ClientTreeMutex.Lock()
		defer state.ClientTreeMutex.Unlock()

		if state.ClientTree == nil || state.ClientTree.Size() == 0 {
			fmt.Println("No registered victims found.")
			return
		}

		it := state.ClientTree.Iterator()
		for it.Next() {
			hash := it.Key().(string)
			victim := it.Value().(types.VictimInfo)

			fmt.Printf("- Hash: %s\n", hash)

			if len(victim.Commands) == 0 {
				fmt.Println("  No queued commands.")
			} else {
				for i, cmd := range victim.Commands {
					fmt.Printf("  [%d] CommandType: %s, Args: %v\n",
						i+1, cmd.Commandytype, cmd.Arguments)
				}
			}

			fmt.Println() // spacing between victims
		}
	}
}
