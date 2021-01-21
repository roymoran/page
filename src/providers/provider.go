package providers

import (
	hosts "builtonpage.com/main/providers/hosts"
	registrars "builtonpage.com/main/providers/registrars"
)

type Provider struct {
	Actions            map[string]func(IProvider, string) (bool, string)
	Providers          map[string]IProvider
	HostProviders      map[string]HostProvider
	RegistrarProviders map[string]RegistrarProvider
}

type IProvider interface {
	Add(string) (bool, string)
	List(string) (bool, string)
}

type RegistrarProvider struct {
	Supported map[string]IRegistrar
}

type HostProvider struct {
	Supported map[string]IHost
}

var SupportedProviders = Provider{
	Actions: map[string]func(IProvider, string) (bool, string){"add": IProvider.Add, "list": IProvider.List},
	Providers: map[string]IProvider{
		"host": HostProvider{
			Supported: map[string]IHost{
				"page": hosts.PageHost{},
			},
		},
		"registrar": RegistrarProvider{
			Supported: map[string]IRegistrar{
				"namecheap": registrars.Namecheap{},
				"page":      registrars.Page{},
			},
		},
	},
}
