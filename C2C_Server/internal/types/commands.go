package types

// could also split all this into files, if wanted to
type VictimInfo struct {
	Hash     string
	Commands []VictimCommand
}

type IncomingVictimRequest struct {
	RequesterHash string //who is making this request
}

type VictimCommand struct {
	Commandytype string
	Arguments    []string
}

type IncomingSuperCatData struct {
	Filename string
	Content  string
	Hash     string
}
