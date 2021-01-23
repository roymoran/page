package providers

import "fmt"

type AmazonWebServices struct {
}

func (aws AmazonWebServices) Deploy() bool {
	return true
}

func (aws AmazonWebServices) ConfigureHost() bool {
	fmt.Println("configured aws host")
	return true
}
