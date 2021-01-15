/*
Package command provides an entry point for
the execution of all cli commands e.g. init,
up, etc.
*/
package command

type ICommand interface {
	ValidArgs() bool
	Execute() bool
	Output() string
}

type Command struct {
	Name string
}

// commandLookup creates a mapping of each command
// name with the struct that implements ICommand.
// "none" is special case where no command name
// is passed in which case the return is the
// version info of the cli tool
var commandLookup = map[string]ICommand{
	"none": None{},
	"init": Init{},
}

// Handle is the entry point that begins execution
// of all commands. It parses the command line args
// and calls execute on the appropriate command.
func Handle(args []string) string {
	var command ICommand

	if len(args) == 0 {
		// special case, cli tool is executed
		// with no arguments
		command = commandLookup["none"]
	} else {
		command = commandLookup[args[0]]
	}

	command.Execute()
	return command.Output()
}