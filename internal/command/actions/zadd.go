package actions

import (
	"strconv"

	"github.com/miladbarzideh/goldis/internal/datastore"
)

type ZAddCommand struct {
	dataStore *datastore.DataStore
}

func NewZAddCommand(dataStore *datastore.DataStore) *ZAddCommand {
	return &ZAddCommand{dataStore: dataStore}
}

func (c *ZAddCommand) Execute(args []string) string {
	if len(args) == 3 {
		score, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return SyntaxErrorMsg
		}
		return c.dataStore.ZAdd(args[0], score, args[2])
	}
	return SyntaxErrorMsg
}
