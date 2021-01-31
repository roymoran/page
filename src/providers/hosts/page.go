package providers

import "fmt"

type PageHost struct {
}

func (p PageHost) Deploy() bool {
	return true
}

func (p PageHost) ConfigureHost(alias string, definitionFilePath string, stateFilePath string) (bool, error) {
	fmt.Println("configured page host")
	return true, nil
}

func (p PageHost) HostProviderDefinition() []byte {
	return []byte{}
}
