package cliinit

type PageConfig struct {
	TfPath       string           `json:"tfPath"`
	TfExecPath   string           `json:"tfExecPath"`
	TFVersion    string           `json:"tfVersion"`
	Providers    []ProviderConfig `json:"providers"`
	ConfigStatus bool             `json:"configStatus"`
}

type ProviderConfig struct {
	Id           string `json:"id"`
	Type         string `json:"type"` // 'registrar' or 'host'
	Name         string `json:"name"`
	Auth         string `json:"auth"`
	Default      bool   `json:"default"`
	TfConfigPath string `json:"tfConfigPath"`
	TfStatePath  string `json:"tfStatePath"`
}
