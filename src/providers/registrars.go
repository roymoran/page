package providers

import "fmt"

type IRegistrar interface {
	RegisterDomain() bool
	ConfigureDns() bool
	ConfigureRegistrar() bool
}

func (rp RegistrarProvider) Add(name string, channel chan string) (bool, string) {
	registrarProvider := SupportedProviders.Providers["registrar"].(RegistrarProvider)
	registrar := registrarProvider.Supported[name]
	registrar.ConfigureRegistrar()
	return true, fmt.Sprintln()
}

func (rp RegistrarProvider) List(name string, channel chan string) (bool, string) {
	supportedRegistrars := fmt.Sprint()
	for _, registrarName := range SupportedRegistrars {
		supportedRegistrars += fmt.Sprintln(registrarName)
	}
	supportedRegistrars += fmt.Sprintln()
	return true, supportedRegistrars
}
