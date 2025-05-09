package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

// setup struct for the thread
type CNCFuncSetup struct {
	timeout    time.Duration
	retryCount int
	retryDelay time.Duration
	port       int
	targetHost string
}

func CNCHandler(wg *sync.WaitGroup, receiveClient chan AssimilateMessage, setup CNCFuncSetup) {
	defer wg.Done() //signal main that you are done when function ends
	println("Entered CNCHandler")

	//setup the network IO
	readAccumulator := new(bytes.Buffer)            //queue structure for the datastream
	readBuffer := make([]byte, 1024)                //input buffer
	conn, connStatus := AssimilateTCPConnect(setup) //setup comms with server
	if connStatus == false {                        //if connection failed kill thread otherwise defer the close
		println("CNCHANDLER:Error connecting to server, terminating thread")
		return
	}
	defer conn.Close()

	//main connection loop
	for {

		//make message to send
		tosend := AssimilateMessage{
			SenderHash: "testing",
			Command:    "send",
			Arguments:  []string{setup.targetHost},
		}

		//convert to JSON
		jsonData, err := json.Marshal(tosend)

		/////IMPORTANT CAP ALL MESSAGES WITH "\r\n"/////
		finalData := append(jsonData, byte('\r'))
		finalData = append(finalData, byte('\n'))

		if err != nil {
			fmt.Println("CNCHANDLER:Error encoding to JSON:", err)
			return
		}

		for {
			//write stuff out
			conn.Write(finalData)
			time.Sleep(5 * time.Second)

			//read stuff in
			conn.SetReadDeadline(time.Now().Add(2 * time.Millisecond))
			bytesRead, readerr := conn.Read(readBuffer)
			if readerr == nil && bytesRead > 0 {

				println("attempting to read bytes length of %d", bytesRead)
				//only read stuff in if it get actual bytes
				readAccumulator.Write(readBuffer)
				//handle data
				println(string(readBuffer[:bytesRead]))
			}
		}
	}

	return
}

func AssimilateTCPConnect(setup CNCFuncSetup) (net.Conn, bool) {

	retryCounter := setup.retryCount

	for {
		//attempt to connect
		fmt.Printf("Attempting to connect to %s...\n", setup.targetHost)
		mergedTargetHost := setup.targetHost + ":" + fmt.Sprintf("%d", setup.port)
		conn, err := net.DialTimeout("tcp", mergedTargetHost, setup.timeout*time.Second)

		//if conn returned nil this is a failure wait and retry in a bit
		if err != nil {
			fmt.Printf("CNCHANDLER:Failed to connect: %v\n", err)
			fmt.Printf("CNCHANDLER:Retrying in %s...\n", setup.retryDelay)
			time.Sleep(setup.retryDelay)
			retryCounter--

			//if counter is out of retries kill program
			if retryCounter <= 0 {
				return nil, false
			}

			continue // Go to the next iteration of the loop to retry
		}

		println("CNCHANDLER:Successfully connected to server")
		return conn, true
	}
}
