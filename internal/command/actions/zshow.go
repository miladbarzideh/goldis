package actions

import (
	"github.com/miladbarzideh/goldis/internal/datastore"
)

type ZShowCommand struct {
	dataStore datastore.DataStore
}

func NewZShowCommand(dataStore datastore.DataStore) *ZShowCommand {
	return &ZShowCommand{dataStore: dataStore}
}

func (c *ZShowCommand) Execute(args []string) string {
	if len(args) == 1 {
		return c.dataStore.ZShow(args[0])
	}
	return SyntaxErrorMsg
}
