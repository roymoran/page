package providers

// TODO: Add logic so that there are character restrictions
// to resource names for terraform resources. Ensure these are
// enforced to avoid errors when tf files are programatiically
// created

// ProviderTemplate defines minimum fields required to
// create a new terraform host directory. This data is
// written to disk as a json file and 'terraform init'
// is used to initialize the directory and dowload the
// provider plugin
// More info on terraform providers can found at link
// https://registry.terraform.io/browse/providers
type ProviderTemplate struct {
	Terraform RequiredProviders `json:"terraform,omitempty"`
}

// ProviderConfigTemplate defines configuration for
// a provider such as the region resources will be
// deployed to.
// https://registry.terraform.io/browse/providers
type ProviderConfigTemplate struct {
	Provider map[string]interface{} `json:"provider,omitempty"`
}

// BaseInfraTemplate defines the resources required to
// create the infrastructure on which all sites
// will be hosted.
type BaseInfraTemplate struct {
	Resource map[string]interface{} `json:"resource,omitempty"`
	Output   map[string]interface{} `json:"output,omitempty"`
}

// SiteTemplate defines the resources required to
// create a site on existing infrastructure
type SiteTemplate struct {
	Site map[string]interface{} `json:"resource,omitempty"`
}

type RequiredProviders struct {
	RequiredProvider map[string]Provider `json:"required_providers"`
}

type Provider struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}

type ProviderConfig struct {
	Profile string `json:"profile"`
	Region  string `json:"region"`
}

type ModuleTemplate struct {
	Module map[string]interface{} `json:"module,omitempty"`
	Output map[string]interface{} `json:"output,omitempty"`
	Input  map[string]interface{} `json:"input,omitempty"`
}

type HostModuleTemplate struct {
	Module map[string]HostModuleProperties `json:"module,omitempty"`
}

type HostModuleProperties struct {
	Certificates map[string]Certificate `json:"certificates,omitempty"`
	Source       string                 `json:"source,omitempty"`
}
type Certificate struct {
	CertificateChain string `json:"certificate_chain,omitempty"`
	CertificatePem   string `json:"certificate_pem,omitempty"`
	PrivateKeyPem    string `json:"private_key_pem,omitempty"`
}
