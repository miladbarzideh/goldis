package actions

import (
	"github.com/miladbarzideh/goldis/internal/datastore"
)

type ZRemCommand struct {
	dataStore datastore.DataStore
}

func NewZRemCommand(dataStore datastore.DataStore) *ZRemCommand {
	return &ZRemCommand{dataStore: dataStore}
}

func (c *ZRemCommand) Execute(args []string) string {
	if len(args) == 2 {
		return c.dataStore.ZRemove(args[0], args[1])
	}
	return SyntaxErrorMsg
}
