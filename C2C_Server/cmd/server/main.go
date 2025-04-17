package main

import (
	"assimilate/c2c_server/v2/internal/menu"
	"assimilate/c2c_server/v2/internal/model"
	"assimilate/c2c_server/v2/internal/server"
	"sync"
	"time"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

func main() {
	initialServerState := model.ServerState{
		ClientTree:        rbt.NewWithStringComparator(),
		SuperCatTree:      rbt.NewWithStringComparator(),
		ClientTreeMutex:   sync.Mutex{},
		SuperCatTreeMutex: sync.Mutex{},
		GlobalDebug:       true,
	}

	go startMenu(&initialServerState)

	go server.StartServer(&initialServerState)

	select {}
}

func startMenu(state *model.ServerState) {
	for {
		menu.Handler(state)
		time.Sleep(500 * time.Millisecond)
	}
}
