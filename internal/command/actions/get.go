package actions

import (
	"github.com/miladbarzideh/goldis/internal/datastore"
)

type GetCommand struct {
	dataStore *datastore.DataStore
}

func NewGetCommand(dataStore *datastore.DataStore) *GetCommand {
	return &GetCommand{dataStore: dataStore}
}

func (c *GetCommand) Execute(args []string) string {
	if len(args) == 1 {
		return c.dataStore.Get(args[0])
	}
	return SyntaxErrorMsg
}
