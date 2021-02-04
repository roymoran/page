package providers

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"builtonpage.com/main/cliinit"
	"github.com/hashicorp/terraform-exec/tfexec"
)

type IHost interface {
	ConfigureHost() bool
	AddHost(alias string, definitionPath string, statePath string) error
	ProviderTemplate() []byte
	ProviderConfigTemplate() []byte
}

func (hp HostProvider) Add(name string, channel chan string) error {
	// TODO Check if alias for host has already been added. if so return with
	// error
	alias := "alias"
	hostProvider := SupportedProviders.Providers["host"].(HostProvider)
	host := hostProvider.Supported[name]
	hostPath := filepath.Join(cliinit.TfInstallPath, name)
	providerTemplatePath := filepath.Join(hostPath, "provider.tf.json")
	providerConfigTemplatePath := filepath.Join(hostPath, alias+"_providerconfig.tf.json")
	stateDefinitionPath := filepath.Join(hostPath, alias+".tfstate")
	if !HostDirectoryConfigured(hostPath) {
		// TODO: This logic only allows for single host configured per host type e.g. aws, azure, etc.
		// please allow for user to configure multilple aws hosts (so there would be a single provider
		// tf file per directlry that contains provider details and the multiple tf files per host
		// configuration). This implies that each host config must be uniquely identified by a unique
		// alias. The alias will have be unique among all host configurations.
		// Once this is implemented the tf apply command would have to be run per host config so this
		// logic must be modified to
		channel <- fmt.Sprint("Applying ", name, " resource changes...")
		hostDirErr := os.MkdirAll(hostPath, os.ModePerm)
		if hostDirErr != nil {
			log.Fatalln("error creating host config directory for", hostPath, hostDirErr)
		}
		InstallTerraformProvider(alias, hostPath, host, providerTemplatePath, providerConfigTemplatePath, stateDefinitionPath)
	}

	// TODO: Get host alias from stdin
	channel <- fmt.Sprint("Adding ", name, " host configuration...")
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
	result := true
	_, err := os.Stat(hostPath)
	if err != nil {
		result = false
		return result
	}
	return result
}

func InstallTerraformProvider(providerId string, hostPath string, host IHost, providerTemplatePath string, providerConfigTemplatePath string, stateDefinitionPath string) {
	err := ioutil.WriteFile(providerTemplatePath, host.ProviderTemplate(), 0644)
	err = ioutil.WriteFile(providerConfigTemplatePath, host.ProviderConfigTemplate(), 0644)
	if err != nil {
		fmt.Println("failed ioutil.WriteFile for provider template", err)
	}

	tf, err := tfexec.NewTerraform(hostPath, cliinit.TfExecPath)
	if err != nil {
		fmt.Println(tf.Output(context.Background()))
		fmt.Println("error creating NewTerraform", hostPath, cliinit.TfInstallPath)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))

	if err != nil {
		fmt.Println(tf.Output(context.Background()))
		fmt.Println("error initializing tf directory", hostPath, cliinit.TfInstallPath, err)
	}
}
