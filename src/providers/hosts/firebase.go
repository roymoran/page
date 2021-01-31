package providers

import "fmt"

type Firebase struct {
}

func (f Firebase) Deploy() bool {
	return true
}

func (f Firebase) ConfigureHost(alias string, definitionFilePath string, stateFilePath string) (bool, error) {
	fmt.Println("configured Firebase host")
	return true, nil
}

func (f Firebase) HostProviderDefinition() []byte {
	return []byte{}
}
