package providers

import (
	"fmt"
	"sort"

	"builtonpage.com/main/cliinit"
	hosts "builtonpage.com/main/providers/hosts"
	registrars "builtonpage.com/main/providers/registrars"
)

type Provider struct {
	Actions   map[string]func(IProvider, string) (bool, string)
	Providers map[string]IProvider
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
	Actions: map[string]func(IProvider, string) (bool, string){"add": AddProvider, "list": IProvider.List},
	Providers: map[string]IProvider{
		"host": HostProvider{
			Supported: map[string]IHost{
				"page":  hosts.PageHost{},
				"azure": hosts.Azure{},
				"aws":   hosts.AmazonWebServices{},
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

func AddProvider(provider IProvider, providerName string) (bool, string) {
	fmt.Println("cliinit.CliInitialized()", cliinit.CliInitialized())
	if !cliinit.CliInitialized() {
		cliinit.CliInit()
	}

	// TODO: Change
	_, _ = provider.Add(providerName)

	return true, ""
}

var SupportedProviderTypes []string = BuildSupportedProviderTypes()
var SupportedAction []string = BuildSupportedActions()
var SupportedRegistrars []string = BuildSupportedRegistrars()
var SupportedHosts []string = BuildSupportedHosts()

func BuildSupportedHosts() []string {
	hostMap := SupportedProviders.Providers["host"].(HostProvider).Supported
	hosts := make([]string, len(hostMap))
	index := 0
	for host := range hostMap {
		hosts[index] = fmt.Sprint(host)
		index++
	}

	sort.Strings(hosts)

	return hosts
}

func BuildSupportedRegistrars() []string {
	registrarMap := SupportedProviders.Providers["registrar"].(RegistrarProvider).Supported
	registrars := make([]string, len(registrarMap))
	index := 0
	for registrar := range registrarMap {
		registrars[index] = fmt.Sprint(registrar)
		index++
	}

	sort.Strings(registrars)

	return registrars
}

func BuildSupportedActions() []string {
	actions := make([]string, len(SupportedProviders.Actions))
	index := 0
	for action := range SupportedProviders.Actions {
		actions[index] = fmt.Sprint(action)
		index++
	}

	sort.Strings(actions)

	return actions
}

func BuildSupportedProviderTypes() []string {
	providerTypes := make([]string, len(SupportedProviders.Providers))
	index := 0
	for providerType := range SupportedProviders.Providers {
		providerTypes[index] = fmt.Sprint(providerType)
		index++
	}

	sort.Strings(providerTypes)

	return providerTypes
}
