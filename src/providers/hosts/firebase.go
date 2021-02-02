package providers

import "builtonpage.com/main/cliinit"

type Firebase struct {
}

func (f Firebase) Deploy() bool {
	return true
}

func (f Firebase) CoAddHostnfigureHost(alias string, definitionFilePath string, stateFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             "firebase",
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      stateFilePath,
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

func (f Firebase) HostProviderDefinition() []byte {
	return []byte{}
}
