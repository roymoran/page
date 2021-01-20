package commands

import (
	"fmt"

	"builtonpage.com/main/definition"
)

type Init struct {
	DisplayName              string
	ExecutionOutput          string
	ExecutionOk              bool
	MinimumExpectedArgs      int
	MaximumExpectedArguments int
}

var initCommand CommandInfo = CommandInfo{
	DisplayName:              "init",
	ExecutionOutput:          "",
	ExecutionOk:              true,
	MinimumExpectedArgs:      0,
	MaximumExpectedArguments: 0,
}

func (i Init) LoadArgs() {

}

func (i Init) UsageInfoShort() string {
	return "creates a new page.yml definition file with default values"
}

func (i Init) UsageInfoExpanded() string {
	extendedUsage := fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Summary:")
	extendedUsage += fmt.Sprintln(initCommand.DisplayName, "-", i.UsageInfoShort())
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Description:")
	extendedUsage += fmt.Sprintln("creates a new page.yml definition file using the default registrar and host provider that have been configured using the conf command. If a default registrar or host have not been configured, the 'page' registrar/host will be the default. The page.yml template contains the minimum fields required to create a new page. By default, the page.yml template is created in the directory where the command is executed.")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Additonal arguments and options:")
	extendedUsage += fmt.Sprintln("does not require any additional arguments or options")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Example usage:")
	extendedUsage += fmt.Sprintln("page", initCommand.DisplayName)
	extendedUsage += fmt.Sprintln()
	return extendedUsage
}

func (i Init) UsageCategory() int {
	return 0
}

func (i Init) Execute() {
	definition.WriteDefinitionFile()
}

func (i Init) Output() string {
	return ""
}
