package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"pagecli.com/main/cliinit"
	"pagecli.com/main/definition"
	"pagecli.com/main/logging"
	"pagecli.com/main/providers"
	"pagecli.com/main/providers/hosts"
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
	extendedUsage += fmt.Sprintln("Publishes your files to a live website using the .yml definition file in")
	extendedUsage += fmt.Sprintln("the current directory. It uses the default host/registrar specified in the")
	extendedUsage += fmt.Sprintln(".yml file. If neither a default host/registrar is found for the provided values,")
	extendedUsage += fmt.Sprintln("the command will fail. The domain must exists on the configured registrar.")
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
	return 0
}

func (u Up) BindArgs() {

}

func (u Up) Execute() {
	if !up.ExecutionOk {
		return
	}

	if !cliinit.CliInitialized() {
		cliinit.CliInit()
	}

	logMessage := ""
	OutputChannel <- "\n"

	pageDefinition, err := definition.ReadDefinitionFile()

	if err != nil {
		logging.SendLog(logging.LogRecord{
			Level:   "critical",
			Message: err.Error(),
		})
		up.ExecutionOutput += err.Error()
		return
	}

	pageDefinitionConfig, err := definition.ProccessDefinitionFile(&pageDefinition)

	if err != nil {
		up.ExecutionOk = false
		logging.SendLog(logging.LogRecord{
			Level:   "critical",
			Message: err.Error(),
		})
		up.ExecutionOutput += err.Error()
		return
	}

	formattedDomain := strings.Replace(pageDefinition.Domain, ".", "_", -1)
	siteDir := filepath.Join(cliinit.SiteFilesPath, formattedDomain)
	os.RemoveAll(siteDir)

	siteDirErr := os.MkdirAll(siteDir, 0755)

	if siteDirErr != nil {
		up.ExecutionOk = false
		up.ExecutionOutput += siteDirErr.Error()
		logging.SendLog(logging.LogRecord{
			Level:   "critical",
			Message: siteDirErr.Error(),
		})
		return
	}

	registrarProvider := providers.SupportedProviders.Providers["registrar"]
	hostProvider := providers.SupportedProviders.Providers["host"]

	var hostProviderConcrete providers.HostProvider
	var registrarProviderConcrete providers.RegistrarProvider

	var registrar providers.IRegistrar = nil
	var registrarAlias string = pageDefinition.Registrar
	registrarProviderConcrete = registrarProvider.(providers.RegistrarProvider)
	registrarName, err := cliinit.FindRegistrarByAlias(registrarAlias)
	registrar = registrarProviderConcrete.Supported[registrarName]

	// alias may not exist, in which case assume
	// name of registrar was provided.
	// Retrieve the default alias for that registrar.
	if err != nil {
		registrar = registrarProviderConcrete.Supported[pageDefinition.Registrar]
		registrarAlias, _ = cliinit.FindDefaultAliasForRegistrar(pageDefinition.Registrar)
	}

	var host providers.IHost = nil
	var hostAlias string = pageDefinition.Host
	hostProviderConcrete = hostProvider.(providers.HostProvider)
	hostName, err := cliinit.FindHostByAlias(hostAlias)
	host = hostProviderConcrete.Supported[hostName]

	// alias may not exist, in which case assume
	// name of host provider was provided. Retrieve
	// the default alias for that host.
	if err != nil {
		host = hostProviderConcrete.Supported[pageDefinition.Host]
		hostAlias, _ = cliinit.FindDefaultAliasForHost(pageDefinition.Host)
	}

	if pageDefinitionConfig.FilesSource == definition.GitURL {
		_, err = git.PlainClone(siteDir, false, &git.CloneOptions{
			URL: pageDefinition.Files,
		})
		os.RemoveAll(filepath.Join(siteDir, ".git"))
	} else {
		// its a file path so just walk the directory and copy the files to the temp dir
		err = filepath.WalkDir(pageDefinition.Files, func(path string, d os.DirEntry, err error) error {

			if err != nil {
				return err
			}

			// if it is a directory, copy the directory structure
			// to the temp dir and return
			if d.IsDir() {
				rel, err := filepath.Rel(pageDefinition.Files, path)
				if err != nil {
					return err
				}

				err = os.MkdirAll(filepath.Join(siteDir, rel), 0755)
				if err != nil {
					return err
				}

				return nil
			}

			rel, err := filepath.Rel(pageDefinition.Files, path)
			if err != nil {
				return err
			}

			// create the directory structure in the temp dir
			// and copy the file to that location
			err = os.MkdirAll(filepath.Join(siteDir, filepath.Dir(rel)), 0755)
			if err != nil {
				return err
			}

			// now read the file contents
			data, err := os.ReadFile(path)

			if err != nil {
				data = make([]byte, 0)
			}

			// now write the file and contents
			err = os.WriteFile(filepath.Join(siteDir, rel), data, 0644)
			if err != nil {
				return err
			}

			return nil
		})
	}

	if err != nil {
		up.ExecutionOk = false
		logMessage = fmt.Sprint("Error fetching files at " + pageDefinition.Files + ". (details: " + err.Error() + ")")
		OutputChannel <- logMessage
		logging.SendLog(logging.LogRecord{
			Level:   "error",
			Message: logMessage,
		})
		return
	}

	if err != nil {
		up.ExecutionOk = false
		up.ExecutionOutput += err.Error()
		logging.SendLog(logging.LogRecord{
			Level:   "critical",
			Message: err.Error(),
		})
		return
	}

	err = host.ConfigureHost(hostAlias, siteDir, pageDefinition)
	if err != nil {
		up.ExecutionOk = false
		up.ExecutionOutput += err.Error()
		logging.SendLog(logging.LogRecord{
			Level:   "critical",
			Message: "Failed to configure host, " + err.Error(),
		})
		return
	}

	err = registrar.ConfigureCertificate(registrarAlias, pageDefinition)
	if err != nil {
		up.ExecutionOk = false
		up.ExecutionOutput += err.Error()
		logging.SendLog(logging.LogRecord{
			Level:   "critical",
			Message: "Failed to configure certificate, " + err.Error(),
		})
		return
	}

	// now adjust host module terraform template to provide
	// certificate details as input to the host module
	err = addCertificatesToHostModule(hostAlias, registrarAlias, formattedDomain)

	err = host.ConfigureWebsite(hostAlias, siteDir, pageDefinition)
	if err != nil {
		up.ExecutionOk = false
		up.ExecutionOutput += err.Error()
		logging.SendLog(logging.LogRecord{
			Level:   "critical",
			Message: "Failed to configure website, " + err.Error(),
		})
		return
	}

	// now adjust host module terraform template to include
	// output values for the host module
	err = addOutputsToHostModule(hostAlias, registrarAlias, formattedDomain)

	// now adjust registrar module terraform template
	err = writeCnameDomainToRegistrarModule(hostAlias, registrarAlias, formattedDomain)

	err = registrar.ConfigureDNS(registrarAlias, pageDefinition)
	if err != nil {
		up.ExecutionOk = false
		up.ExecutionOutput += err.Error()
		logging.SendLog(logging.LogRecord{
			Level:   "critical",
			Message: "Failed to configure registrar, " + err.Error(),
		})
		return
	}

	hostOutputKey := "host_" + hostAlias + "_" + formattedDomain
	certExpiration, _ := hosts.TfSingleOutput[string](hostOutputKey, "certificate_expiry")
	dayBeforeRenew, _ := hosts.TfSingleOutput[int](hostOutputKey, "certificate_min_days_before_renew")
	// domain, _ := hosts.TfSingleOutput[string](hostOutputKey, "domain")

	// Parse the time from the string
	var userReadableExpiration string
	var renewInfo string
	expirationTime, err := time.Parse(time.RFC3339, certExpiration)
	if err != nil {
		userReadableExpiration = certExpiration
		renewInfo = "Certificate Renewal: run 'page up' within " + strconv.Itoa(dayBeforeRenew) + " days of expiration"

	} else {
		// Format the time to be more user-readable
		userReadableExpiration = expirationTime.Format("January 2, 2006 at 3:04pm (MST)")
		// Calculate the renewal window start time
		renewalWindowStart := expirationTime.AddDate(0, 0, -dayBeforeRenew)
		// Calculate the time remaining before the renewal window starts
		timeUntilRenewal := time.Until(renewalWindowStart)
		daysUntilRenewal := int(timeUntilRenewal.Hours() / 24)
		renewInfo = "Certificate Renewal: run 'page up' in " + strconv.Itoa(daysUntilRenewal) + " days"

	}

	up.ExecutionOutput += fmt.Sprintln("")
	up.ExecutionOutput += fmt.Sprintln("Page details")
	up.ExecutionOutput += fmt.Sprintln("Domain: https://" + pageDefinition.Domain)
	up.ExecutionOutput += fmt.Sprintln("Certificate Expires: " + userReadableExpiration)
	up.ExecutionOutput += fmt.Sprintln(renewInfo)
}

