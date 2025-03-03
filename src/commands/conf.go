package commands

import (
	"fmt"
	"strings"

	"pagecli.com/main/logging"
	"pagecli.com/main/providers"
)

type Conf struct {
}

type ConfArgs struct {
	Action   func(providers.IProvider, string, chan string) error
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

func (c Conf) BindArgs() {
	if !conf.ExecutionOk {
		return
	}

	logMessage := ""

	provider, ok := providers.SupportedProviders.Providers[conf.ArgValues["providerType"]]

	if !ok {
		logMessage = fmt.Sprint("unrecognized value '", conf.ArgValues["providerType"], "'. Expected either registrar or host\n\n")
		conf.ExecutionOk = false
		conf.ExecutionOutput = logMessage
		logging.SendLog(logging.LogRecord{
			Level:   "error",
			Message: logMessage,
		})
		return
	}

	confArgs.Provider = provider
	action, actionExists := providers.SupportedProviders.Actions[conf.ArgValues["actionName"]]

	if !actionExists {
		logMessage = fmt.Sprint("unrecognized value '", conf.ArgValues["actionName"], "'. Expected either add or list\n\n")
		conf.ExecutionOk = false
		conf.ExecutionOutput = logMessage
		logging.SendLog(logging.LogRecord{
			Level:   "error",
			Message: logMessage,
		})
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
			logMessage = fmt.Sprint("unrecognized value '", conf.ArgValues["providerName"], "' for ", conf.ArgValues["providerType"], ". See 'page ", conf.ArgValues["providerType"], " list' for currently supported ", conf.ArgValues["providerType"], "s\n\n")
			conf.ExecutionOk = false
			conf.ExecutionOutput = logMessage
			logging.SendLog(logging.LogRecord{
				Level:   "error",
				Message: logMessage,
			})
			return
		}
	}
}

func (c Conf) UsageInfoShort() string {
	return "configures defaults for domain registrar and host"
}

func (c Conf) UsageInfoExpanded() string {
	supportedProviderTypesConcat := strings.Join(providers.SupportedProviderTypes[:], ", ")
	supportedActionsConcat := strings.Join(providers.SupportedAction[:], ", ")
	supportedRegistrarProvidersConcat := strings.Join(providers.SupportedRegistrars[:], ", ")
	supportedHostProvidersConcat := strings.Join(providers.SupportedHosts[:], ", ")

	extendedUsage := fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Summary")
	extendedUsage += fmt.Sprintln(conf.DisplayName, "-", c.UsageInfoShort())
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Description")
	extendedUsage += fmt.Sprintln("Lets you configure the default host and domain name registrar for all page projects. You can also use the command to view currently supported hosts/registrars and currently configured hosts/registrars.")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Arguments")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "[provider type] [action] [provider name]")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("[provider type] - indicates what provider you are working with (either", supportedProviderTypesConcat, "). Any other value will result in an error.")
	extendedUsage += fmt.Sprintln("[action] - indicates what action you would like to perform on provider type either", supportedActionsConcat, ". Any other value will result in an error.")
	extendedUsage += fmt.Sprintln("[provider name] - required only if action is specified as 'add'. Provider name must be", supportedRegistrarProvidersConcat, " if specifying a registrar or", supportedHostProvidersConcat, " for a host. See supported hosts/registrars with 'page conf [provider type] list'. Any other value will result in an error.")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Example usage")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "host list")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "registrar list")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "host add aws")
	extendedUsage += fmt.Sprintln("page", conf.DisplayName, "registrar add namecheap")
	extendedUsage += fmt.Sprintln()
	return extendedUsage
}

func (c Conf) UsageCategory() int {
	return 1
}

func (c Conf) Execute() {
	if !conf.ExecutionOk {
		return
	}

	err := confArgs.Action(confArgs.Provider, conf.ArgValues["providerName"], OutputChannel)
	if err != nil {
		conf.ExecutionOk = false
		logging.SendLog(logging.LogRecord{
			Level:   "error",
			Message: err.Error(),
		})
	}
}

func (c Conf) Output() string {
	return conf.ExecutionOutput
}
