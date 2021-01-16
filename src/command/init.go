package command

import "builtonpage.com/main/definition"

type Init struct {
}

func (i Init) UsageInfoShort() string {
	return "creates a new page.yml definition file"
}

func (i Init) UsageInfoExpanded() string {
	return ""
}

func (i Init) UsageCategory() int {
	return 0
}

func (i Init) Execute() bool {
	ok := definition.WriteDefinitionFile()
	return ok
}

func (i Init) Output() string {
	return ""
}

func (i Init) ValidArgs() bool {
	return true
}
