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
	ConfigureHost() bool
	HostProviderDefinition() string
}

func (hp HostProvider) Add(name string) (bool, string) {
	hostProvider := SupportedProviders.Providers["host"].(HostProvider)
	host := hostProvider.Supported[name]

	if !HostDirectoryConfigured(name) {
		ConfigureHostDirectory(name, host)
	}

	host.ConfigureHost()
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

func ConfigureHostDirectory(hostName string, host IHost) {
	hostDirErr := os.MkdirAll(filepath.Join(cliinit.TfInstallPath, hostName), os.ModePerm)

	if hostDirErr != nil {
		log.Fatalln("error creating host config directory for", hostName, hostDirErr)
	}

	InstallTerraformPlugin(hostName, host)
}

func InstallTerraformPlugin(hostName string, host IHost) {
	hostPath := filepath.Join(cliinit.TfInstallPath, hostName)
	hostDirErr := os.MkdirAll(hostPath, os.ModePerm)

	if hostDirErr != nil {
		log.Fatalln("error creating host config directory for", hostName, hostDirErr)
	}

	fmt.Println(cliinit.TfInstallPath)
	_ = ioutil.WriteFile(filepath.Join(hostPath, hostName+".tf"), []byte(host.HostProviderDefinition()), 0644)
	tf, err := tfexec.NewTerraform(hostPath, cliinit.TfExecPath)
	if err != nil {
		log.Fatalln("error creating NewTerraform", hostPath, cliinit.TfInstallPath)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))

	if err != nil {
		log.Fatalln("error initializing tf directory", hostPath, cliinit.TfInstallPath)
	}

}

func HostDirectoryConfigured(hostName string) bool {
	result := true
	_, err := os.Stat(filepath.Join(cliinit.TfInstallPath, hostName))
	if err != nil {
		result = false
		return result
	}
	return result
}
