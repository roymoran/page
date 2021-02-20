package providers

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/providers/hosts"
	"builtonpage.com/main/providers/registrars"
)

type Provider struct {
	Actions   map[string]func(IProvider, string, chan string) error
	Providers map[string]IProvider
}

type IProvider interface {
	Add(string, chan string) error
	List(string, chan string) error
}

type RegistrarProvider struct {
	Supported map[string]IRegistrar
}

type HostProvider struct {
	Supported map[string]IHost
}

var SupportedProviders = Provider{
	Actions: map[string]func(IProvider, string, chan string) error{"add": AddProvider, "list": IProvider.List},
	Providers: map[string]IProvider{
		"host": HostProvider{
			Supported: map[string]IHost{
				"page": hosts.PageHost{
					HostName: "page",
				},
				"azure": hosts.Azure{
					HostName: "azure",
				},
				"aws": hosts.AmazonWebServices{
					HostName: "aws",
				},
			},
		},
		"registrar": RegistrarProvider{
			Supported: map[string]IRegistrar{
				"namecheap": registrars.Namecheap{
					RegistrarName: "namecheap",
				},
				"page": registrars.Page{},
			},
		},
	},
}

// TODO: Maybe we don't want the channel to be propegated?
// if there is an error in Adding the provider or initializing
// the cli we can let conf.go receive the message via err.Error()
// have conf.go/handler.go communicate that over the channel
func AddProvider(provider IProvider, providerName string, channel chan string) error {
	if !cliinit.CliInitialized() {
		channel <- fmt.Sprintln("Performing one time cli configuration...")
		cliinit.CliInit()
	}

	err := provider.Add(providerName, channel)
	return err
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

// AssignAliasName allows user to enter an alternate
// name for their provider either host or registrar
func AssignAliasName(providerType string) string {
	var alias string
	var supportedProviders []string = []string{}

	if providerType == "registrar" {
		supportedProviders = BuildSupportedRegistrars()
	} else if providerType == "host" {
		supportedProviders = BuildSupportedHosts()
	}

	for {
		valid := true
		fmt.Print("Give your " + providerType + " an alias: ")
		fmt.Scanln(&alias)
		for _, providerName := range supportedProviders {
			if providerName == alias {
				valid = !valid
				// TODO: WHAT IF AN ALIAS IS ADDED THAT BECOMES AN IVALID NAME IN THE FUTURE?
				// for example cli currently does not support firebase host so user can use 'firebase'
				// as their alias. Once firebase is supported this alias may become unsupported.

				// TODO: We'll need to provide a mechanism to rename an alias
				fmt.Println("alias should not be the same as the name of a " + providerType + " provider (" + strings.Join(supportedProviders[:], ", ") + ")")

				break
			}
		}

		if valid {
			break
		}
	}
	return alias

}

// ProviderDirectoryConfigured is utility function that
// check if a directory already exists given a path
func AliasDirectoryConfigured(aliasPath string) bool {
	exists := true
	_, err := os.Stat(aliasPath)
	if err != nil {
		return !exists
	}
	return exists
}
