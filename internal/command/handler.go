package command

import (
	"log"
	"regexp"
	"strings"

	"github.com/miladbarzideh/goldis/internal/datastore"
)

const (
	getCommand     = "get"
	setCommand     = "set"
	delCommand     = "del"
	syntaxErrorMsg = "(error) ERR syntax error"
)

type Handler struct {
	dataSource *datastore.DataStore
}

func NewHandler() *Handler {
	return &Handler{dataSource: datastore.NewDataStore()}
}

func (h *Handler) Execute(input string) string {
	commandParts := extractCommandParts(input)
	if commandParts == nil || len(commandParts) <= 1 {
		return syntaxErrorMsg
	}
	command, args := commandParts[0], commandParts[1:]
	log.Printf("Command %s will be executed", command)
	switch {
	case strings.EqualFold(command, setCommand) && len(args) == 2:
		return h.dataSource.Set(args[0], args[1])
	case strings.EqualFold(command, getCommand) && len(args) == 1:
		return h.dataSource.Get(args[0])
	case strings.EqualFold(command, delCommand) && len(args) == 1:
		return h.dataSource.Delete(args[0])
	}
	return syntaxErrorMsg
}

func extractCommandParts(input string) []string {
	regex := regexp.MustCompile(`\w+`)
	return regex.FindAllString(input, -1)
}
