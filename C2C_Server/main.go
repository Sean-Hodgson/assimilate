package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olekukonko/tablewriter"
)

// Edit later
type victim struct {
	ID   string `json:"id"`
	Ip   string `json:"ip"`
	User string `json:"user"`
	Port string `json:"port"`
}

// Sample data
var victims = []victim{
	{ID: "1", Ip: "127.23.44", User: "Fred", Port: ""},
	{ID: "2", Ip: "127.23.43", User: "bro", Port: "1988"},
}

type option struct {
	Cmd    string
	Action func()
}

// edit later
var options = []option{
	{"/help", specificHelp},
	{"/clear", CallClear},
	{"/exit", exitFunction},
	{"/whoami", whoami},
	{"/list", victimList},
	{"/floodping", victimList},
	{"/supercat", victimList},
	{"/writefile", victimList},
}

// lists victims
func getVictims(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, victims)

}

// makes a new veictim
func createVictim(c *gin.Context) {
	var newVictim victim

	if err := c.BindJSON(&newVictim); err != nil {
		return
	}

	victims = append(victims, newVictim)
	c.IndentedJSON(http.StatusCreated, newVictim)

}

func getVictimById(id string) (*victim, error) {
	for i, b := range victims {
		if b.ID == id {
			return &victims[i], nil
		}
	}

	return nil, errors.New("not found")
}

// Searches for victims in list
func VictimById(c *gin.Context) {
	id := c.Param("id")
	victim, err := getVictimById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "victim not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, victim)
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

// prints when an invalid command is entered
func genericHelp() {
	fmt.Println("Invalid Command, use the /help command for details of all the commands")
}

func specificHelp() {
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
	table.Append([]string{"/whoami", "None", "Returns your device details"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/list", "None", "Lists all victims currently infected"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/floodping", "TBD", "does something i think"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/supercat", "TBD", "something else haha"})
	table.Append([]string{"", "", ""})
	table.Append([]string{"/writefile", "TBD", "hate this command"})

	// Render the table
	table.Render()
}

// just filler command
func whoami() {
	fmt.Println("This is just filler, but I'm Steve from minecraft")
}

// Displays vicitim list
func victimList() {
	for _, v := range victims {
		fmt.Printf("ID: %s | IP: %s | User: %s | Port: %s\n", v.ID, v.Ip, v.User, v.Port)
	}
}

// exits the program
func exitFunction() {
	exitStr := "Goodbye\n"
	for i := 0; i < len(exitStr); i++ {
		fmt.Print(string(exitStr[i]))
		time.Sleep(150 * time.Millisecond)
	}
	os.Exit(0)
	//return 1
}

// Handler for menu
func Handler() {

	var command string
	fmt.Print(("(C2C)> "))
	fmt.Scanln(&command)
	for _, b := range options {
		if b.Cmd == command {
			b.Action()
			return
		}
	}

	//provide the user a default help for entering commands
	genericHelp()

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

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(\n")
	}
}

//	============================
// 	Use goCurl? or other options
//	Use curl to get information from victims?
//	============================

//	==========================

// runs in a goroutine
func startServer() {
	router := gin.Default()
	router.GET("/victims", getVictims)
	router.GET("/victims/:id", VictimById)
	router.POST("/createvictim", createVictim)

	fmt.Println("Running on localhost:8888")
	router.Run("localhost:8888")
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
