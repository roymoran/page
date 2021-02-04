package cliinit

type PageConfig struct {
	// TODO: Is it worth storing these paths?
	// Could we not access these programatically
	// with os.UserHomeDir() each time we need
	// them? Values stored in config may be brittle
	// to filesystem changes
	TfPath       string           `json:"tfPath"`
	TfExecPath   string           `json:"tfExecPath"`
	TFVersion    string           `json:"tfVersion"`
	Providers    []ProviderConfig `json:"providers"`
	ConfigStatus bool             `json:"configStatus"`
}

type ProviderConfig struct {
	Alias            string `json:"alias"`
	Type             string `json:"type"` // 'registrar' or 'host'
	Name             string `json:"name"`
	Auth             string `json:"auth"`
	Default          bool   `json:"default"`
	TfDefinitionPath string `json:"tfDefinitionPath"`
	TfStatePath      string `json:"tfStatePath"`
}
