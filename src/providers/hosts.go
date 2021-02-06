package providers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"builtonpage.com/main/cliinit"
	providers "builtonpage.com/main/providers/hosts"
)

type IHost interface {
	ConfigureAuth() error
	ConfigureHost(alias string) error
	AddHost(alias string, definitionPath string, statePath string) error
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
	hostPath := cliinit.HostPath(name)
	providerTemplatePath := filepath.Join(hostPath, "provider.tf.json")
	// This doesn't work with multiple aliases since
	// provider config file is created only once on host dir configuration
	providerConfigTemplatePath := filepath.Join(hostPath, alias+"_providerconfig.tf.json")
	// TODO: Should each alias have its own tf state file?
	// TODO: Consider leaving state file path as default name
	stateDefinitionPath := filepath.Join(hostPath, alias+".tfstate")
	if !HostDirectoryConfigured(hostPath) {
		channel <- fmt.Sprint("Configuring ", name, " host...")
		hostDirErr := os.MkdirAll(hostPath, os.ModePerm)
		if hostDirErr != nil {
			log.Fatalln("error creating host config directory for", hostPath, hostDirErr)
		}
		InstallTerraformProvider(alias, hostPath, host, providerTemplatePath, providerConfigTemplatePath, stateDefinitionPath)
	}

	// TODO: Get host alias from stdin
	channel <- fmt.Sprint("Saving ", name, " host configuration...")
	host.AddHost(alias, providerTemplatePath, stateDefinitionPath)

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

func InstallTerraformProvider(providerId string, hostPath string, host IHost, providerTemplatePath string, providerConfigTemplatePath string, stateDefinitionPath string) {
	err := ioutil.WriteFile(providerTemplatePath, host.ProviderTemplate(), 0644)
	err = ioutil.WriteFile(providerConfigTemplatePath, host.ProviderConfigTemplate(), 0644)
	if err != nil {
		fmt.Println("failed ioutil.WriteFile for provider template", err)
	}

	err = providers.TfInit(hostPath)
	if err != nil {
		os.Remove(providerTemplatePath)
		os.Remove(providerConfigTemplatePath)
		fmt.Println("failed init on new terraform directory", hostPath)
	}
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
