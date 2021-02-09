package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	var alias string = assignAliasName()
	hostProvider := SupportedProviders.Providers["host"].(HostProvider)
	host := hostProvider.Supported[name]
	host.ConfigureAuth()
	providerTemplatePath := filepath.Join(cliinit.HostAliasPath(name, alias), "provider.tf.json")
	// This doesn't work with multiple aliases since
	// provider config file is created only once on host dir configuration
	providerConfigTemplatePath := filepath.Join(cliinit.HostAliasPath(name, alias), "providerconfig.tf.json")
	moduleTemplatePath := filepath.Join(cliinit.HostPath(name), alias+".tf.json")
	if !HostDirectoryConfigured(cliinit.HostAliasPath(name, alias)) {
		channel <- fmt.Sprint("Configuring ", name, " host...")
		err := InstallTerraformProvider(alias, cliinit.HostPath(name), cliinit.HostAliasPath(name, alias), host, providerTemplatePath, providerConfigTemplatePath, moduleTemplatePath)
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

func HostDirectoryConfigured(hostPath string) bool {
	exists := true
	_, err := os.Stat(hostPath)
	if err != nil {
		return !exists
	}
	return exists
}

func InstallTerraformProvider(alias string, hostPath string, hostAliasPath string, host IHost, providerTemplatePath string, providerConfigTemplatePath string, moduleTemplatePath string) error {
	hostDirErr := os.MkdirAll(hostAliasPath, os.ModePerm)
	if hostDirErr != nil {
		os.Remove(hostAliasPath)
		log.Fatalln("error creating host config directory for", hostAliasPath, hostDirErr)
		return hostDirErr
	}

	moduleTemplatePathErr := ioutil.WriteFile(moduleTemplatePath, moduleTemplate(alias), 0644)
	providerTemplatePathErr := ioutil.WriteFile(providerTemplatePath, host.ProviderTemplate(), 0644)
	providerConfigTemplatePathErr := ioutil.WriteFile(providerConfigTemplatePath, host.ProviderConfigTemplate(), 0644)

	if moduleTemplatePathErr != nil || providerTemplatePathErr != nil || providerConfigTemplatePathErr != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(hostAliasPath)
		fmt.Println("failed ioutil.WriteFile for provider template")
		return fmt.Errorf("failed ioutil.WriteFile for provider template")
	}

	err := providers.TfInit(hostPath)
	if err != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(hostAliasPath)
		fmt.Println("failed init on new terraform directory", hostPath)
		return err
	}

	return nil
}

func moduleTemplate(alias string) []byte {

	var awsProviderTemplate providers.ModuleTemplate = providers.ModuleTemplate{
		Module: map[string]interface{}{
			alias: map[string]interface{}{
				"source": "./" + alias,
			},
		},
		Output: map[string]interface{}{
			alias + "_bucket": map[string]interface{}{
				"value": "${module.alias1.bucket}",
			},
			alias + "_bucket_regional_domain_name": map[string]interface{}{
				"value": "${module.alias1.bucket_regional_domain_name}",
			},
		},
	}

	file, _ := json.MarshalIndent(awsProviderTemplate, "", " ")
	return file
}

func assignAliasName() string {
	var alias string
	for {
		valid := true
		fmt.Print("Give your host an alias: ")
		fmt.Scanln(&alias)
		for _, hostName := range SupportedHosts {
			if hostName == alias {
				valid = !valid
				// TODO: WHAT IF AN ALIAS IS ADDED THAT BECOMES AN IVALID NAME IN THE FUTURE?
				// for example cli currently does not support firebase host so user can use 'firebase'
				// as their alias. Once firebase is supported this alias may become unsupported.

				// TODO: We'll need to provide a mechanism to rename an alias
				fmt.Println("alias cannot be the same as the name of a host provider (" + strings.Join(SupportedHosts[:], ", ") + ")")

				break
			}
		}

		if valid {
			break
		}
	}
	return alias

}
