package main

import (
	"os"
	"os/exec"
	"time"
)


type option struct {
	Cmd        string
	ParamCount int //set the value to -1 if we do not care about the amount
	Action     func(params ...interface{})
	//in order to make this work, all functions that are in options have this as a parameter, but will be secure
	//because we are checking that the amount of params to pass equals the amount we set, and since we are always
	//passing in an array we must break that array apart in the function itself.
	//We set this to be params ...interface{} so it can also work with empty calls
}

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// runs in a goroutine
func startMenu() {
	for {
		Handler()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	//	Start of the menu func
	displayMenu() //display the menu before the start here
	go startMenu()

	//	start of the server func
	go startServer()

	
	select {}
}
