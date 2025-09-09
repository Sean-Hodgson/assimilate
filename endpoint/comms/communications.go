package communications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mal/hash"
	"net/http"
	"runtime"
)

const RegistrationURL = "https://example.com/register"

type RegistrationData struct {
	HardwareHash string `json:"hardwarehash"`
	OS           string `json:"OS"`
}

func Register() {
	// Gather hardware information
	cpuID := hash.GetCPUID()
	mac := hash.GetMACAddress()
	diskSerial := hash.GetDiskSerial()

	hardwareHash := hash.GenerateHardwareHash(cpuID, mac, diskSerial)

	data := RegistrationData{
		HardwareHash: hardwareHash,
		OS:           runtime.GOOS,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	resp, err := http.Post(RegistrationURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Registration successful, status:", resp.Status)
	} else {
		fmt.Println("Registration failed, status:", resp.Status)
	}
}

func Heartbeat() {
	//TODO:
}
