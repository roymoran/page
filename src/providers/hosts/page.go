package providers

import "builtonpage.com/main/cliinit"

type PageHost struct {
}

func (p PageHost) Deploy() bool {
	return true
}

func (p PageHost) AddHost(alias string, definitionFilePath string, stateFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             "page",
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      stateFilePath,
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

func (p PageHost) HostProviderDefinition() []byte {
	return []byte{}
}
