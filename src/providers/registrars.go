package providers

type IRegistrar interface {
	RegisterDomain() bool
	ConfigureDns() bool
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
