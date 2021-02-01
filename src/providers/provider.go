package providers

import (
	"fmt"
	"sort"

	"builtonpage.com/main/cliinit"
	hosts "builtonpage.com/main/providers/hosts"
	registrars "builtonpage.com/main/providers/registrars"
)

type Provider struct {
	Actions   map[string]func(IProvider, string, chan string) (bool, string)
	Providers map[string]IProvider
}

type IProvider interface {
	Add(string, chan string) (bool, string)
	List(string, chan string) (bool, string)
}

type RegistrarProvider struct {
	Supported map[string]IRegistrar
}

type HostProvider struct {
	Supported map[string]IHost
}

var SupportedProviders = Provider{
	Actions: map[string]func(IProvider, string, chan string) (bool, string){"add": AddProvider, "list": IProvider.List},
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

func AddProvider(provider IProvider, providerName string, channel chan string) (bool, string) {
	fmt.Println("AddProvider call")
	if !cliinit.CliInitialized() {
		channel <- fmt.Sprintln("Performing one time cli configuration...")
		fmt.Println("1 time init")
		cliinit.CliInit()
		fmt.Println("finish 1 time init")
	}

	// TODO: Change
	_, _ = provider.Add(providerName, channel)
	fmt.Println("End AddProvider call")
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
