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
	Action          string
	Registrar       providers.IRegistrar
	Host            providers.IHost
	OrderedArgLabel []string
	ArgValues       map[string]string
}

var conf Conf = Conf{
	DisplayName:              "conf",
	ExecutionOutput:          "",
	ExecutionOk:              true,
	MinimumExpectedArgs:      2,
	MaximumExpectedArguments: 3,
}

var confArgs ConfArgs = ConfArgs{
	Action:          "",
	Registrar:       nil,
	Host:            nil,
	OrderedArgLabel: []string{"providerType", "action", "providerName"},
	ArgValues: map[string]string{
		"providerType": "",
		"action":       "",
		"providerName": "",
	},
}

func (c Conf) LoadArgs(args []string) {
	if len(args) < conf.MinimumExpectedArgs {
		conf.ExecutionOk = false
		conf.ExecutionOutput = fmt.Sprintln(conf.DisplayName, "expects at least", conf.MinimumExpectedArgs, "arguments, received", len(args))
		return
	}

	if len(args) > conf.MaximumExpectedArguments {
		conf.ExecutionOk = false
		conf.ExecutionOutput = fmt.Sprintln(conf.DisplayName, "expects at most", conf.MaximumExpectedArguments, "arguments, received", len(args))
		return
	}

	for i, arg := range args {
		confArgs.ArgValues[confArgs.OrderedArgLabel[i]] = arg
	}

	// TODO: Validate args passed for providerType, action, providerName are valid values
	// otherwise set conf.ExecutionOk = false and conf.ExecutionOutput with error message
	conf.ExecutionOutput = fmt.Sprintln(args, confArgs.ArgValues)
}

func (c Conf) UsageInfoShort() string {
	return "configures defaults for domain registrar and host provider"
}

func (c Conf) UsageInfoExpanded() string {
	extendedUsage := fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Summary:")
	extendedUsage += fmt.Sprintln("conf - ", c.UsageInfoShort())
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Description:")
	extendedUsage += fmt.Sprintln("conf lets you configure the default host and domain name registrar for all page projects. You can also use the command to view currently supported hosts and registrars. This command also lets you change/manage the default value of the host or registrar.")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Additonal arguments and options:")
	extendedUsage += fmt.Sprintln("conf does not require any additional arguments or options")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Example usage:")
	extendedUsage += fmt.Sprintln("page conf registrar add namecheap")
	extendedUsage += fmt.Sprintln("page conf host add namecheap")
	extendedUsage += fmt.Sprintln("page conf registrar list")
	extendedUsage += fmt.Sprintln("page conf host list")
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

	conf.ExecutionOutput = fmt.Sprintln(conf.ExecutionOk)
}

func (c Conf) Output() string {
	return conf.ExecutionOutput
}

func (c Conf) ValidArgs() bool {
	return true
}

func (c Conf) AddHost() {

}

func (c Conf) AddRegistrar() {

}

func (c Conf) ListRegistrars() {

}

func (c Conf) ListHosts() {

}
