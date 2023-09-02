package actions

import (
	"github.com/miladbarzideh/goldis/internal/datastore"
)

type TTLCommand struct {
	dataStore *datastore.DataStore
}

func NewTTLCommand(dataStore *datastore.DataStore) *TTLCommand {
	return &TTLCommand{dataStore: dataStore}
}

func (c *TTLCommand) Execute(args []string) string {
	if len(args) == 1 {
		return c.dataStore.Ttl(args[0])
	}
	return SyntaxErrorMsg
}
