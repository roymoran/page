package command

import (
	"fmt"

	"builtonpage.com/main/constants"
)

type None struct {
}

func (i None) Execute() bool {
	return true
}

func (i None) Output() string {
	return fmt.Sprint("page version ", constants.AppVersion(), "\n")
}

func (i None) ValidArgs() bool {
	return true
}
