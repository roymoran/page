/*
Package commands provides an entry point for
the execution of all cli commands e.g. init,
up, etc.
*/
package commands

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"pagecli.com/main/logging"
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
	"create",
	"configure",
}

// commandLookup creates a mapping of each command
// name with the struct that implements ICommand.
// Empty string is special case where no command name
// is provided to the program. The output for this case
//
//	is usage info on available commands
var commandLookup = map[string]ICommand{
	new.DisplayName:     New{},
	up.DisplayName:      Up{},
	conf.DisplayName:    Conf{},
	help.DisplayName:    Help{},
	builder.DisplayName: Builder{},
	"":                  None{},
	"infra":             Infra{}, // Hidden command
}

var commandInfoMap = map[string]*CommandInfo{
	new.DisplayName:     &new,
	conf.DisplayName:    &conf,
	help.DisplayName:    &help,
	none.DisplayName:    &none,
	up.DisplayName:      &up,
	infra.DisplayName:   &infra,
	builder.DisplayName: &builder,
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

// OutputChannel used to send intermediate output
// for certain commands that require several seconds
// to run like 'page up'
var OutputChannel chan string = nil

// Handle is the entry point that begins execution
// of all commands. It parses the command line args
// and calls execute on the appropriate command.
func Handle(args []string, channel chan string) {
	OutputChannel = channel
	for i, arg := range args {
		if i > len(programArgs.OrderedArgLabel)-1 {
			programArgs.AdditionalArgValues = args[i:]
			break
		}
		programArgs.ArgValues[programArgs.OrderedArgLabel[i]] = arg
	}

	command, commandValid := commandLookup[programArgs.ArgValues["command"]]

	if programArgs.ArgValues["command"] == "" {
		logging.NewAnalytics().FireEvent("command", logging.EventParams{
			"name": "(none)",
			"args": "(none)",
		})
	} else {
		logging.NewAnalytics().FireEvent("command", logging.EventParams{
			"name": programArgs.ArgValues["command"],
			"args": strings.Join(programArgs.AdditionalArgValues, " "),
		})
	}

	if !commandValid {
		OutputChannel <- fmt.Sprint("unrecognized command ", programArgs.ArgValues["command"], ". See 'page' for list of valid commands.\n")
		logging.SendLog(logging.LogRecord{
			Level:   "error",
			Message: "Unrecognized command " + programArgs.ArgValues["command"],
		})
		close(OutputChannel)
		return
	}

	ValidateArgs(commandInfoMap[programArgs.ArgValues["command"]], programArgs.AdditionalArgValues)
	command.BindArgs()
	command.Execute()
	OutputChannel <- command.Output()
	close(OutputChannel)
}

func BuildUsageInfo() string {
	usageInfo := "common commands:\n"
	var b bytes.Buffer
	tabwriter.NewWriter(&b, 0, 8, 1, '\t', tabwriter.AlignRight)
	for catergoryID, category := range usageCategories {
		usageInfo += fmt.Sprint("\n", category, "\n")

		sortedCommandNames := make([]string, 0)
		for k := range commandLookup {
			sortedCommandNames = append(sortedCommandNames, k)
		}

		sort.Strings(sortedCommandNames)

		for _, commandName := range sortedCommandNames {
			if commandName == "none" || commandName == "infra" {
				continue
			}

			if commandLookup[commandName].UsageCategory() == catergoryID {
				fmt.Fprint(&b, "   ", commandName, "\t\t", commandLookup[commandName].UsageInfoShort(), "\n")
				usageInfo += fmt.Sprint(b.String())
				b.Reset()
			}
		}
	}

	return fmt.Sprint(usageInfo, "\n\n", "For specific command usage use 'page help <command>'\n")
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
		commandInfo.ExecutionOutput += fmt.Sprintln()
		commandInfo.ExecutionOutput += fmt.Sprint("See 'page help ", commandInfo.DisplayName, "' for usage info.")
		commandInfo.ExecutionOutput += fmt.Sprintln()
		logging.SendLog(logging.LogRecord{
			Level:   "error",
			Message: "invalid command arguements",
		})
		return
	}

	for i, arg := range args {
		commandInfo.ArgValues[commandInfo.OrderedArgLabel[i]] = arg
	}
}
