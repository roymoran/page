package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
	"builtonpage.com/main/logging"
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
}

// TODO: What if the same site gets redeployed twice?
func (u Up) Execute() {
	if !up.ExecutionOk {
		return
	}

	logMessage := ""

	pageDefinition, err := definition.ReadDefinitionFile()
	if err != nil {
		logging.LogException(err.Error(), true)
		up.ExecutionOutput += err.Error()
		return
	}

	tempDir, tempDirErr := ioutil.TempDir("", "template")
	defer os.RemoveAll(tempDir)

	if tempDirErr != nil {
		up.ExecutionOk = false
		up.ExecutionOutput += err.Error()
		logging.LogException(tempDirErr.Error(), true)
		return
	}

	registrarProvider := providers.SupportedProviders.Providers["registrar"]
	hostProvider := providers.SupportedProviders.Providers["host"]

	var hostProviderConcrete providers.HostProvider
	var registrarProviderConcrete providers.RegistrarProvider
	var providerSupported bool

	var registrar providers.IRegistrar = nil
	var registrarAlias string = pageDefinition.Registrar
	registrarProviderConcrete = registrarProvider.(providers.RegistrarProvider)
	registrarName, err := cliinit.FindRegistrarByAlias(registrarAlias)
	registrar, providerSupported = registrarProviderConcrete.Supported[registrarName]

	// alias may not exist, in which case assume
	// name of registrar was provided.
	// Retrieve the default alias for that registrar.
	if err != nil {
		registrar, providerSupported = registrarProviderConcrete.Supported[pageDefinition.Registrar]
		registrarAlias, _ = cliinit.FindDefaultAliasForRegistrar(pageDefinition.Registrar)
		if !providerSupported {
			up.ExecutionOk = false
			logMessage = fmt.Sprint("Provided unsupported registrar (" + pageDefinition.Registrar + "). See 'page conf registrar list' for supported registrars.")
			OutputChannel <- logMessage
			logging.LogException(logMessage, false)
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
			logMessage = fmt.Sprint("Provided unsupported host or non-existing alias (" + pageDefinition.Host + "). See 'page conf host list' for supported hosts.")
			OutputChannel <- logMessage
			logging.LogException(logMessage, false)
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
		up.ExecutionOk = false
		logMessage = fmt.Sprint("Error fetching template at " + pageDefinition.Template + ". (details: " + err.Error() + ")")
		OutputChannel <- logMessage
		logging.LogException(logMessage, false)
		return
	}

	err = writeCertificateToHostModule(hostAlias, registrarAlias, pageDefinition.Domain)
	err = writeCnameDomainToRegistrarModule(hostAlias, registrarAlias, pageDefinition.Domain)

	if err != nil {
		up.ExecutionOk = false
		up.ExecutionOutput += err.Error()
		logging.LogException(err.Error(), true)
		return
	}

	err = registrar.ConfigureRegistrar(registrarAlias, pageDefinition)
	if err != nil {
		up.ExecutionOk = false
		up.ExecutionOutput += err.Error()
		logging.LogException("Failed to congiure registrar, "+err.Error(), true)
		return
	}

	err = host.ConfigureHost(hostAlias, tempDir, pageDefinition)
	if err != nil {
		up.ExecutionOk = false
		up.ExecutionOutput += err.Error()
		logging.LogException("Failed to congiure host, "+err.Error(), true)
		return
	}
}

func (u Up) Output() string {
	return up.ExecutionOutput
}

func writeCertificateToHostModule(hostAlias string, registrarAlias string, siteDomain string) error {
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)
	var hostModuleTemplate hosts.HostModuleTemplate
	templateData, readFileErr := ioutil.ReadFile(cliinit.ModuleTemplatePath("host", hostAlias))
	if readFileErr != nil {
		return readFileErr
	}

	err := json.Unmarshal(templateData, &hostModuleTemplate)
	if err != nil {
		return err
	}

	newCert := hosts.Certificate{
		CertificateChain: "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.certificate_chain}",
		PrivateKeyPem:    "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.private_key_pem}",
		CertificatePem:   "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.certificate_pem}",
	}

	hostModuleTemplate.Module["host_"+hostAlias].Certificates[formattedDomain+"_certificate"] = newCert

	file, err := json.MarshalIndent(hostModuleTemplate, "", " ")
	err = ioutil.WriteFile(cliinit.ModuleTemplatePath("host", hostAlias), []byte(file), 0644)
	if err != nil {
		return err
	}

	return nil
}

func writeCnameDomainToRegistrarModule(hostAlias string, registrarAlias string, siteDomain string) error {
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)
	var registrarModuleTemplate hosts.RegistrarModuleTemplate
	templateData, readFileErr := ioutil.ReadFile(cliinit.ModuleTemplatePath("registrar", registrarAlias))

	if readFileErr != nil {
		return readFileErr
	}

	err := json.Unmarshal(templateData, &registrarModuleTemplate)
	if err != nil {
		return err
	}

	cnameDomain := hosts.Domain{
		Domain: "${module.host_" + hostAlias + "." + formattedDomain + "_domain}",
	}

	registrarModuleTemplate.Module["registrar_"+registrarAlias].Domains[formattedDomain+"_domain"] = cnameDomain

	file, err := json.MarshalIndent(registrarModuleTemplate, "", " ")
	err = ioutil.WriteFile(cliinit.ModuleTemplatePath("registrar", registrarAlias), []byte(file), 0644)
	if err != nil {
		return err
	}

	return nil
}
