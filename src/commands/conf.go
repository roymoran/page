package commands

import (
	"fmt"

	"builtonpage.com/main/providers"
)

type Conf struct {
	DisplayName              string
	ExecutionOutput          string
	ExecutionOk              bool
	MinimumExpectedArgs      int
	MaximumExpectedArguments int
}

type ConfArgs struct {
	Action          func(providers.IProvider) bool
	Provider        providers.IProvider
	OrderedArgLabel []string
	ArgValues       map[string]string
}

var conf CommandInfo = CommandInfo{
	DisplayName:              "conf",
	ExecutionOutput:          "",
	ExecutionOk:              true,
	MinimumExpectedArgs:      2,
	MaximumExpectedArguments: 3,
	OrderedArgLabel:          []string{"providerType", "actionName", "providerName"},
	ArgValues: map[string]string{
		"providerType": "",
		"actionName":   "",
		"providerName": "",
	},
}

var confArgs ConfArgs = ConfArgs{
	Action:          nil,
	OrderedArgLabel: []string{"providerType", "actionName", "providerName"},
	ArgValues: map[string]string{
		"providerType": "",
		"actionName":   "",
		"providerName": "",
	},
}

func (c Conf) LoadArgs() {
	provider, ok := providers.SupportedProviders.Providers[confArgs.ArgValues["providerType"]]
	if !ok {
		conf.ExecutionOk = false
		conf.ExecutionOutput = fmt.Sprint("unrecognized value '", confArgs.ArgValues["providerType"], "'. Expected either registrar or host")
		return
	}

	confArgs.Provider = provider

	action, actionExists := providers.SupportedProviders.Actions[confArgs.ArgValues["actionName"]]

	if !actionExists {
		conf.ExecutionOk = false
		conf.ExecutionOutput = fmt.Sprint("unrecognized value '", confArgs.ArgValues["actionName"], "'. Expected either add or list")
		return
	}

	confArgs.Action = action
	var hostProviderConcrete providers.HostProvider
	var registrarProviderConcrete providers.RegistrarProvider
	var providerSupported bool

	if confArgs.ArgValues["actionName"] != "list" {
		if confArgs.ArgValues["providerType"] == "registrar" {
			registrarProviderConcrete = provider.(providers.RegistrarProvider)
			_, providerSupported = registrarProviderConcrete.Supported[confArgs.ArgValues["providerName"]]
		} else {
			hostProviderConcrete = provider.(providers.HostProvider)
			_, providerSupported = hostProviderConcrete.Supported[confArgs.ArgValues["providerName"]]
		}

		if !providerSupported {
			conf.ExecutionOk = false
			conf.ExecutionOutput = fmt.Sprint("unrecognized value '", confArgs.ArgValues["providerName"], "' for ", confArgs.ArgValues["providerType"], ". See 'page ", confArgs.ArgValues["providerType"], " list' for currently supported ", confArgs.ArgValues["providerType"], "s")
			return
		}
	}

	conf.ExecutionOutput = fmt.Sprintln(confArgs.ArgValues)
}

func (c Conf) UsageInfoShort() string {
	return "configures defaults for domain registrar and host provider"
}

func (c Conf) UsageInfoExpanded() string {
	extendedUsage := fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Summary:")
	extendedUsage += fmt.Sprintln(conf.DisplayName, "-", c.UsageInfoShort())
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Description:")
	extendedUsage += fmt.Sprintln("lets you configure the default host and domain name registrar for all page projects.")
	extendedUsage += fmt.Sprintln("You can also use the command to view currently supported hosts and registrars.")
	extendedUsage += fmt.Sprintln("This command also lets you change/manage the default value of the host or registrar.")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Arguments:")
	extendedUsage += fmt.Sprintln("does not require any additional arguments or options")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Example usage:")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "registrar add namecheap")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "host add namecheap")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "registrar list")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "host list")
	extendedUsage += fmt.Sprintln()
	return extendedUsage
}

func (c Conf) UsageCategory() int {
	return 2
}

func (c Conf) Execute() {
	if !conf.ExecutionOk {
		conf.ExecutionOutput += fmt.Sprintln("")
		conf.ExecutionOutput += fmt.Sprint("See 'page help ", conf.DisplayName, "' for usage info.\n")
		return
	}

	conf.ExecutionOutput = fmt.Sprintln(conf.ExecutionOutput, conf.ExecutionOk)
}

func (c Conf) Output() string {
	return conf.ExecutionOutput
}

func (c Conf) AddHost() {

}

func (c Conf) AddRegistrar() {

}

func (c Conf) ListRegistrars() {

}

func (c Conf) ListHosts() {

}
