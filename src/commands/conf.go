package command

type Conf struct {
}

func (c Conf) UsageInfoShort() string {
	return "configures defaults for domain registrar and host provider"
}

func (c Conf) UsageInfoExpanded() string {
	return ""
}

func (c Conf) UsageCategory() int {
	return 2
}

func (c Conf) Execute() bool {
	return true
}

func (c Conf) Output() string {
	return ""
}

func (c Conf) ValidArgs() bool {
	return true
}
