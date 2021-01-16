/*
Package command provides an entry point for
the execution of all cli commands e.g. init,
up, etc.
*/
package command

import (
	"fmt"
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
	UsageCategory() int
}

var usageCategories = []string{
	"start a new page project",
	"publish page project",
	"configure domain registrar and host provider for your projects",
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

// TODO: Use tabwriter for fixed widths
func BuildUsageInfo() string {
	usageInfo := "Common Page commands:\n"

	for catergoryId, category := range usageCategories {
		usageInfo += fmt.Sprint("\n", category, "\n")

		for commandName, command := range commandLookup {
			if commandName == "none" {
				continue
			}

			if command.UsageCategory() == catergoryId {
				usageInfo += fmt.Sprint("  ", commandName, "\t\t", command.UsageInfoShort(), "\n")
			}
		}
	}

	return fmt.Sprint(usageInfo, "\n\n", "For specific command usage use 'page [command_name] --help'")
}
