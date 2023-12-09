package commands

import (
	"fmt"
	"strings"

	"pagecli.com/main/constants"
)

type None struct {
}

var none CommandInfo = CommandInfo{
	DisplayName:              "",
	ExecutionOutput:          "",
	ExecutionOk:              true,
	MinimumExpectedArgs:      0,
	MaximumExpectedArguments: 0,
}

func (n None) UsageInfoShort() string {
	return ""
}

func (n None) UsageInfoExpanded() string {
	return ""
}

func (n None) UsageCategory() int {
	return -1
}

func (n None) Execute() {
}

func (n None) Output() string {
	return fmt.Sprint(constants.AppName(), " ", constants.AppVersion(), "\n", "author(s); "+strings.Join(constants.AppAuthors(), ", "), "\n\n", BuildUsageInfo(), "\n")
}

func (n None) BindArgs() {
}
