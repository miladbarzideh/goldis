package command

import (
	"log"
	"regexp"

	"github.com/miladbarzideh/goldis/internal/repository"
)

const syntaxErrorMsg = "(error) ERR syntax error"

type Handler struct {
	dataSource *repository.DataStore
}

func NewHandler() *Handler {
	return &Handler{dataSource: repository.NewDataStore()}
}

func (h *Handler) Execute(input string) string {
	commandParts := extractCommandParts(input)
	if commandParts == nil && len(commandParts) <= 1 {
		return syntaxErrorMsg
	}
	command, args := commandParts[0], commandParts[1:]
	log.Printf("Command %s will be executed", command)
	if command == "set" && len(args) == 2 {
		return h.dataSource.Set(args[0], args[1])
	} else if command == "get" && len(args) == 1 {
		return h.dataSource.Get(args[0])
	} else if command == "del" && len(args) == 1 {
		return h.dataSource.Delete(args[0])
	}
	return syntaxErrorMsg
}

func extractCommandParts(input string) []string {
	regex := regexp.MustCompile(`\w+`)
	return regex.FindAllString(input, -1)
}
