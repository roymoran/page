package providers

import "fmt"

type IHost interface {
	Deploy() bool
	ConfigureHost() bool
}

func (hp HostProvider) Add(name string) (bool, string) {
	hostProvider := SupportedProviders.Providers["host"].(HostProvider)
	host := hostProvider.Supported[name]
	host.ConfigureHost()
	return true, fmt.Sprintln()
}

func (hp HostProvider) List(name string) (bool, string) {
	supportedHosts := fmt.Sprintln("Supported hosts")

	hostProvider := SupportedProviders.Providers["host"].(HostProvider)

	for hostName := range hostProvider.Supported {
		supportedHosts += fmt.Sprintln(hostName)
	}
	supportedHosts += fmt.Sprintln()
	return true, supportedHosts
}
