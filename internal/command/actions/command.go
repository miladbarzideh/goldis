package actions

const SyntaxErrorMsg = "(error) ERR syntax error"

type Command interface {
	Execute(args []string) string
}
