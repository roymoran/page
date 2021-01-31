package providers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"builtonpage.com/main/cliinit"
)

type IHost interface {
	Deploy() bool
	ConfigureHost(alias string) (bool, error)
	HostProviderDefinition() []byte
}

func (hp HostProvider) Add(name string) (bool, string) {
	hostProvider := SupportedProviders.Providers["host"].(HostProvider)
	host := hostProvider.Supported[name]
	hostPath := filepath.Join(cliinit.TfInstallPath, name)
	if !HostDirectoryConfigured(hostPath) {
		hostDirErr := os.MkdirAll(hostPath, os.ModePerm)

		if hostDirErr != nil {
			log.Fatalln("error creating host config directory for", hostPath, hostDirErr)
		}
	}

	host.ConfigureHost("alias2")
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
	result := true
	_, err := os.Stat(hostPath)
	if err != nil {
		result = false
		return result
	}
	return result
}
