package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {

	//setup variables

	// CNCConf := CNCConf()
	// ClientConf := ClientConf()

	//start the CNC Hanlder
	go CNCHanlder()

	//start the client handler
	go clientHanlder()
}
