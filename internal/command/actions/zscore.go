package actions

import (
	"github.com/miladbarzideh/goldis/internal/datastore"
)

type ZScoreCommand struct {
	dataStore datastore.DataStore
}

func NewZScoreCommand(dataStore datastore.DataStore) *ZScoreCommand {
	return &ZScoreCommand{dataStore: dataStore}
}

func (c *ZScoreCommand) Execute(args []string) string {
	if len(args) == 2 {
		return c.dataStore.ZScore(args[0], args[1])
	}
	return SyntaxErrorMsg
}
