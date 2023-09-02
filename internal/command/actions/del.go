package actions

import (
	"github.com/miladbarzideh/goldis/internal/datastore"
)

type DelCommand struct {
	dataStore datastore.DataStore
}

func NewDelCommand(dataStore datastore.DataStore) *DelCommand {
	return &DelCommand{dataStore: dataStore}
}

func (c *DelCommand) Execute(args []string) string {
	if len(args) == 1 {
		return c.dataStore.Delete(args[0])
	}
	return SyntaxErrorMsg
}
