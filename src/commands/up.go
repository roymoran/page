package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
	"builtonpage.com/main/providers"
	"builtonpage.com/main/providers/hosts"
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
	extendedUsage += fmt.Sprintln("the command will fail. If the domain must exists on the configured registrar")
	extendedUsage += fmt.Sprintln("For self-hosting on cloud platforms like azure, aws, etc. if the infrastructure")
	extendedUsage += fmt.Sprintln("will be created (this is typically a storage resource exposed for public access).")
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

// TODO: What if the same site gets redeployed twice?
func (u Up) Execute() {
	if !up.ExecutionOk {
		return
	}

	pageDefinition, err := definition.ReadDefinitionFile()
	if err != nil {
		log.Fatalln("error:" + err.Error())
		up.ExecutionOutput += err.Error()
		return
	}

	tempDir, _ := ioutil.TempDir("", "template")
	defer os.RemoveAll(tempDir)

	registrarProvider, _ := providers.SupportedProviders.Providers["registrar"]
	hostProvider, _ := providers.SupportedProviders.Providers["host"]

	var hostProviderConcrete providers.HostProvider
	var registrarProviderConcrete providers.RegistrarProvider
	var providerSupported bool

	var registrar providers.IRegistrar = nil
	var registrarAlias string = pageDefinition.Registrar
	registrarProviderConcrete = registrarProvider.(providers.RegistrarProvider)
	registrarName, err := cliinit.FindRegistrarByAlias(registrarAlias)
	registrar, providerSupported = registrarProviderConcrete.Supported[registrarName]

	// alias may not exist, in which case assume
	// name of registrar povider was provided.
	// Retrieve the default alias for that registrar.
	if err != nil {
		registrar, providerSupported = registrarProviderConcrete.Supported[pageDefinition.Registrar]
		registrarAlias, _ = cliinit.FindDefaultAliasForRegistrar(pageDefinition.Registrar)
		if !providerSupported {
			up.ExecutionOk = false
			OutputChannel <- "Provided unsupported registrar (" + pageDefinition.Registrar + "). See 'page conf registrar list' for supported registrars."
			return
		}
	}

	var host providers.IHost = nil
	var hostAlias string = pageDefinition.Host
	hostProviderConcrete = hostProvider.(providers.HostProvider)
	hostName, err := cliinit.FindHostByAlias(hostAlias)
	host, providerSupported = hostProviderConcrete.Supported[hostName]

	// alias may not exist, in which case assume
	// name of host provider was provided. Retrieve
	// the default alias for that host.
	if err != nil {
		host, providerSupported = hostProviderConcrete.Supported[pageDefinition.Host]
		hostAlias, _ = cliinit.FindDefaultAliasForHost(pageDefinition.Host)

		if !providerSupported {
			up.ExecutionOk = false
			OutputChannel <- "Provided unsupported host or non-existing alias (" + pageDefinition.Host + "). See 'page conf host list' for supported hosts."
			return
		}
	}

	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: pageDefinition.Template,
	})
	// TODO: Remove system files to avoid uploading those (e.g. .DS_Store for macos)
	// remove .git directory to avoid uploading .git content
	os.RemoveAll(filepath.Join(tempDir, ".git"))

	if err != nil {
		OutputChannel <- "Error fetching template at " + pageDefinition.Template + ". (details: " + err.Error() + ")"
		return
	}

	writeCertificateToHostModule(hostAlias, registrarAlias, pageDefinition.Domain)
	writeCnameDomainToRegistrarModule(hostAlias, registrarAlias, pageDefinition.Domain)

	registrar.ConfigureRegistrar(registrarAlias, pageDefinition)

	err = host.ConfigureHost(hostAlias, tempDir, pageDefinition)
	if err != nil {
		up.ExecutionOutput += err.Error()
		return
	}

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

func writeCertificateToHostModule(hostAlias string, registrarAlias string, siteDomain string) {
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)
	var hostModuleTemplate hosts.HostModuleTemplate
	templateData, _ := ioutil.ReadFile(cliinit.ModuleTemplatePath("host", hostAlias))
	_ = json.Unmarshal(templateData, &hostModuleTemplate)
	newCert := hosts.Certificate{
		CertificateChain: "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.certificate_chain}",
		PrivateKeyPem:    "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.private_key_pem}",
		CertificatePem:   "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.certificate_pem}",
	}

	hostModuleTemplate.Module["host_"+hostAlias].Certificates[formattedDomain+"_certificate"] = newCert

	file, _ := json.MarshalIndent(hostModuleTemplate, "", " ")
	_ = ioutil.WriteFile(cliinit.ModuleTemplatePath("host", hostAlias), []byte(file), 0644)
}

func writeCnameDomainToRegistrarModule(hostAlias string, registrarAlias string, siteDomain string) {
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)
	var registrarModuleTemplate hosts.RegistrarModuleTemplate
	templateData, _ := ioutil.ReadFile(cliinit.ModuleTemplatePath("registrar", registrarAlias))
	_ = json.Unmarshal(templateData, &registrarModuleTemplate)
	cnameDomain := hosts.Domain{
		Domain: "${module.host_" + hostAlias + "." + formattedDomain + "_domain}",
	}

	registrarModuleTemplate.Module["registrar_"+registrarAlias].Domains[formattedDomain+"_domain"] = cnameDomain

	file, _ := json.MarshalIndent(registrarModuleTemplate, "", " ")
	_ = ioutil.WriteFile(cliinit.ModuleTemplatePath("registrar", registrarAlias), []byte(file), 0644)
}
