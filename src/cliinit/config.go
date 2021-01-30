package cliinit

type PageConfigJson struct {
	TfPath       string           `json:"tfpath"`
	TFVersion    string           `json:"tfversion"`
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
