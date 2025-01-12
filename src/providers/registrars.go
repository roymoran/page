package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"pagecli.com/main/cliinit"
	"pagecli.com/main/definition"
	"pagecli.com/main/providers/hosts"
)

type IRegistrar interface {
	ConfigureAuth() (cliinit.Credentials, error)
	ConfigureCertificate(string, definition.PageDefinition) error
	ConfigureDNS(string, definition.PageDefinition) error
	AddRegistrar(string, cliinit.Credentials) error
	ProviderDefinition() (string, hosts.Provider)
	ProviderConfig(string, string) map[string]string
}

func (rp RegistrarProvider) Add(name string, channel chan string) error {
	var alias string = AssignAliasName("registrar")

	registrarProvider := SupportedProviders.Providers["registrar"].(RegistrarProvider)
	registrar := registrarProvider.Supported[name]

	credentials, authErr := registrar.ConfigureAuth()
	registrarProviderName, registrarProviderDefinition := registrar.ProviderDefinition()
	registrarProviderConfig := registrar.ProviderConfig(credentials.Username, credentials.Password)

	if authErr != nil {
		return authErr
	}

	providerTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "provider.tf.json")
	// This doesn't work with multiple aliases since
	// provider config file is created only once on host dir configuration
	providerConfigTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "providerconfig.tf.json")
	// acmeRegistrationTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "acmeregistration.tf.json")
	registrarVariablesTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "registrarvar.tf.json")

	moduleTemplatePath := cliinit.ModuleTemplatePath("registrar", alias)

	if !AliasDirectoryConfigured(cliinit.ProviderAliasPath(name, alias)) {
		channel <- fmt.Sprintln("Configuring", name, "registrar...")
		err := InstallRegistrarTerraformProvider(name, alias, cliinit.ProviderAliasPath(name, alias), providerTemplatePath, providerConfigTemplatePath, moduleTemplatePath, registrarProviderName, registrarProviderDefinition, registrarProviderConfig, registrarVariablesTemplatePath)
		if err != nil {
			return err
		}
	}

	channel <- fmt.Sprintln("Saving", name, "registrar configuration...")
	addRegistrarErr := registrar.AddRegistrar(alias, credentials)
	if addRegistrarErr != nil {
		return addRegistrarErr
	}

	return nil
}

func (rp RegistrarProvider) List(name string, channel chan string) error {
	channel <- fmt.Sprint("Supported Registrars\n")
	for _, registrarName := range SupportedRegistrars {
		channel <- fmt.Sprintln(registrarName)
	}
	channel <- fmt.Sprintln("")
	channel <- fmt.Sprint("Configured Registrars\n")
	configuredAliases, _ := cliinit.FindAllRegistrarAliases()
	for _, alias := range configuredAliases {
		channel <- fmt.Sprintln(alias)
	}

	return nil
}

func InstallRegistrarTerraformProvider(name string, alias string, providerAliasPath string, providerTemplatePath string, providerConfigTemplatePath string, moduleTemplatePath string, registrarProviderName string, registrarProviderDefinition hosts.Provider, registrarProviderConfig map[string]string, registrarVariablesTemplatePath string) error {
	hostDirErr := os.MkdirAll(providerAliasPath, os.ModePerm)
	if hostDirErr != nil {
		os.Remove(providerAliasPath)
		log.Fatalln("error creating host config directory for", providerAliasPath, hostDirErr)
		return hostDirErr
	}

	registrarProviderTemplate := registrarProviderTemplate(registrarProviderName, registrarProviderDefinition)
	moduleTemplatePathErr := os.WriteFile(moduleTemplatePath, registrarModuleTemplate(name, alias), 0644)
	providerTemplatePathErr := os.WriteFile(providerTemplatePath, registrarProviderTemplate, 0644)
	providerConfigTemplatePathErr := os.WriteFile(providerConfigTemplatePath, registrarProviderConfigTemplate(registrarProviderName, registrarProviderConfig), 0644)
	registrarVariablesTemplatePathErr := os.WriteFile(registrarVariablesTemplatePath, hostRegistrarVariablesTemplate(), 0644)

	if moduleTemplatePathErr != nil ||
		providerTemplatePathErr != nil ||
		providerConfigTemplatePathErr != nil ||
		registrarVariablesTemplatePathErr != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(providerAliasPath)
		fmt.Println("failed os.WriteFile for provider template")
		return fmt.Errorf("failed os.WriteFile for provider template")
	}

	err := hosts.TfInit(cliinit.ProvidersPath)
	if err != nil {
		// os.Remove(moduleTemplatePath)
		// os.RemoveAll(providerAliasPath)
		fmt.Println("failed init on new terraform directory", cliinit.ProvidersPath)
		return err
	}

	return nil
}

func registrarModuleTemplate(providerName string, alias string) []byte {

	var registrarProviderTemplate hosts.ModuleTemplate = hosts.ModuleTemplate{
		Module: map[string]interface{}{
			"registrar_" + alias: map[string]interface{}{
				"source": "./" + providerName + "/" + alias,
			},
		},
	}

	file, _ := json.MarshalIndent(registrarProviderTemplate, "", " ")
	return file
}

func registrarProviderTemplate(registrarProviderName string, registrarProviderDefinition hosts.Provider) []byte {
	var providerTemplate hosts.ProviderTemplate = hosts.ProviderTemplate{
		Terraform: hosts.RequiredProviders{
			RequiredProvider: map[string]hosts.Provider{
				registrarProviderName: registrarProviderDefinition,
			},
		},
	}

	file, _ := json.MarshalIndent(providerTemplate, "", " ")
	return file
}

func registrarProviderConfigTemplate(registrarProviderName string, registrarProviderConfig map[string]string) []byte {
	var providerConfigTemplate hosts.ProviderConfigTemplate = hosts.ProviderConfigTemplate{
		Provider: map[string]interface{}{
			registrarProviderName: registrarProviderConfig,
		},
	}

	file, _ := json.MarshalIndent(providerConfigTemplate, "", " ")
	return file
}

func hostRegistrarVariablesTemplate() []byte {

	var hostRegistrarVariablesTemplate map[string]interface{} = map[string]interface{}{
		"variable": map[string]interface{}{
			"domains": map[string]interface{}{
				"type":     "map(object({domain = string}))",
				"nullable": "true",
				"default":  nil,
			},
			"dns_records": map[string]interface{}{
				"type":     "map(object({records = list(object({ host = string, type = string, value = string, ttl = number})), record_vars = list(string)}))",
				"nullable": "true",
				"default":  nil,
			},
		},
	}

	file, _ := json.MarshalIndent(hostRegistrarVariablesTemplate, "", " ")
	return file
}