func (u Up) Output() string {
	return up.ExecutionOutput
}

func addOutputsToHostModule(hostAlias string, registrarAlias string, formattedSiteDomain string) error {
	var hostModuleTemplate hosts.HostModuleTemplate
	templateData, readFileErr := os.ReadFile(cliinit.ModuleTemplatePath("host", hostAlias))
	if readFileErr != nil {
		return readFileErr
	}

	err := json.Unmarshal(templateData, &hostModuleTemplate)
	if err != nil {
		return err
	}

	hostModuleOutputValueProperties := hosts.HostModuleOutputValueProperties{
		Domain:                        "${module.host_" + hostAlias + "." + formattedSiteDomain + "_domain}",
		CertificateMinDaysBeforeRenew: "${module.registrar_" + registrarAlias + "." + formattedSiteDomain + "_certificate.min_days_remaining}",
		CertificateExpiry:             "${module.registrar_" + registrarAlias + "." + formattedSiteDomain + "_certificate.certificate_not_after}",
	}

	if hostModuleTemplate.Output == nil {
		hostModuleTemplate.Output = map[string]hosts.HostModuleOutputValue{}
	}

	hostModuleTemplate.Output["host_"+hostAlias+"_"+formattedSiteDomain] = hosts.HostModuleOutputValue{
		Value: hostModuleOutputValueProperties,
	}

	file, err := json.MarshalIndent(hostModuleTemplate, "", " ")
	err = os.WriteFile(cliinit.ModuleTemplatePath("host", hostAlias), []byte(file), 0644)
	if err != nil {
		return err
	}

	return nil
}

