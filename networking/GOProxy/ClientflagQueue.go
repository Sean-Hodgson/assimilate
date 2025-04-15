package main

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"sync"
)

type IndividualClientFlag struct {
	clientLock               sync.Mutex
	clientHasSingletMessages bool
}

type ClientFlagArray struct {
	treeLock   sync.Mutex
	LookupTree *rbt.Tree
}

// flag handler function
func (q *ClientFlagArray) Setup() {
	q.treeLock.Lock()
	defer q.treeLock.Unlock()
	q.LookupTree = rbt.NewWithStringComparator()
}

// Function to insert data into the tree
func (q *ClientFlagArray) Insert(clientKey string, flag IndividualClientFlag) {
	// Lock the tree before modifying
	q.treeLock.Lock()
	defer q.treeLock.Unlock()

	// Insert the flag data with the clientKey as the key
	q.LookupTree.Put(clientKey, flag)
}

// Function to find data in the tree
func (q *ClientFlagArray) Find(clientKey string) (*IndividualClientFlag, bool) {
	// Lock the tree before reading
	q.treeLock.Lock()
	defer q.treeLock.Unlock()

	// Try to get the value from the tree
	value, found := q.LookupTree.Get(clientKey)
	if found {
		// Type assertion to get the value as IndividualClientFlag
		if flag, ok := value.(IndividualClientFlag); ok {
			return &flag, true
		}
	}
	return nil, false
}

//find in DS
