package providers

import (
	"fmt"

	"builtonpage.com/main/definition"
)

type IRegistrar interface {
	ConfigureRegistrar(definition.PageDefinition) bool
	ConfigureDns() bool
	AddRegistrar(string) error
}

func (rp RegistrarProvider) Add(name string, channel chan string) error {
	registrarProvider := SupportedProviders.Providers["registrar"].(RegistrarProvider)
	registrar := registrarProvider.Supported[name]
	registrar.AddRegistrar("namecheap_main")
	return nil
}

func (rp RegistrarProvider) List(name string, channel chan string) error {
	for _, registrarName := range SupportedRegistrars {
		channel <- fmt.Sprint(registrarName)
	}
	return nil
}
