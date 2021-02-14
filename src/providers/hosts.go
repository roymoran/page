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
	providers "builtonpage.com/main/providers/hosts"
)

type IHost interface {
	ConfigureAuth() error
	ConfigureHost(alias string, templatePath string, page definition.PageDefinition) error
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

	moduleTemplatePath := filepath.Join(cliinit.ProvidersPath, name+"_"+alias+".tf.json")
	if !AliasDirectoryConfigured(cliinit.ProviderAliasPath(name, alias)) {
		channel <- fmt.Sprint("Configuring ", name, " host...")
		err := InstallHostTerraformProvider(name, alias, cliinit.ProviderAliasPath(name, alias), host, providerTemplatePath, providerConfigTemplatePath, moduleTemplatePath)
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

func InstallHostTerraformProvider(name string, alias string, providerAliasPath string, host IHost, providerTemplatePath string, providerConfigTemplatePath string, moduleTemplatePath string) error {
	hostDirErr := os.MkdirAll(providerAliasPath, os.ModePerm)
	if hostDirErr != nil {
		os.Remove(providerAliasPath)
		log.Fatalln("error creating host config directory for", providerAliasPath, hostDirErr)
		return hostDirErr
	}

	moduleTemplatePathErr := ioutil.WriteFile(moduleTemplatePath, hostModuleTemplate(name, alias), 0644)
	providerTemplatePathErr := ioutil.WriteFile(providerTemplatePath, host.ProviderTemplate(), 0644)
	providerConfigTemplatePathErr := ioutil.WriteFile(providerConfigTemplatePath, host.ProviderConfigTemplate(), 0644)

	if moduleTemplatePathErr != nil || providerTemplatePathErr != nil || providerConfigTemplatePathErr != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(providerAliasPath)
		fmt.Println("failed ioutil.WriteFile for provider template")
		return fmt.Errorf("failed ioutil.WriteFile for provider template")
	}

	err := providers.TfInit(cliinit.ProvidersPath)
	if err != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(providerAliasPath)
		fmt.Println("failed init on new terraform directory", cliinit.ProvidersPath)
		return err
	}

	return nil
}

func hostModuleTemplate(providerName string, alias string) []byte {

	var awsProviderTemplate providers.ModuleTemplate = providers.ModuleTemplate{
		Module: map[string]interface{}{
			"host_" + alias: map[string]interface{}{
				"source": "./" + providerName + "/" + alias,
			},
		},
		Output: map[string]interface{}{
			alias + "_bucket": map[string]interface{}{
				"value": "${module.host_" + alias + ".bucket}",
			},
			alias + "_bucket_regional_domain_name": map[string]interface{}{
				"value": "${module.host_" + alias + ".bucket_regional_domain_name}",
			},
		},
	}

	file, _ := json.MarshalIndent(awsProviderTemplate, "", " ")
	return file
}
