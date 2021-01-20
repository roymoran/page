package commands

import "fmt"

type Help struct {
	DisplayName              string
	ExecutionOutput          string
	ExecutionOk              bool
	MinimumExpectedArgs      int
	MaximumExpectedArguments int
	OrderedArgLabel          []string
	ArgValues                map[string]string
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

func (h Help) LoadArgs() {

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

	command := commandLookup[help.ArgValues["commandName"]]
	help.ExecutionOutput = command.UsageInfoExpanded()
}

func (h Help) Output() string {
	return help.ExecutionOutput
}
