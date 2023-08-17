package command

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/miladbarzideh/goldis/internal/datastore"
)

const (
	getCommand     = "get"
	setCommand     = "set"
	delCommand     = "del"
	zaddCommand    = "zadd"
	zremCommand    = "zrem"
	zscoreCommand  = "zscore"
	zqueryCommand  = "zquery"
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
	case strings.EqualFold(command, zaddCommand) && len(args) == 3:
		score, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return syntaxErrorMsg
		}
		return h.dataSource.ZAdd(args[0], score, args[2])
	case strings.EqualFold(command, zremCommand) && len(args) == 2:
		return h.dataSource.ZRemove(args[0], args[1])
	case strings.EqualFold(command, zscoreCommand) && len(args) == 2:
		return h.dataSource.ZScore(args[0], args[1])
	case strings.EqualFold(command, zqueryCommand) && len(args) == 5:
		score, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return syntaxErrorMsg
		}
		offset, err := strconv.Atoi(args[3])
		if err != nil {
			return syntaxErrorMsg
		}
		limit, err := strconv.Atoi(args[4])
		if err != nil {
			return syntaxErrorMsg
		}
		return h.dataSource.ZQuery(args[0], score, args[2], uint32(offset), uint32(limit))

	}
	return syntaxErrorMsg
}

func extractCommandParts(input string) []string {
	regex := regexp.MustCompile(`[^\s]+`)
	return regex.FindAllString(input, -1)
}
