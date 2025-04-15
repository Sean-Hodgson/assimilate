package hash

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

// Get MAC Address
func GetMACAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range interfaces {
		if iface.HardwareAddr != nil {
			return iface.HardwareAddr.String()
		}
	}
	return ""
}

// Get CPU ID
func GetCPUID() string {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("wmic", "cpu", "get", "ProcessorId")
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("sh", "-c", "cat /proc/cpuinfo | grep -m 1 'Serial' | awk '{print $3}'")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
	} else {
		return "Unsupported OS"
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "Error fetching CPU ID"
	}

	return strings.TrimSpace(out.String())
}

// Get Disk Serial Number
func GetDiskSerial() string {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("wmic", "diskdrive", "get", "serialnumber")
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("lsblk", "-no", "SERIAL")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("diskutil", "info", "/")
	} else {
		return "Unsupported OS"
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "Error fetching Disk Serial"
	}

	// macOS specific
	if runtime.GOOS == "darwin" {
		// Extract serial number from output. diskutil throws out a lot of garbage so we have to find the specific line.
		output := out.String()
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Volume UUID:") {
				parts := strings.Fields(line) // Split the line into fields
				if len(parts) > 0 {
					return parts[len(parts)-1] // Return the serial number part
				}
			}
		}
	}

	return strings.TrimSpace(out.String())
}


func GenerateHardwareHash(cpuID, mac, diskSerial string) string {
	data := cpuID + mac + diskSerial
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}
