package actions

import (
	"github.com/miladbarzideh/goldis/internal/datastore"
)

type SetCommand struct {
	dataStore datastore.DataStore
}

func NewSetCommand(dataStore datastore.DataStore) *SetCommand {
	return &SetCommand{dataStore: dataStore}
}

func (c *SetCommand) Execute(args []string) string {
	if len(args) == 2 {
		return c.dataStore.Set(args[0], args[1])
	}
	return SyntaxErrorMsg
}
