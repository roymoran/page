package providers

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/logging"
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
	Actions: map[string]func(IProvider, string, chan string) error{"add": AddProvider, "list": ListProvider},
	Providers: map[string]IProvider{
		"host": HostProvider{
			Supported: map[string]IHost{
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
	if err != nil {
		logging.LogException("Failed to add provider. Details: "+err.Error(), true)
	}

	return err
}

// TODO: Maybe we don't want the channel to be propegated?
// if there is an error in Adding the provider or initializing
// the cli we can let conf.go receive the message via err.Error()
// have conf.go/handler.go communicate that over the channel
func ListProvider(provider IProvider, providerName string, channel chan string) error {
	err := provider.List(providerName, channel)
	if err != nil {
		logging.LogException("Failed to list provider. Details: "+err.Error(), true)
	}

	return err
}

// SupportedProviderTypes returns a string slice containing the different provider
// types e.g. registrar, host
var SupportedProviderTypes []string = BuildSupportedProviderTypes()

// SupportedAction returns a string slice containing the different actions
// that can be performed on a provider e.g. add, list
var SupportedAction []string = BuildSupportedActions()

// SupportedRegistrars returns a string slice containing the
// currently supported registrars
var SupportedRegistrars []string = BuildSupportedRegistrars()

// SupportedHosts returns a string slice containing the
// currently supported hosts
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
		if AliasForProviderExists(providerType, alias) {
			fmt.Print("The alias should be unique across all aliases for your " + providerType + "s\n\n")
			continue
		}
		for _, providerName := range supportedProviders {
			if providerName == alias {
				valid = !valid
				// TODO: WHAT IF AN ALIAS IS ADDED THAT BECOMES AN IVALID NAME IN THE FUTURE?
				// for example cli currently does not support firebase host so user can use 'firebase'
				// as their alias. Once firebase is supported this alias may become unsupported.

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

// AliasForProviderExists checks if the alias already
// exists for the provider type,
func AliasForProviderExists(providerType string, alias string) bool {
	var providerExistsErr error

	if providerType == "registrar" {
		_, providerExistsErr = cliinit.FindRegistrarByAlias(alias)
	} else if providerType == "host" {
		_, providerExistsErr = cliinit.FindHostByAlias(alias)
	}

	// alias already exists for provider
	if providerExistsErr == nil {
		return true
	}

	return false
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
