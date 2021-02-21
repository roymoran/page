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

type IHost interface {
	ConfigureAuth() error
	ConfigureHost(hostAlias string, templatePath string, page definition.PageDefinition) error
	AddHost(alias string, definitionPath string) error
	ProviderTemplate() []byte
	ProviderConfigTemplate() []byte
}

func (hp HostProvider) Add(name string, channel chan string) error {
	// TODO Check if alias for host has already been added. if so return with
	// error
	var alias string = AssignAliasName("host")
	hostProvider := SupportedProviders.Providers["host"].(HostProvider)
	host := hostProvider.Supported[name]
	host.ConfigureAuth()
	providerTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "provider.tf.json")
	// This doesn't work with multiple aliases since
	// provider config file is created only once on host dir configuration
	providerConfigTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "providerconfig.tf.json")
	certificatesVariableTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "certificatesvar.tf.json")

	moduleTemplatePath := cliinit.ModuleTemplatePath("host", alias)
	if !AliasDirectoryConfigured(cliinit.ProviderAliasPath(name, alias)) {
		channel <- fmt.Sprint("Configuring ", name, " host...")
		err := InstallHostTerraformProvider(name, alias, cliinit.ProviderAliasPath(name, alias), host, providerTemplatePath, providerConfigTemplatePath, moduleTemplatePath, certificatesVariableTemplatePath)
		if err != nil {
			return err
		}
	}

	channel <- fmt.Sprint("Saving ", name, " host configuration...")
	host.AddHost(alias, providerTemplatePath)

	return nil
}

func (hp HostProvider) List(name string, channel chan string) error {
	for _, hostName := range SupportedHosts {
		channel <- fmt.Sprint(hostName)
	}
	return nil
}

func InstallHostTerraformProvider(name string, alias string, providerAliasPath string, host IHost, providerTemplatePath string, providerConfigTemplatePath string, moduleTemplatePath string, certificatesVariableTemplatePath string) error {
	hostDirErr := os.MkdirAll(providerAliasPath, os.ModePerm)
	if hostDirErr != nil {
		os.Remove(providerAliasPath)
		log.Fatalln("error creating host config directory for", providerAliasPath, hostDirErr)
		return hostDirErr
	}

	moduleTemplatePathErr := ioutil.WriteFile(moduleTemplatePath, hostModuleTemplate(name, alias), 0644)
	providerTemplatePathErr := ioutil.WriteFile(providerTemplatePath, host.ProviderTemplate(), 0644)
	providerConfigTemplatePathErr := ioutil.WriteFile(providerConfigTemplatePath, host.ProviderConfigTemplate(), 0644)
	hostCertificatesVariableTemplatePathErr := ioutil.WriteFile(certificatesVariableTemplatePath, hostCertificatesVariableTemplate(alias), 0644)

	if moduleTemplatePathErr != nil ||
		providerTemplatePathErr != nil ||
		providerConfigTemplatePathErr != nil ||
		hostCertificatesVariableTemplatePathErr != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(providerAliasPath)
		fmt.Println("failed ioutil.WriteFile for provider template")
		return fmt.Errorf("failed ioutil.WriteFile for provider template")
	}

	err := hosts.TfInit(cliinit.ProvidersPath)
	if err != nil {
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
