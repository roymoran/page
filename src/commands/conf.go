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
	Action   func(providers.IProvider, string) (bool, string)
	Provider providers.IProvider
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
	Action:   nil,
	Provider: nil,
}

func (c Conf) LoadArgs() {
	if !conf.ExecutionOk {
		return
	}

	provider, ok := providers.SupportedProviders.Providers[conf.ArgValues["providerType"]]
	if !ok {
		conf.ExecutionOk = false
		conf.ExecutionOutput = fmt.Sprint("unrecognized value '", conf.ArgValues["providerType"], "'. Expected either registrar or host\n")
		return
	}

	confArgs.Provider = provider

	action, actionExists := providers.SupportedProviders.Actions[conf.ArgValues["actionName"]]

	if !actionExists {
		conf.ExecutionOk = false
		conf.ExecutionOutput = fmt.Sprint("unrecognized value '", conf.ArgValues["actionName"], "'. Expected either add or list\n")
		return
	}

	confArgs.Action = action
	var hostProviderConcrete providers.HostProvider
	var registrarProviderConcrete providers.RegistrarProvider
	var providerSupported bool

	if conf.ArgValues["actionName"] != "list" {
		if conf.ArgValues["providerType"] == "registrar" {
			registrarProviderConcrete = provider.(providers.RegistrarProvider)
			_, providerSupported = registrarProviderConcrete.Supported[conf.ArgValues["providerName"]]
		} else {
			hostProviderConcrete = provider.(providers.HostProvider)
			_, providerSupported = hostProviderConcrete.Supported[conf.ArgValues["providerName"]]
		}

		if !providerSupported {
			conf.ExecutionOk = false
			conf.ExecutionOutput = fmt.Sprint("unrecognized value '", conf.ArgValues["providerName"], "' for ", conf.ArgValues["providerType"], ". See 'page ", conf.ArgValues["providerType"], " list' for currently supported ", conf.ArgValues["providerType"], "s\n")
			return
		}
	}
}

func (c Conf) UsageInfoShort() string {
	return "configures defaults for domain registrar and host"
}

func (c Conf) UsageInfoExpanded() string {
	extendedUsage := fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Summary")
	extendedUsage += fmt.Sprintln(conf.DisplayName, "-", c.UsageInfoShort())
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Description")
	extendedUsage += fmt.Sprintln("Lets you configure the default host and domain name registrar for all page projects. You can also use the command to view currently supported hosts and registrars.")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Arguments")
	extendedUsage += fmt.Sprintln("Expects at least", conf.MinimumExpectedArgs, "arguments depending on whether you want to add a host/registrar or list the supported hosts/registrars. See example usage below.")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Example usage")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "host list")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "registrar list")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "host add namecheap")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "registrar add namecheap")
	extendedUsage += fmt.Sprintln()
	return extendedUsage
}

func (c Conf) UsageCategory() int {
	return 2
}

func (c Conf) Execute() {
	if !conf.ExecutionOk {
		return
	}

	_, conf.ExecutionOutput = confArgs.Action(confArgs.Provider, conf.ArgValues["providerName"])
}

func (c Conf) Output() string {
	return conf.ExecutionOutput
}
