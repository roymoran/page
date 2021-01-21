package providers

import "fmt"

type IRegistrar interface {
	RegisterDomain() bool
	ConfigureDns() bool
	ConfigureRegistrar() bool
}

func (rp RegistrarProvider) Add(name string) bool {
	registrarProvider := SupportedProviders.Providers["registrar"].(RegistrarProvider)
	registrar := registrarProvider.Supported[name]
	registrar.ConfigureRegistrar()
	return true
}

func (rp RegistrarProvider) List(name string) bool {
	fmt.Println("registrar list")
	return true
}
