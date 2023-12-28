package commands

import (
	"fmt"

	"gitlab.nasapps.net/page/builder/buildercore"
)

type Builder struct {
}

var builder CommandInfo = CommandInfo{
	DisplayName:              "build",
	ExecutionOutput:          "",
	ExecutionOk:              true,
	MinimumExpectedArgs:      0,
	MaximumExpectedArguments: 0,
}

func (b Builder) BindArgs() {

}

func (b Builder) Execute() {
	buildercore.Start()
}

func (b Builder) Output() string {
	return builder.ExecutionOutput
}

func (b Builder) UsageCategory() int {
	return 0
}

func (b Builder) UsageInfoExpanded() string {
	extendedUsage := fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Summary")
	extendedUsage += fmt.Sprintln(builder.DisplayName, "-", b.UsageInfoShort())
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Description")
	extendedUsage += fmt.Sprintln("Automates the creation of webpages using textual prompts as input. For example, you could input, 'Create a webpage for my project, starting with a headline 'Hello, World!' centered both vertically and horizontally.' Based on your input, it generates HTML, CSS, and JavaScript files in your current directory, which are then opened in your default browser for review. You have the flexibility to manually modify the generated code, or continue evolving the webpage by providing additional input. This tool is designed to kickstart your webpage projects without in depth knowledge of HTML, CSS, or JavaScript.")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Arguments")
	extendedUsage += fmt.Sprintln("no arguments for", builder.DisplayName)
	extendedUsage += fmt.Sprintln()
	return extendedUsage
}

func (b Builder) UsageInfoShort() string {
	return "create a webpage with generative AI via text input"
}
