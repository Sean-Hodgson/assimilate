package main

import "sync"

func clientHandler(wg *sync.WaitGroup, recieveCNC chan AssimilateMessage) {
	println("clientHandler")
	defer wg.Done() //signal main that you are done when function ends


}
