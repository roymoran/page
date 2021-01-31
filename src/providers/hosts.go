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
	Deploy() bool
	ConfigureHost(alias string, definitionPath string, statePath string) (bool, error)
	HostProviderDefinition() []byte
}

func (hp HostProvider) Add(name string) (bool, string) {
	fmt.Println("Add Host Provider in hosts.go")
	alias := "alias"
	hostProvider := SupportedProviders.Providers["host"].(HostProvider)
	host := hostProvider.Supported[name]
	hostPath := filepath.Join(cliinit.TfInstallPath, name)
	definitionPath := filepath.Join(hostPath, alias+".tf.json")
	stateDefinitionPath := filepath.Join(hostPath, alias+".tfstate")
	if !HostDirectoryConfigured(hostPath) {
		fmt.Println("1 time Directory Config in hosts.go")
		hostDirErr := os.MkdirAll(hostPath, os.ModePerm)
		if hostDirErr != nil {
			log.Fatalln("error creating host config directory for", hostPath, hostDirErr)
		}
		InstallTerraformPlugin(alias, hostPath, host, definitionPath, stateDefinitionPath)
		fmt.Println("finish 1 time Directory Config in hosts.go")
	}

	// TODO: Get host alias from stdin
	host.ConfigureHost(alias, definitionPath, stateDefinitionPath)

	return true, fmt.Sprintln()
}

func (hp HostProvider) List(name string) (bool, string) {
	supportedHosts := fmt.Sprint()
	for _, hostName := range SupportedHosts {
		supportedHosts += fmt.Sprintln(hostName)
	}
	supportedHosts += fmt.Sprintln()
	return true, supportedHosts
}

func HostDirectoryConfigured(hostPath string) bool {
	fmt.Println("HostDirectoryConfigured check in hosts.go")
	result := true
	_, err := os.Stat(hostPath)
	if err != nil {
		result = false
		return result
	}
	fmt.Println("Finish HostDirectoryConfigured check in hosts.go")
	return result
}

func InstallTerraformPlugin(providerId string, hostPath string, host IHost, definitionPath string, stateDefinitionPath string) {
	fmt.Println("in InstallTerraformPlugin")
	_ = ioutil.WriteFile(definitionPath, host.HostProviderDefinition(), 0644)
	tf, err := tfexec.NewTerraform(hostPath, cliinit.TfExecPath)
	if err != nil {
		log.Fatalln("error creating NewTerraform", hostPath, cliinit.TfInstallPath)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))

	if err != nil {
		fmt.Println(tf.Output(context.Background()))
		log.Fatalln("error initializing tf directory", hostPath, cliinit.TfInstallPath, err)
	}

	tf.Apply(context.Background(), tfexec.State(stateDefinitionPath))
	fmt.Println("finish InstallTerraformPlugin")
}
