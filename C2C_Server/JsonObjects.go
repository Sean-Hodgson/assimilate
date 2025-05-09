package main

import "time"

/*
case "response/json"       paired with ResponseCode
case "victimregister/json" paired with VictimRegister
case "requestCommand/json" paired with IncomingVictimRequest struct
case "incomingData/json"   paired with incoming superCat data
*/


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
	OutputHash 	 	 string
	GenerationSource string
	Data 		 	 string
}


