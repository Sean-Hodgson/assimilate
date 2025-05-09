package main

type AssimilateMessage struct {
	SenderHash string
	Command    string
	Arguments  []string
}
