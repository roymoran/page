package command

type Up struct {
}

func (u Up) UsageInfoShort() string {
	return "publishes the page using the page definition file provided"
}

func (u Up) UsageInfoExpanded() string {
	return ""
}

func (u Up) UsageCategory() int {
	return 1
}

func (u Up) Execute() bool {
	return true
}

func (u Up) Output() string {
	return ""
}

func (u Up) ValidArgs() bool {
	return true
}
