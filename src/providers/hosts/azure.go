package providers

import "fmt"

type Azure struct {
}

func (a Azure) Deploy() bool {
	return true
}

func (a Azure) ConfigureHost() bool {
	fmt.Println("configured azure host")
	return true
}
