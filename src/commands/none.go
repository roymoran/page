package commands

import (
	"fmt"

	"builtonpage.com/main/constants"
)

type None struct {
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

func (n None) Execute() bool {
	return true
}

func (n None) Output() string {
	return fmt.Sprint("page version ", constants.AppVersion(), "\n\n", BuildUsageInfo(), "\n")
}

func (n None) ValidArgs() bool {
	return true
}
