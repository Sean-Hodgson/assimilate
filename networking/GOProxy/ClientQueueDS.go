package main

import (
	"sync"
)

type IndividualClientDataQueue struct {
	clientLock       sync.Mutex
	queueIncrementer int
	clientHash       string
	messageStack     []string
}

//receiver functions

func (q *IndividualClientDataQueue) Setup(hash string) {
	q.clientLock.Lock()
	defer q.clientLock.Unlock()

	q.queueIncrementer = 0
	q.clientHash = hash
	q.messageStack = make([]string, 0)
}

func (q *IndividualClientDataQueue) Enqueue(givenString string) {
	q.clientLock.Lock()
	defer q.clientLock.Unlock()

	q.messageStack = append(q.messageStack, givenString) // Add new command to the queue
	q.queueIncrementer++
}

func (q *IndividualClientDataQueue) Dequeue() interface {
}
