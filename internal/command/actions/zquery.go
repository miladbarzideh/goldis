package actions

import (
	"strconv"

	"github.com/miladbarzideh/goldis/internal/datastore"
)

type ZQueryCommand struct {
	dataStore *datastore.DataStore
}

func NewZQueryCommand(dataStore *datastore.DataStore) *ZQueryCommand {
	return &ZQueryCommand{dataStore: dataStore}
}

func (c *ZQueryCommand) Execute(args []string) string {
	if len(args) == 5 {
		score, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return SyntaxErrorMsg
		}
		offset, err := strconv.Atoi(args[3])
		if err != nil {
			return SyntaxErrorMsg
		}
		limit, err := strconv.Atoi(args[4])
		if err != nil {
			return SyntaxErrorMsg
		}
		return c.dataStore.ZQuery(args[0], score, args[2], int32(uint32(offset)), uint32(limit))
	}
	return SyntaxErrorMsg
}
