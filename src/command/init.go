package command

import "builtonpage.com/main/definition"

type Init struct {
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
