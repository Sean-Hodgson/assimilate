package main

import (
	//"bufio"
	//"bytes"
	//"encoding/base64"
	//"encoding/json"
	"fmt"
	"io"
	"net/http"
	//"path/filepath"
	//"runtime"
	//"strings"
	//"github.com/olekukonko/tablewriter"
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

	//read the message sent
	message, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}

	//differentiate get vs post
	if r.Method == http.MethodGet {
		// Respond to a GET request
		fmt.Fprintln(w, "This is a GET request response.")
		fmt.Println("GET request received")
	} else if r.Method == http.MethodPost {
		println("this is a post request")
	}


	fmt.Println("Message: ", string(message))

}



