package command

import (
	"log"
	"regexp"

	"github.com/miladbarzideh/goldis/internal/command/actions"
	"github.com/miladbarzideh/goldis/internal/datastore"
)

const (
	getCommand    = "get"
	setCommand    = "set"
	delCommand    = "del"
	keysCommand   = "keys"
	zaddCommand   = "zadd"
	zremCommand   = "zrem"
	zscoreCommand = "zscore"
	zqueryCommand = "zquery"
	zshowCommand  = "zshow"
	expireCommand = "pexpire"
	ttlCommand    = "pttl"
)

type Executor struct {
	dataSource *datastore.DataStore
	commands   map[string]actions.Command
}

func NewExecutor() *Executor {
	ds := datastore.NewDataStore()
	handler := &Executor{
		dataSource: ds,
		commands:   make(map[string]actions.Command),
	}
	handler.RegisterCommand(setCommand, actions.NewSetCommand(*ds))
	handler.RegisterCommand(getCommand, actions.NewGetCommand(*ds))
	handler.RegisterCommand(delCommand, actions.NewDelCommand(*ds))
	handler.RegisterCommand(keysCommand, actions.NewKeysCommand(*ds))
	handler.RegisterCommand(zaddCommand, actions.NewZAddCommand(*ds))
	handler.RegisterCommand(zremCommand, actions.NewZRemCommand(*ds))
	handler.RegisterCommand(zscoreCommand, actions.NewZScoreCommand(*ds))
	handler.RegisterCommand(zqueryCommand, actions.NewZQueryCommand(*ds))
	handler.RegisterCommand(zshowCommand, actions.NewZShowCommand(*ds))
	handler.RegisterCommand(expireCommand, actions.NewExpireCommand(*ds))
	handler.RegisterCommand(ttlCommand, actions.NewTTLCommand(*ds))
	return handler
}

func (h *Executor) RegisterCommand(key string, command actions.Command) {
	h.commands[key] = command
}

func (h *Executor) Execute(input string) string {
	commandParts := extractCommandParts(input)
	if commandParts == nil || len(commandParts) < 1 {
		return actions.SyntaxErrorMsg
	}
	commandKey, args := commandParts[0], commandParts[1:]
	log.Printf("Command %s will be executed", commandKey)
	if command, ok := h.commands[commandKey]; ok {
		return command.Execute(args)
	}
	return actions.SyntaxErrorMsg
}

func extractCommandParts(input string) []string {
	regex := regexp.MustCompile(`[^\s]+`)
	return regex.FindAllString(input, -1)
}
