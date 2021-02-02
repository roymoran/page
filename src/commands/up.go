package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"builtonpage.com/main/definition"
	"builtonpage.com/main/providers"
	"github.com/go-git/go-git/v5"
)

type Up struct {
}

var up CommandInfo = CommandInfo{
	DisplayName:              "up",
	ExecutionOutput:          "",
	ExecutionOk:              true,
	MinimumExpectedArgs:      0,
	MaximumExpectedArguments: 0,
}

func (u Up) UsageInfoShort() string {
	return "publishes the page using the page definition file provided"
}

func (u Up) UsageInfoExpanded() string {
	extendedUsage := fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Summary")
	extendedUsage += fmt.Sprintln(up.DisplayName, "-", u.UsageInfoShort())
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Description")
	extendedUsage += fmt.Sprintln("Publishes your template to a live website using the .yml definition file in")
	extendedUsage += fmt.Sprintln("the current directory. It uses the default host/registrar specified in the")
	extendedUsage += fmt.Sprintln(".yml file. If neither a default host/registrar is found for the provided values,")
	extendedUsage += fmt.Sprintln("the command will fail. If the domain specified does not exists on the configured")
	extendedUsage += fmt.Sprintln("registrar then it will be acquired. For self-hosting on cloud platforms like azure,")
	extendedUsage += fmt.Sprintln("aws, etc. if the infrastructure does not exists it will be created (this is")
	extendedUsage += fmt.Sprintln("typically a storage resource exposed for public access).")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Arguments")
	extendedUsage += fmt.Sprintln("Expects", up.MinimumExpectedArgs, "additional arguments.")
	extendedUsage += fmt.Sprintln()
	extendedUsage += fmt.Sprintln("Example usage")
	extendedUsage += fmt.Sprintln("page", up.DisplayName)
	extendedUsage += fmt.Sprintln()
	return extendedUsage
}

func (u Up) UsageCategory() int {
	return 1
}

func (u Up) BindArgs() {
	if !up.ExecutionOk {
		return
	}
}

func (u Up) Execute() {
	if !up.ExecutionOk {
		return
	}

	pageDefinition, message, err := definition.ReadDefinitionFile()
	if err != nil {
		log.Fatalln("error:" + err.Error())
		up.ExecutionOutput += message
		return
	}

	tempDir, _ := ioutil.TempDir("", "template")
	defer os.RemoveAll(tempDir)

	registrarProvider, _ := providers.SupportedProviders.Providers["registrar"]
	hostProvider, _ := providers.SupportedProviders.Providers["host"]

	var hostProviderConcrete providers.HostProvider
	var registrarProviderConcrete providers.RegistrarProvider
	var providerSupported bool

	registrarProviderConcrete = registrarProvider.(providers.RegistrarProvider)
	registrar, providerSupported := registrarProviderConcrete.Supported[pageDefinition.Registrar]

	if !providerSupported {
		up.ExecutionOk = false
		OutputChannel <- "Provided unsupported registrar (" + pageDefinition.Registrar + "). See 'page conf registrar list' for supported registrars."
		return
	}

	hostProviderConcrete = hostProvider.(providers.HostProvider)
	host, providerSupported := hostProviderConcrete.Supported[pageDefinition.Host]

	if !providerSupported {
		up.ExecutionOk = false
		OutputChannel <- "Provided unsupported host (" + pageDefinition.Host + "). See 'page conf host list' for supported registrars."
		return
	}

	// TODO: Consider cloning to in-memory storage
	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: pageDefinition.Template,
	})

	if err != nil {
		OutputChannel <- "Error fetching template at " + pageDefinition.Template + ". (details: " + err.Error() + ")"
		return
	}

	host.ConfigureHost()
	registrar.ConfigureRegistrar()

	// Resolve template url, is it valid?
	// Download template from url, build static assets as needed,
	// then read build files into memory. Take into consideration the size
	// of the built assets - will it be ok to store in memory until deploy?
	// or maybe copy these one by one into a deploy directory (zip if needed)?
	// maintaining a flag that signals deploy step once assets are ready.

	// Get default host for host_value on yaml file. Does infrastructure
	// exist to deploy assets? If not create infrastructure with message
	// 'Creating infrastructure on [host_value]...'
	// Infrastructure could potentially be defined and created with
	// Infrastructure as Code tool e.g. terraform (this logic)
	// may need to be done 'page conf host...' command

	// Get default registrar for registrar_value on yaml file,
	// does domain exist on registrar? if not register with message
	// 'Registering domain.com with [registrar_value]...'
	// configure dns records as needed so that the custom domain
	// points to the host infrastructure

	// Take assets from deploy directory, and execute depoyment via host
	// cli
	up.ExecutionOutput += fmt.Sprintln("deployed")
}

func (u Up) Output() string {
	return up.ExecutionOutput
}
