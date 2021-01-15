/*
Package command provides an entry point for
the execution of all cli commands e.g. init,
up, etc.
*/
package command

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type Command struct {
	Name string
}
type ICommand interface {
	ValidArgs() bool
	Execute() bool
	Output() string
	UsageInfoShort() string
	UsageInfoExpanded() string
}

// commandLookup creates a mapping of each command
// name with the struct that implements ICommand.
// "none" is special case where no command name
// is passed in which case the return is the
// version info of the cli tool
var commandLookup = map[string]ICommand{
	"init":      Init{},
	"up":        Up{},
	"configure": Configure{},
	"none":      None{},
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

// make your page live\n  up	publishes your page using the page defintion provided
func BuildUsageInfo() string {
	usageInfo := "Common Page commands. For specific command usage use 'page [command_name] --help':\n\nstart, publish, and configure a new page project\n"
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	for commandName, command := range commandLookup {
		if commandName == "none" {
			continue
		}
		usageInfo += fmt.Sprint("\n")
		usageInfo += fmt.Sprint(commandName, "\t", command.UsageInfoShort())
	}
	writer.Flush()
	return fmt.Sprint(usageInfo)
}
