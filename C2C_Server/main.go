package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	{"/exit", exitFunction},
	{"/whoami", whoami},
	{"/list", victimList},
	{"^C", victimList},
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
	fmt.Println("_______  _______  _______ _________ _______ _________ _        _______ _________ _______")
	fmt.Println("(  ___  )(  ____ \\(  ____ \\__   __/(       )\\__   __/( \\      (  ___  )\\__   __/(  ____ \\")
	fmt.Println("| (   ) || (    \\/| (    \\/   ) (   | () () |   ) (   | (      | (   ) |   ) (   | (    \\/")
	fmt.Println("| (___) || (_____ | (_____    | |   | || || |   | |   | |      | (___) |   | |   | (__    ")
	fmt.Println("|  ___  |(_____  )(_____  )   | |   | |(_)| |   | |   | |      |  ___  |   | |   |  __)   ")
	fmt.Println("| (   ) |      ) |      ) |   | |   | |   | |   | |   | |      | (   ) |   | |   | (      ")
	fmt.Println("| )   ( |/\\____) |/\\____) |___) (___| )   ( |___) (___| (____/\\| )   ( |   | |   | (____/\\")
	fmt.Println("|/     \\|\\_______)\\_______)\\_______/|/     \\|\\_______/(_______/|/     \\|   )_(   (_______/")

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
	fmt.Println("Ight imma head out ToT")
	os.Exit(0)
	//return 1
}

// Handler for menu
func Handler() {

	var command string
	fmt.Scanln(&command)
	for _, b := range options {
		if b.Cmd == command {
			b.Action()
			return
		}
	}

	fmt.Println("No command like that sorry uwu")

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
	displayMenu()

	//	Start of the menu func
	go startMenu()

	//	start of the server func
	go startServer()

	//	Makes it so the server and
	// 	menu run concurrently together
	select {}

}
