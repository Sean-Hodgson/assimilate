package main

import (
	"bytes"
	"fmt"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {

	//setup variables
	//ProxyDSSetup := ProxyDSSetup()

	//start the CNC Hanlder
	go CNCHandler()

	//start the client handler
	go clientHandler()
}
