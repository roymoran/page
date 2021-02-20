package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
	"builtonpage.com/main/providers/hosts"
)

type IRegistrar interface {
	ConfigureAuth() (cliinit.Credentials, error)
	ConfigureRegistrar(string, string, definition.PageDefinition) error
	ConfigureDns() bool
	AddRegistrar(string, cliinit.Credentials) error
}

func (rp RegistrarProvider) Add(name string, channel chan string) error {
	var alias string = AssignAliasName("registrar")

	registrarProvider := SupportedProviders.Providers["registrar"].(RegistrarProvider)
	registrar := registrarProvider.Supported[name]
	credentials, authErr := registrar.ConfigureAuth()

	if authErr != nil {
		return authErr
	}

	providerTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "provider.tf.json")
	// This doesn't work with multiple aliases since
	// provider config file is created only once on host dir configuration
	providerConfigTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "providerconfig.tf.json")
	acmeRegistrationTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "acmeregistration.tf.json")

	moduleTemplatePath := cliinit.ModuleTemplatePath("registrar", alias)

	if !AliasDirectoryConfigured(cliinit.ProviderAliasPath(name, alias)) {
		channel <- fmt.Sprint("Configuring ", name, " registrar...")
		err := InstallAcmeTerraformProvider(name, alias, cliinit.ProviderAliasPath(name, alias), providerTemplatePath, providerConfigTemplatePath, moduleTemplatePath, acmeRegistrationTemplatePath)
		if err != nil {
			return err
		}
	}

	addRegistrarErr := registrar.AddRegistrar(alias, credentials)
	if addRegistrarErr != nil {
		return addRegistrarErr
	}

	return nil
}

func (rp RegistrarProvider) List(name string, channel chan string) error {
	for _, registrarName := range SupportedRegistrars {
		channel <- fmt.Sprint(registrarName)
	}
	return nil
}

func InstallAcmeTerraformProvider(name string, alias string, providerAliasPath string, providerTemplatePath string, providerConfigTemplatePath string, moduleTemplatePath string, acmeRegistrationTemplatePath string) error {
	hostDirErr := os.MkdirAll(providerAliasPath, os.ModePerm)
	if hostDirErr != nil {
		os.Remove(providerAliasPath)
		log.Fatalln("error creating host config directory for", providerAliasPath, hostDirErr)
		return hostDirErr
	}

	moduleTemplatePathErr := ioutil.WriteFile(moduleTemplatePath, registrarModuleTemplate(name, alias), 0644)
	providerTemplatePathErr := ioutil.WriteFile(providerTemplatePath, acmeProviderTemplate(), 0644)
	providerConfigTemplatePathErr := ioutil.WriteFile(providerConfigTemplatePath, acmeProviderConfigTemplate(), 0644)
	// TODO: Read registrartion email from user
	acmeRegistrationTemplatePathErr := ioutil.WriteFile(acmeRegistrationTemplatePath, acmeRegistrationTemplate("roymoran20@gmail.com"), 0644)

	if moduleTemplatePathErr != nil || providerTemplatePathErr != nil || providerConfigTemplatePathErr != nil || acmeRegistrationTemplatePathErr != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(providerAliasPath)
		fmt.Println("failed ioutil.WriteFile for provider template")
		return fmt.Errorf("failed ioutil.WriteFile for provider template")
	}

	err := hosts.TfInit(cliinit.ProvidersPath)
	if err != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(providerAliasPath)
		fmt.Println("failed init on new terraform directory", cliinit.ProvidersPath)
		return err
	}

	return nil
}

func registrarModuleTemplate(providerName string, alias string) []byte {

	var awsProviderTemplate hosts.ModuleTemplate = hosts.ModuleTemplate{
		Module: map[string]interface{}{
			"registrar_" + alias: map[string]interface{}{
				"source": "./" + providerName + "/" + alias,
			},
		},
	}

	file, _ := json.MarshalIndent(awsProviderTemplate, "", " ")
	return file
}

func acmeProviderTemplate() []byte {
	var providerTemplate hosts.ProviderTemplate = hosts.ProviderTemplate{
		Terraform: hosts.RequiredProviders{
			RequiredProvider: map[string]hosts.Provider{
				"acme": {
					Source:  "vancluever/acme",
					Version: "2.0.0",
				},
				"tls": {
					Source:  "hashicorp/tls",
					Version: "3.0.0",
				},
			},
		},
	}

	file, _ := json.MarshalIndent(providerTemplate, "", " ")
	return file
}

func acmeProviderConfigTemplate() []byte {
	var providerConfigTemplate hosts.ProviderConfigTemplate = hosts.ProviderConfigTemplate{
		Provider: map[string]interface{}{
			"acme": map[string]interface{}{
				"server_url": "https://acme-staging-v02.api.letsencrypt.org/directory",
			},
		},
	}

	file, _ := json.MarshalIndent(providerConfigTemplate, "", " ")
	return file
}

func acmeRegistrationTemplate(registrationEmail string) []byte {
	var acmeRegistrationTemplate map[string]interface{} = map[string]interface{}{
		"resource": map[string]interface{}{
			"tls_private_key": map[string]interface{}{
				"private_key": map[string]interface{}{
					"algorithm": "RSA",
				},
			},

			"acme_registration": map[string]interface{}{
				"reg": map[string]interface{}{
					"account_key_pem": "${tls_private_key.private_key.private_key_pem}",
					"email_address":   registrationEmail,
				},
			},
		},
	}

	file, _ := json.MarshalIndent(acmeRegistrationTemplate, "", " ")
	return file
}
