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

func NewExecutor(dataStore *datastore.DataStore) *Executor {
	handler := &Executor{
		dataSource: dataStore,
		commands:   make(map[string]actions.Command),
	}
	handler.RegisterCommand(setCommand, actions.NewSetCommand(dataStore))
	handler.RegisterCommand(getCommand, actions.NewGetCommand(dataStore))
	handler.RegisterCommand(delCommand, actions.NewDelCommand(dataStore))
	handler.RegisterCommand(keysCommand, actions.NewKeysCommand(dataStore))
	handler.RegisterCommand(zaddCommand, actions.NewZAddCommand(dataStore))
	handler.RegisterCommand(zremCommand, actions.NewZRemCommand(dataStore))
	handler.RegisterCommand(zscoreCommand, actions.NewZScoreCommand(dataStore))
	handler.RegisterCommand(zqueryCommand, actions.NewZQueryCommand(dataStore))
	handler.RegisterCommand(zshowCommand, actions.NewZShowCommand(dataStore))
	handler.RegisterCommand(expireCommand, actions.NewExpireCommand(dataStore))
	handler.RegisterCommand(ttlCommand, actions.NewTTLCommand(dataStore))
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
