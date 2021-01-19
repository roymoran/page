package commands

import "fmt"

type Help struct {
	DisplayName              string
	ExecutionOutput          string
	ExecutionOk              bool
	MinimumExpectedArgs      int
	MaximumExpectedArguments int
}

type HelpArgs struct {
	OrderedArgLabel []string
	ArgValues       map[string]string
}

var help Help = Help{
	DisplayName:              "help",
	ExecutionOutput:          "",
	ExecutionOk:              true,
	MinimumExpectedArgs:      1,
	MaximumExpectedArguments: 1,
}

var helpArgs HelpArgs = HelpArgs{
	OrderedArgLabel: []string{"commandName"},
	ArgValues: map[string]string{
		"commandName": "",
	},
}

func (h Help) LoadArgs(args []string) {
	if len(args) < help.MinimumExpectedArgs {
		help.ExecutionOk = false
		help.ExecutionOutput += fmt.Sprintln(help.DisplayName, "expects at least", help.MinimumExpectedArgs, "arguments, received", len(args))
		return
	}

	if len(args) > help.MaximumExpectedArguments {
		help.ExecutionOk = false
		help.ExecutionOutput += fmt.Sprintln(help.DisplayName, "expects at most", help.MaximumExpectedArguments, "arguments, received", len(args))
		return
	}

	for i, arg := range args {
		helpArgs.ArgValues[helpArgs.OrderedArgLabel[i]] = arg
	}
}

func (h Help) UsageInfoShort() string {
	return "creates a new page.yml definition file with default values"
}

func (h Help) UsageInfoExpanded() string {
	return ""
}

func (h Help) UsageCategory() int {
	return -1
}

func (h Help) Execute() {
	if !help.ExecutionOk {
		help.ExecutionOutput += fmt.Sprintln("")
		help.ExecutionOutput += fmt.Sprint("See 'page help ", help.DisplayName, "' for usage info.\n")
		return
	}

	command := commandLookup[helpArgs.ArgValues["commandName"]]
	help.ExecutionOutput = command.UsageInfoExpanded()
}

func (h Help) Output() string {
	return help.ExecutionOutput
}

func (h Help) ValidArgs() bool {
	return true
}
