package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func CNCHandler(wg *sync.WaitGroup, receiveClient chan AssimilateMessage, setup CNCFuncSetup) {
	defer wg.Done() //signal main that you are done when function ends
	println("Entered CNCHandler")

	readAccumulator := new(bytes.Buffer)
	readBuffer := make([]byte, 1024)

	for {
		conn, connStatus := AssimilateTCPConnect(setup)
		if connStatus {
			//if the connection has failed return and stop
			return
		}

		//make message to send
		tosend := AssimilateMessage{
			SenderHash: "testing",
			Command:    "send",
			Arguments:  []string{setup.targetHost},
		}

		//convert to JSON
		jsonData, err := json.Marshal(tosend)
		if err != nil {
			fmt.Println("CNCHANDLER:Error encoding to JSON:", err)
			return
		}

		for {

			conn.Write(jsonData)
			time.Sleep(5 * time.Second)
			conn.Read(readBuffer)
			println(readBuffer)
		}
	}

	return
}

func AsimilateTCPConnect(setup CNCFuncSetup) bool {

	retryCounter := setup.retryCount

	for {
		fmt.Printf("Attempting to connect to %s...\n", setup.targetHost)
		mergedTargetHost := setup.targetHost + ":" + fmt.Sprintf("%d", setup.port)
		conn, err := net.DialTimeout("tcp", mergedTargetHost, setup.timeout*time.Second)

		if retryCounter <= 0 {
			return false
		}

		if err != nil {
			fmt.Printf("CNCHANDLER:Failed to connect: %v\n", err)
			fmt.Printf("CNCHANDLER:Retrying in %s...\n", setup.retryDelay)
			time.Sleep(setup.retryDelay)
			continue // Go to the next iteration of the loop to retry
		}

		err = conn.SetReadDeadline(time.Now().Add(2 * time.Nanosecond))
		if err != nil {
			log.Println("Error setting read deadline:", err)
			continue
		}

		return true
	}
}
