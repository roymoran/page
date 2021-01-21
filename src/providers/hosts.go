package providers

import "fmt"

type IHost interface {
	Deploy() bool
	ConfigureHost() bool
}

func (hp HostProvider) Add(name string) bool {
	hostProvider := SupportedProviders.Providers["host"].(HostProvider)
	host := hostProvider.Supported[name]
	host.ConfigureHost()
	return true
}

func (hp HostProvider) List(name string) bool {
	fmt.Println("host list")
	return true
}
