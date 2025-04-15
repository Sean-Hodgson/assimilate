package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/net/proxy"
)


var (
	cncPort      = flag.Int("socks5port", 9050, "port for the local socks5 proxy")
	servingPort  = flag.Int("ingressport", 5555, "ingress port for visitors")
	remotePort   = flag.Int("remoteport", 7777, "remote port for the .onion service")
	onionDomain  = flag.String("onion", "azgbvwd2j47yqsli5tdwbhms2cmswttdmj7xk5bwipoqmstwc33j63yd.onion", "onion domain name to access")
)


func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	println("successfully received message")
	defer r.Body.Close()

	proxyReq, err := http.NewRequest(r.Method, (*onionDomain+r.RequestURI), r.Body)
	if(err != nil){
		println("ERROR:ProxyHandler failed to get http request")
		return
	}


	// SOCKS5 proxy setup
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		log.Println("ERROR: Failed to create SOCKS5 proxy dialer:", err)
		http.Error(w, "Failed to create SOCKS5 proxy dialer", http.StatusInternalServerError)
		return
	}

	// Create an HTTP client with the SOCKS5 proxy dialer
	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}

	println("got message: ", r.Body)
	
	resp, err := client.Do(proxyReq)
	if(err != nil){
		println("ERROR:proxyHandler failed to send proxied messages")
		return
	}
	defer resp.Body.Close()

	println("got cnc response: ", resp.Body)


	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	w.Write(bodyBytes)
}



func main() {
	println("Starting main")
	flag.Parse()
	*onionDomain = (*onionDomain+ ":" + strconv.Itoa(*remotePort))

	// Parse the flags
	
	http.HandleFunc("/proxy", ProxyHandler)
	fmt.Println("Serving on port", *servingPort, "going towards socks5 localproxy:127.0.0.1:", *cncPort,  "sending to .onion address: ", *onionDomain)
	
	err := http.ListenAndServe((":" + strconv.Itoa(*servingPort)), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}


