package actions

import "github.com/miladbarzideh/goldis/internal/datastore"

type KeysCommand struct {
	dataStore datastore.DataStore
}

func NewKeysCommand(dataStore datastore.DataStore) *KeysCommand {
	return &KeysCommand{dataStore: dataStore}
}

func (c *KeysCommand) Execute(args []string) string {
	return c.dataStore.Keys()
}
