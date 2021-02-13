package providers

import (
	"fmt"

	"builtonpage.com/main/definition"
	providers "builtonpage.com/main/providers/hosts"
)

var providerTemplate providers.ProviderTemplate = providers.ProviderTemplate{
	Terraform: providers.RequiredProviders{
		RequiredProvider: map[string]providers.Provider{
			"acme": {
				Source:  "vancluever/acme",
				Version: "2.0.0",
			},
		},
	},
}

var acmeResourceTemplate map[string]interface{} = map[string]interface{}{
	"resource": map[string]interface{}{
		"tls_private_key": map[string]interface{}{
			"private_key": map[string]interface{}{
				"algorithm": "RSA",
			},
		},

		"acme_registration": map[string]interface{}{
			"reg": map[string]interface{}{
				"account_key_pem": "${tls_private_key.private_key.private_key_pem}",
				// TODO: Read from user
				"email_address": "romoran1@outlook.com",
			},
		},
	},
}

type IRegistrar interface {
	ConfigureAuth() error
	ConfigureRegistrar(definition.PageDefinition) bool
	ConfigureDns() bool
	AddRegistrar(string) error
}

func (rp RegistrarProvider) Add(name string, channel chan string) error {
	var alias string = AssignAliasName("registrar")

	registrarProvider := SupportedProviders.Providers["registrar"].(RegistrarProvider)
	registrar := registrarProvider.Supported[name]
	registrar.ConfigureAuth()
	// acmeResourceTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "acme.tf.json")
	registrar.AddRegistrar(alias)

	return nil
}

func (rp RegistrarProvider) List(name string, channel chan string) error {
	for _, registrarName := range SupportedRegistrars {
		channel <- fmt.Sprint(registrarName)
	}
	return nil
}
