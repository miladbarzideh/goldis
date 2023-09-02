package actions

import (
	"strconv"

	"github.com/miladbarzideh/goldis/internal/datastore"
)

type ExpireCommand struct {
	dataStore *datastore.DataStore
}

func NewExpireCommand(dataStore *datastore.DataStore) *ExpireCommand {
	return &ExpireCommand{dataStore: dataStore}
}

func (c *ExpireCommand) Execute(args []string) string {
	if len(args) == 2 {
		ttl, err := strconv.Atoi(args[1])
		if err != nil {
			return SyntaxErrorMsg
		}
		return c.dataStore.Expire(args[0], int64(ttl))
	}
	return SyntaxErrorMsg
}
