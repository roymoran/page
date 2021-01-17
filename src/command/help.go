package command

type Help struct {
}

var executionOutput string

func (h Help) UsageInfoShort() string {
	return "creates a new page.yml definition file with default values"
}

func (h Help) UsageInfoExpanded() string {
	return ""
}

func (h Help) UsageCategory() int {
	return -1
}

func (h Help) Execute() bool {
	command := commandLookup["init"]
	executionOutput = command.UsageInfoExpanded()
	return true
}

func (h Help) Output() string {
	return executionOutput
}

func (h Help) ValidArgs() bool {
	return true
}
