package commands

import (
	"fmt"

	"builtonpage.com/main/logging"
)

type Help struct {
}

var help CommandInfo = CommandInfo{
	DisplayName:              "help",
	ExecutionOutput:          "",
	ExecutionOk:              true,
	MinimumExpectedArgs:      1,
	MaximumExpectedArguments: 1,
	OrderedArgLabel:          []string{"commandName"},
	ArgValues: map[string]string{
		"commandName": "",
	},
}

func (h Help) BindArgs() {
	logMessage := ""

	_, ok := commandLookup[help.ArgValues["commandName"]]
	if !ok {
		logMessage = fmt.Sprint("unrecognized command '", help.ArgValues["commandName"], "'. Expected a valid command. See 'page' for valid commands.\n")
		help.ExecutionOk = false
		help.ExecutionOutput += logMessage
		logging.LogException(logMessage, false)
		return
	}
}

func (h Help) UsageInfoShort() string {
	return "get expanded description and usage info for any valid command"
}

func (h Help) UsageInfoExpanded() string {
	extendedUsage := fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Summary")
	extendedUsage += fmt.Sprintln(help.DisplayName, "-", h.UsageInfoShort())
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Description")
	extendedUsage += fmt.Sprintln("Prints expanded description for a command, required arguments, and example usage.")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Arguments")
	extendedUsage += fmt.Sprintln("Expects", help.MinimumExpectedArgs, "additional argument (the the name of the command).")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Example usage")
	extendedUsage += fmt.Sprintln("page", help.DisplayName, "up")
	extendedUsage += fmt.Sprintln("page", help.DisplayName, "conf")
	extendedUsage += fmt.Sprintln("page", help.DisplayName, "init")
	extendedUsage += fmt.Sprintln()
	return extendedUsage
}

func (h Help) UsageCategory() int {
	return -1
}

func (h Help) Execute() {
	if !help.ExecutionOk {
		return
	}

	command := commandLookup[help.ArgValues["commandName"]]
	help.ExecutionOutput = command.UsageInfoExpanded()
}

func (h Help) Output() string {
	return help.ExecutionOutput
}
