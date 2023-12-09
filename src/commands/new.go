package commands

import (
	"fmt"

	"pagecli.com/main/definition"
	"pagecli.com/main/logging"
)

type New struct {
}

var new CommandInfo = CommandInfo{
	DisplayName:              "new",
	ExecutionOutput:          "",
	ExecutionOk:              true,
	MinimumExpectedArgs:      0,
	MaximumExpectedArguments: 0,
}

func (n New) BindArgs() {

}

func (n New) UsageInfoShort() string {
	return "creates a new page.yml definition file with default values"
}

func (n New) UsageInfoExpanded() string {
	extendedUsage := fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Summary:")
	extendedUsage += fmt.Sprintln(new.DisplayName, "-", n.UsageInfoShort())
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Description:")
	extendedUsage += fmt.Sprintln("creates a new page.yml definition file using the default registrar and host provider that have been configured using the conf command. If a default registrar or host have not been configured, the 'page' registrar/host will be the default. The page.yml template contains the minimum fields required to create a new page. By default, the page.yml template is created in the directory where the command is executed.")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Additonal arguments and options:")
	extendedUsage += fmt.Sprintln("does not require any additional arguments or options")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Example usage:")
	extendedUsage += fmt.Sprintln("page", new.DisplayName)
	extendedUsage += fmt.Sprintln()
	return extendedUsage
}

func (n New) UsageCategory() int {
	return 0
}

func (n New) Execute() {
	err := definition.WriteDefinitionFile()
	if err != nil {
		logging.LogException(err.Error(), true)
	}
}

func (n New) Output() string {
	return ""
}
