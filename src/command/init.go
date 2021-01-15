package command

type Init struct {
}

func (i Init) Execute() bool {
	return true
}

func (i Init) Output() string {
	return "TODO: implement"
}

func (i Init) ValidArgs() bool {
	return true
}
