/*
Package command provides an entry point for
the execution of all cli commands e.g. init,
up, etc.
*/
package command

type ICommand interface {
	ValidArgs() bool
	Execute()
	Output()
}

type Command struct {
	Name string
}

var commandLookup = map[string]ICommand{
	"init": Init{},
}

func Handle(args []string) {
	// TODO: Implement
	// Arg sanitation/preprocessing
	// Command execution
}
