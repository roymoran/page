package providers

import "fmt"

type PageHost struct {
}

func (p PageHost) Deploy() bool {
	return true
}

func (p PageHost) ConfigureHost() bool {
	fmt.Println("configured page host")
	return true
}

func (p PageHost) HostProviderDefinition() string {
	return ""
}