func addCertificatesToHostModule(hostAlias string, registrarAlias string, formattedSiteDomain string) error {
	var hostModuleTemplate hosts.HostModuleTemplate
	templateData, readFileErr := os.ReadFile(cliinit.ModuleTemplatePath("host", hostAlias))
	if readFileErr != nil {
		return readFileErr
	}

	err := json.Unmarshal(templateData, &hostModuleTemplate)
	if err != nil {
		return err
	}

	newCert := hosts.Certificate{
		CertificateChain: "${module.registrar_" + registrarAlias + "." + formattedSiteDomain + "_certificate.certificate_chain}",
		PrivateKeyPem:    "${module.registrar_" + registrarAlias + "." + formattedSiteDomain + "_certificate.private_key_pem}",
		CertificatePem:   "${module.registrar_" + registrarAlias + "." + formattedSiteDomain + "_certificate.certificate_pem}",
	}

	hostModuleTemplate.Module["host_"+hostAlias].Certificates[formattedSiteDomain+"_certificate"] = newCert

	file, err := json.MarshalIndent(hostModuleTemplate, "", " ")
	err = os.WriteFile(cliinit.ModuleTemplatePath("host", hostAlias), []byte(file), 0644)
	if err != nil {
		return err
	}

	return nil
}

func writeCnameDomainToRegistrarModule(hostAlias string, registrarAlias string, formattedSiteDomain string) error {
	var registrarModuleTemplate hosts.RegistrarModuleTemplate
	templateData, readFileErr := os.ReadFile(cliinit.ModuleTemplatePath("registrar", registrarAlias))

	if readFileErr != nil {
		return readFileErr
	}

	err := json.Unmarshal(templateData, &registrarModuleTemplate)
	if err != nil {
		return err
	}

	cnameDomain := hosts.Domain{
		Domain: "${module.host_" + hostAlias + "." + formattedSiteDomain + "_domain}",
	}

	registrarModuleTemplate.Module["registrar_"+registrarAlias].Domains[formattedSiteDomain+"_domain"] = cnameDomain

	file, err := json.MarshalIndent(registrarModuleTemplate, "", " ")
	err = os.WriteFile(cliinit.ModuleTemplatePath("registrar", registrarAlias), []byte(file), 0644)
	if err != nil {
		return err
	}

	return nil
}
