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

type IHost interface {
	ConfigureAuth() error
	ConfigureHost(hostAlias string, templatePath string, page definition.PageDefinition) error
	ConfigureWebsite(hostAlias string, templatePath string, page definition.PageDefinition) error
	AddHost(alias string, definitionPath string) error
	ProviderTemplate() []byte
	ProviderConfigTemplate() []byte
}

func (hp HostProvider) Add(name string, channel chan string) error {
	var alias string = AssignAliasName("host")
	hostProvider := SupportedProviders.Providers["host"].(HostProvider)
	host := hostProvider.Supported[name]
	host.ConfigureAuth()
	providerTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "provider.tf.json")
	providerConfigTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "providerconfig.tf.json")
	certificatesVariableTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "certificatesvar.tf.json")

	moduleTemplatePath := cliinit.ModuleTemplatePath("host", alias)
	if !AliasDirectoryConfigured(cliinit.ProviderAliasPath(name, alias)) {
		channel <- fmt.Sprintln("Configuring", name, "host...")
		err := InstallHostTerraformProvider(name, alias, cliinit.ProviderAliasPath(name, alias), host, providerTemplatePath, providerConfigTemplatePath, moduleTemplatePath, certificatesVariableTemplatePath)
		if err != nil {
			return err
		}
	}

	channel <- fmt.Sprintln("Saving", name, "host configuration...")
	host.AddHost(alias, providerTemplatePath)

	return nil
}

func (hp HostProvider) List(name string, channel chan string) error {
	channel <- fmt.Sprint("Supported Hosts\n")

	for _, hostName := range SupportedHosts {
		channel <- fmt.Sprintln(hostName)
	}
	channel <- fmt.Sprintln("")
	channel <- fmt.Sprint("Configured Hosts\n")
	configuredAliases, _ := cliinit.FindAllHostAliases()
	for _, alias := range configuredAliases {
		channel <- fmt.Sprintln(alias)
	}

	return nil
}

func InstallHostTerraformProvider(name string, alias string, providerAliasPath string, host IHost, providerTemplatePath string, providerConfigTemplatePath string, moduleTemplatePath string, certificatesVariableTemplatePath string) error {
	hostDirErr := os.MkdirAll(providerAliasPath, os.ModePerm)
	if hostDirErr != nil {
		os.RemoveAll(providerAliasPath)
		log.Fatalln("error creating host config directory for", providerAliasPath, hostDirErr)
		return hostDirErr
	}

	moduleTemplatePathErr := os.WriteFile(moduleTemplatePath, hostModuleTemplate(name, alias), 0644)
	providerTemplatePathErr := os.WriteFile(providerTemplatePath, host.ProviderTemplate(), 0644)
	providerConfigTemplatePathErr := os.WriteFile(providerConfigTemplatePath, host.ProviderConfigTemplate(), 0644)
	hostCertificatesVariableTemplatePathErr := os.WriteFile(certificatesVariableTemplatePath, hostCertificatesVariableTemplate(alias), 0644)

	if moduleTemplatePathErr != nil ||
		providerTemplatePathErr != nil ||
		providerConfigTemplatePathErr != nil ||
		hostCertificatesVariableTemplatePathErr != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(providerAliasPath)
		return fmt.Errorf("failed os.WriteFile for provider template")
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

func hostModuleTemplate(providerName string, alias string) []byte {

	var hostProviderTemplate map[string]interface{} = map[string]interface{}{
		"module": map[string]interface{}{
			"host_" + alias: map[string]interface{}{
				"source":       "./" + providerName + "/" + alias,
				"certificates": map[string]string{},
			},
		},
	}

	file, _ := json.MarshalIndent(hostProviderTemplate, "", " ")
	return file
}

func hostCertificatesVariableTemplate(alias string) []byte {

	var hostCertificatesVariableTemplate map[string]interface{} = map[string]interface{}{
		"variable": map[string]interface{}{
			"certificates": map[string]interface{}{
				"type": "map(object({certificate_pem = string, private_key_pem = string, certificate_chain = string}))",
			},
		},
	}

	file, _ := json.MarshalIndent(hostCertificatesVariableTemplate, "", " ")
	return file
}
