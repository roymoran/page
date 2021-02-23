package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/constants"
	"builtonpage.com/main/definition"
	"builtonpage.com/main/providers/hosts"
)

type IRegistrar interface {
	ConfigureAuth() (cliinit.Credentials, error)
	ConfigureRegistrar(string, definition.PageDefinition) error
	AddRegistrar(string, cliinit.Credentials) error
	ProviderDefinition() (string, hosts.Provider)
	ProviderConfig(string, string) map[string]string
}

func (rp RegistrarProvider) Add(name string, channel chan string) error {
	var alias string = AssignAliasName("registrar")
	var acmeEmail string

	registrarProvider := SupportedProviders.Providers["registrar"].(RegistrarProvider)
	registrar := registrarProvider.Supported[name]

	fmt.Print("ACME registration email: ")
	_, err := fmt.Scanln(&acmeEmail)
	if err != nil {
		return err
	}

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
	acmeRegistrationTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "acmeregistration.tf.json")
	domainsVariableTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "domainsvar.tf.json")

	moduleTemplatePath := cliinit.ModuleTemplatePath("registrar", alias)

	if !AliasDirectoryConfigured(cliinit.ProviderAliasPath(name, alias)) {
		channel <- fmt.Sprint("Configuring ", name, " registrar...")
		err := InstallRegistrarTerraformProvider(name, alias, cliinit.ProviderAliasPath(name, alias), providerTemplatePath, providerConfigTemplatePath, moduleTemplatePath, acmeRegistrationTemplatePath, acmeEmail, registrarProviderName, registrarProviderDefinition, registrarProviderConfig, domainsVariableTemplatePath)
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
	channel <- fmt.Sprint("Supported\n")
	for _, registrarName := range SupportedRegistrars {
		channel <- fmt.Sprintln(registrarName)
	}
	channel <- fmt.Sprintln("")
	channel <- fmt.Sprint("Configured\n")
	configuredAliases, _ := cliinit.FindAllRegistrarAliases()
	for _, alias := range configuredAliases {
		channel <- fmt.Sprintln(alias)
	}

	return nil
}

func InstallRegistrarTerraformProvider(name string, alias string, providerAliasPath string, providerTemplatePath string, providerConfigTemplatePath string, moduleTemplatePath string, acmeRegistrationTemplatePath string, acmeRegistrationEmail string, registrarProviderName string, registrarProviderDefinition hosts.Provider, registrarProviderConfig map[string]string, domainsVariableTemplatePath string) error {
	hostDirErr := os.MkdirAll(providerAliasPath, os.ModePerm)
	if hostDirErr != nil {
		os.Remove(providerAliasPath)
		log.Fatalln("error creating host config directory for", providerAliasPath, hostDirErr)
		return hostDirErr
	}

	registrarProviderTemplate := registrarProviderTemplate(registrarProviderName, registrarProviderDefinition)
	moduleTemplatePathErr := ioutil.WriteFile(moduleTemplatePath, registrarModuleTemplate(name, alias), 0644)
	providerTemplatePathErr := ioutil.WriteFile(providerTemplatePath, registrarProviderTemplate, 0644)
	providerConfigTemplatePathErr := ioutil.WriteFile(providerConfigTemplatePath, registrarProviderConfigTemplate(registrarProviderName, registrarProviderConfig), 0644)
	acmeRegistrationTemplatePathErr := ioutil.WriteFile(acmeRegistrationTemplatePath, acmeRegistrationTemplate(acmeRegistrationEmail), 0644)
	domainsVariableTemplatePathErr := ioutil.WriteFile(domainsVariableTemplatePath, hostDomainsVariableTemplate(), 0644)

	if moduleTemplatePathErr != nil || providerTemplatePathErr != nil || providerConfigTemplatePathErr != nil || acmeRegistrationTemplatePathErr != nil || domainsVariableTemplatePathErr != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(providerAliasPath)
		fmt.Println("failed ioutil.WriteFile for provider template")
		return fmt.Errorf("failed ioutil.WriteFile for provider template")
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

	var awsProviderTemplate hosts.ModuleTemplate = hosts.ModuleTemplate{
		Module: map[string]interface{}{
			"registrar_" + alias: map[string]interface{}{
				"source":  "./" + providerName + "/" + alias,
				"domains": map[string]string{},
			},
		},
	}

	file, _ := json.MarshalIndent(awsProviderTemplate, "", " ")
	return file
}

func registrarProviderTemplate(registrarProviderName string, registrarProviderDefinition hosts.Provider) []byte {
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
			"acme": map[string]interface{}{
				"server_url": constants.AcmeServerUrl(),
			},
			registrarProviderName: registrarProviderConfig,
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

func hostDomainsVariableTemplate() []byte {

	var hostDomainsVariableTemplate map[string]interface{} = map[string]interface{}{
		"variable": map[string]interface{}{
			"domains": map[string]interface{}{
				"type": "map(object({domain = string}))",
			},
		},
	}

	file, _ := json.MarshalIndent(hostDomainsVariableTemplate, "", " ")
	return file
}
