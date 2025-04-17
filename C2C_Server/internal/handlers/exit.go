package handlers

import (
	"fmt"
	"os"
	"time"
)

func ExitFunction() {
	exitStr := "Goodbye\n"

	for i := range len(exitStr) {
		fmt.Print(string(exitStr[i]))
		time.Sleep(100 * time.Millisecond)
	}

	os.Exit(0)
}
