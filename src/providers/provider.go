package providers

import (
	hosts "builtonpage.com/main/providers/hosts"
	registrars "builtonpage.com/main/providers/registrars"
)

type IProvider interface {
	Add() bool
	Remove() bool
	List() bool
}

type RegistrarProvider struct {
	SupportedRegistrars map[string]IRegistrar
}

type HostProvider struct {
	SupportedHosts map[string]IHost
}

var provider = map[string]IProvider{
	"host": HostProvider{
		SupportedHosts: map[string]IHost{
			"page": hosts.PageHost{},
		},
	},
	"registrar": RegistrarProvider{
		SupportedRegistrars: map[string]IRegistrar{
			"namecheap": registrars.Namecheap{},
			"page":      registrars.Page{},
		},
	},
}
