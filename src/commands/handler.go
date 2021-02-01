/*
Package commands provides an entry point for
the execution of all cli commands e.g. init,
up, etc.
*/
package commands

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

type ICommand interface {
	BindArgs()
	Execute()
	Output() string
	UsageInfoShort() string
	UsageInfoExpanded() string
	UsageCategory() int
}

type CommandInfo struct {
	DisplayName              string
	ExecutionOutput          string
	ExecutionOk              bool
	MinimumExpectedArgs      int
	MaximumExpectedArguments int
	OrderedArgLabel          []string
	ArgValues                map[string]string
}

var usageCategories = []string{
	"start new page project",
	"publish page project",
	"configure default registrar/host",
}

// commandLookup creates a mapping of each command
// name with the struct that implements ICommand.
// Empty string is special case where no command name
// is provided to the program. The output for this case
//  is usage info on available commands
var commandLookup = map[string]ICommand{
	initCommand.DisplayName: Init{},
	up.DisplayName:          Up{},
	conf.DisplayName:        Conf{},
	help.DisplayName:        Help{},
	"":                      None{},
}

var commandInfoMap = map[string]*CommandInfo{
	initCommand.DisplayName: &initCommand,
	conf.DisplayName:        &conf,
	help.DisplayName:        &help,
	none.DisplayName:        &none,
	up.DisplayName:          &up,
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
func Handle(args []string, channel chan string) {
	for i, arg := range args {
		if i > len(programArgs.OrderedArgLabel)-1 {
			programArgs.AdditionalArgValues = args[i:]
			break
		}
		programArgs.ArgValues[programArgs.OrderedArgLabel[i]] = arg
	}

	command, commandValid := commandLookup[programArgs.ArgValues["command"]]

	if !commandValid {
		channel <- fmt.Sprint("unrecognized command ", programArgs.ArgValues["command"], ". See 'page' for list of valid commands.\n")
		close(channel)
		return
	}

	ValidateArgs(commandInfoMap[programArgs.ArgValues["command"]], programArgs.AdditionalArgValues)
	command.BindArgs()
	command.Execute()
	channel <- command.Output()
	close(channel)
	return
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

func ValidateArgs(commandInfo *CommandInfo, args []string) {
	if len(args) < commandInfo.MinimumExpectedArgs {
		commandInfo.ExecutionOk = false
		commandInfo.ExecutionOutput += fmt.Sprintln(commandInfo.DisplayName, "expects at least", commandInfo.MinimumExpectedArgs, "arguments, received", len(args))
	}

	if len(args) > commandInfo.MaximumExpectedArguments {
		commandInfo.ExecutionOk = false
		commandInfo.ExecutionOutput += fmt.Sprintln(commandInfo.DisplayName, "expects at most", commandInfo.MaximumExpectedArguments, "arguments, received", len(args))
	}

	// TODO: Generalize so that on 'page help...' message is appended to all failing commands
	if !commandInfo.ExecutionOk {
		conf.ExecutionOutput += fmt.Sprintln()
		conf.ExecutionOutput += fmt.Sprint("See 'page help ", commandInfo.DisplayName, "' for usage info.")
		conf.ExecutionOutput += fmt.Sprintln()
		return
	}

	for i, arg := range args {
		commandInfo.ArgValues[commandInfo.OrderedArgLabel[i]] = arg
	}
}
