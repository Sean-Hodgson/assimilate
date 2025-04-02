package main

import (
	"sync"
	"time"
)

type CNCFuncSetup struct {
	timeout    time.Duration
	retryCount int
	retryDelay time.Duration
	port       int
	targetHost string
}

type ClientFuncSetup struct {
	timeout time.Duration
	port    int
}

func main() {
	println("Starting main")

	//setup thread waiting
	var wg sync.WaitGroup
	wg.Add(1) // Increment the counter for one goroutines

	//recieveCNC := make(chan AssimilateMessage)
	receiveClient := make(chan AssimilateMessage)

	//setup overall DS
	//ProxyDSSetup := ProxyDSSetup()

	//start the CNC Handler
	inputParamsCNC := CNCFuncSetup{
		timeout:    10000,
		retryCount: 5,
		retryDelay: 10 * time.Millisecond,
		port:       7884,
		targetHost: "127.0.0.1",
	}
	go CNCHandler(&wg, receiveClient, inputParamsCNC)

	wg.Wait()
	//start the client handler
	//go clientHandler(&wg, receiveCNC)
}

/*
package main

import (
	"fmt"
	"time"
)

func sender(ch chan string, message string) {
	fmt.Printf("Sender goroutine attempting to send: %q\n", message)
	ch <- message // Send the message to the channel
	fmt.Printf("Sender goroutine sent: %q\n", message)
}

func receiver(ch chan string) {
	fmt.Println("Receiver goroutine is waiting to receive.")
	receivedMessage := <-ch // Receive a message from the channel
	fmt.Printf("Receiver goroutine received: %q\n", receivedMessage)
}

func main() {
	// Create an unbuffered channel that can hold strings [see our conversation history]
	messageChannel := make(chan string)

	// Start the sender goroutine, passing the channel and the message [see our conversation history]
	go sender(messageChannel, "Hello from sender!")

	// Start the receiver goroutine, passing the same channel [see our conversation history]
	go receiver(messageChannel)

	// Keep the main goroutine alive for a short time to allow the others to execute
	time.Sleep(1 * time.Second)

	fmt.Println("Main goroutine finished.")
}*/
