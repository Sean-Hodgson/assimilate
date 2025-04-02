package main

import "sync"

func clientHandler(wg *sync.WaitGroup, recieveCNC chan AssimilateMessage) {
	println("clientHandler")

}
