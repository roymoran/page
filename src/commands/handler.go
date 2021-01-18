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
	LoadArgs(args []string)
	Execute()
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
// Empty string is special case where no command name
// is provided to the program. The output for this case
//  is usage info on available commands
var commandLookup = map[string]ICommand{
	// "init": Init{},
	// "up":   Up{},
	conf.Name: Conf{},
	// "help": Help{},
	"": None{},
}

type ProgramArgs struct {
	OrderedArgLabel     []string
	ArgValues           map[string]string
	AdditionalArgValues []string
}

var programArgs ProgramArgs = ProgramArgs{
	ArgValues: map[string]string{
		"programName": "",
		"command":     "",
	},
	OrderedArgLabel: []string{"programName", "command"},
}

// Handle is the entry point that begins execution
// of all commands. It parses the command line args
// and calls execute on the appropriate command.
func Handle(args []string) string {
	for i, arg := range args {
		if i > len(programArgs.OrderedArgLabel)-1 {
			programArgs.AdditionalArgValues = args[i:]
			break
		}
		programArgs.ArgValues[programArgs.OrderedArgLabel[i]] = arg
	}

	command, commandValid := commandLookup[programArgs.ArgValues["command"]]

	if !commandValid {
		return fmt.Sprint("unrecognized command ", programArgs.ArgValues["command"], ". See 'page' for list of valid commands.\n")
	}

	command.LoadArgs(programArgs.AdditionalArgValues)
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

func LoadArgsIfAny(args []string, command ICommand) bool {

	return true
}
