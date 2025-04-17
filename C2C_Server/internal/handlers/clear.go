package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func CallClear() {
	clear := make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	if fn, ok := clear[runtime.GOOS]; ok {
		fn()
	} else {
		fmt.Println("Cannot clear screen on this operating system.")
	}
}
