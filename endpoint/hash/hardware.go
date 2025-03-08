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

// Using this approach, we can map from an operating system string to a function
// rather than doing a if-else statement. Thus, adding support for a new operating
// system is as simple as adding an entry to the map - requiring no other modifications
// to existing functions.

var cpuidFuncs = map[string]func() string{
	"windows": func() string {
		cmd := exec.Command("wmic", "cpu", "get", "ProcessorId")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return "Error fetching CPU ID"
		}
		return strings.TrimSpace(out.String())
	},
	"linux": func() string {
		cmd := exec.Command("sh", "-c", "cat /proc/cpuinfo | grep -m 1 'Serial' | awk '{print $3}'")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return "Error fetching CPU ID"
		}
		return strings.TrimSpace(out.String())
	},
	"darwin": func() string {
		cmd := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return "Error fetching CPU ID"
		}
		return strings.TrimSpace(out.String())
	},
}

var diskSerialFuncs = map[string]func() string{
	"windows": func() string {
		cmd := exec.Command("wmic", "diskdrive", "get", "serialnumber")
		var out bytes.Buffer
		cmd.Stdout = &out

		if err := cmd.Run(); err != nil {
			return "Error fetching Disk Serial"
		}
		return strings.TrimSpace(out.String())
	},
	"linux": func() string {
		cmd := exec.Command("lsblk", "-no", "SERIAL")
		var out bytes.Buffer
		cmd.Stdout = &out

		if err := cmd.Run(); err != nil {
			return "Error fetching Disk Serial"
		}
		return strings.TrimSpace(out.String())
	},
	"darwin": func() string {
		cmd := exec.Command("diskutil", "info", "/")
		var out bytes.Buffer
		cmd.Stdout = &out

		if err := cmd.Run(); err != nil {
			return "Error fetching Disk Serial"
		}

		output := out.String()
		lines := strings.Split(output, "\n")

		for _, line := range lines {
			if strings.Contains(line, "Volume UUID:") {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					return parts[len(parts)-1]
				}
			}
		}

		return "Error parsing Disk Serial"
	},
}

func GetCPUID() string {
	if f, ok := cpuidFuncs[runtime.GOOS]; ok {
		return f()
	}

	return "Unsupported OS"
}

func GetDiskSerial() string {
	if f, ok := diskSerialFuncs[runtime.GOOS]; ok {
		return f()
	}

	return "Unsupported OS"
}

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

func GenerateHardwareHash(cpuID, mac, diskSerial string) string {
	data := cpuID + mac + diskSerial
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}
