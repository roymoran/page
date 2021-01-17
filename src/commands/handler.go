/*
Package command provides an entry point for
the execution of all cli commands e.g. init,
up, etc.
*/
package commands

import (
	"bytes"
	"fmt"
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
	UsageCategory() int
}

var usageCategories = []string{
	"start a new page project",
	"publish a page project",
	"configure defaults for registrar/host provider",
}

// commandLookup creates a mapping of each command
// name with the struct that implements ICommand.
// "none" is special case where no command name
// is passed in which case the return is the
// version info of the cli tool
var commandLookup = map[string]ICommand{
	"init": Init{},
	"up":   Up{},
	"conf": Conf{},
	"help": Help{},
	"none": None{},
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
		// check if valid command otherwise
		// return with invalid command message
		if val, ok := commandLookup[args[0]]; ok {
			command = val
		} else {
			return fmt.Sprint("unrecognized command ", args[0], ". See 'page' for list of valid commands.\n")
		}
	}

	command.Execute()
	return command.Output()
}

func BuildUsageInfo() string {
	usageInfo := "Common Page commands:\n"
	var b bytes.Buffer
	tabwriter.NewWriter(&b, 0, 8, 1, '\t', tabwriter.AlignRight)
	for catergoryId, category := range usageCategories {
		usageInfo += fmt.Sprint("\n", category, "\n")

		for commandName, command := range commandLookup {
			if commandName == "none" {
				continue
			}

			if command.UsageCategory() == catergoryId {
				fmt.Fprint(&b, "   ", commandName, "\t\t", command.UsageInfoShort(), "\n")
				usageInfo += fmt.Sprint(b.String())
				b.Reset()
			}
		}
	}

	return fmt.Sprint(usageInfo, "\n\n", "For specific command usage use 'page help <command>'")
}
