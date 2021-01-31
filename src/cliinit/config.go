package cliinit

type PageConfig struct {
	TfPath       string           `json:"tfPath"`
	TfExecPath   string           `json:"tfExecPath"`
	TFVersion    string           `json:"tfVersion"`
	Providers    []ProviderConfig `json:"providers"`
	ConfigStatus bool             `json:"configStatus"`
}

type ProviderConfig struct {
	Id               string `json:"id"`
	Alias            string `json:"alias"`
	Type             string `json:"type"` // 'registrar' or 'host'
	HostName         string `json:"hostName"`
	Auth             string `json:"auth"`
	Default          bool   `json:"default"`
	TfDefinitionPath string `json:"tfDefinitionPath"`
	TfStatePath      string `json:"tfStatePath"`
}
