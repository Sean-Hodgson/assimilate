package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// .onion service
	var cncServerAddr = "http://127.0.0.1:7777"

	//port must be the local port for the tor service
	var cncPort = 7777

	//serving Port
	var servingPort = 5555

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	proxyReq, err := http.NewRequest(r.Method, (cncServerAddr+r.RequestURI), r.Body)
	if(err != nil){
		println("ERROR:ProxyHandler failed to get http request")
		return
	}
	
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if(err != nil){
		println("ERROR:proxyHandler failed to send proxied messages")
		return
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	w.Write(bodyBytes)

}

func main() {
	println("Starting main")

	http.HandleFunc("/proxy", ProxyHandler)
	fmt.Println("Serving on port", servingPort, "sending to .onion serving: ", cncServerAddr)
	
	err := http.ListenAndServe((":" + strconv.Itoa(servingPort)), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}


