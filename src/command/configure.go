package command

type Configure struct {
}

func (c Configure) UsageInfoShort() string {
	return "configures a new page.yml definition file"
}

func (c Configure) UsageInfoExpanded() string {
	return ""
}

func (c Configure) UsageCategory() int {
	return 2
}

func (c Configure) Execute() bool {
	return true
}

func (c Configure) Output() string {
	return ""
}

func (c Configure) ValidArgs() bool {
	return true
}
