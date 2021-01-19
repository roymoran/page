package providers

import (
	registrars "builtonpage.com/main/providers/registrars"
)

type IRegistrar interface {
	RegisterDomain() bool
	ConfigureDns() bool
}

type RegistrarProvider struct {
	SupportedRegistrars map[string]IRegistrar
}

var registrarProvider RegistrarProvider = RegistrarProvider{
	SupportedRegistrars: map[string]IRegistrar{
		"namecheap": registrars.Namecheap{},
		"page":      registrars.Page{},
	},
}

func (rp RegistrarProvider) Add() bool {
	return true
}

func (rp RegistrarProvider) Remove() bool {
	return true
}

func (rp RegistrarProvider) List() bool {
	return true
}
