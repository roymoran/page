package providers

import "fmt"

type IRegistrar interface {
	RegisterDomain() bool
	ConfigureDns() bool
	ConfigureRegistrar() bool
}

func (rp RegistrarProvider) Add(name string) (bool, string) {
	registrarProvider := SupportedProviders.Providers["registrar"].(RegistrarProvider)
	registrar := registrarProvider.Supported[name]
	registrar.ConfigureRegistrar()
	return true, fmt.Sprintln()
}

func (rp RegistrarProvider) List(name string) (bool, string) {
	supportedRegistrars := fmt.Sprintln("Supported registrars")
	registrarProvider := SupportedProviders.Providers["registrar"].(RegistrarProvider)

	for registrarName := range registrarProvider.Supported {
		supportedRegistrars += fmt.Sprintln(registrarName)
	}
	supportedRegistrars += fmt.Sprintln()

	return true, supportedRegistrars
}
