package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	// lets create the base infrastructure required to host
	// the website. This varies by host.
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

	// If the host is capable of certificate management, lets rely on it.
	// Relying on host for this capability offsets the need for manual
	// intervention for certificate renewals. Otherwise, we rely on ACME/LetsEncrypt
	// for local certificate management and upload the certificate to the host.
	if host.IsManagedCertificateCapable() {
		err = host.ConfigureCertificate(hostAlias, pageDefinition)
		if err != nil {
			up.ExecutionOk = false
			up.ExecutionOutput += err.Error()
			return
		}

		// we've generated the managed certificate, but we need to surface
		// the DNS settings so that the registrar can finalize domain verification
		err = writeDNSRecordVarToRegistrarModule(hostAlias, registrarAlias, formattedDomain, formattedDomain+"_dns_records")
		if err != nil {
			up.ExecutionOk = false
			up.ExecutionOutput += err.Error()
			return
		}
	} else {
		err = registrar.ConfigureCertificate(registrarAlias, pageDefinition)
		if err != nil {
			up.ExecutionOk = false
			up.ExecutionOutput += err.Error()
			return
		}
		// now adjust host module terraform template to provide
		// certificate details as input to the host module
		err = addCertificatesToHostModule(hostAlias, registrarAlias, formattedDomain)
		if err != nil {
			up.ExecutionOk = false
			up.ExecutionOutput += err.Error()
			return
		}
	}

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

	// now adjust registrar module terraform template
	err = writeDNSRecordVarToRegistrarModule(hostAlias, registrarAlias, formattedDomain, formattedDomain+"_domain")
	if err != nil {
		up.ExecutionOk = false
		up.ExecutionOutput += err.Error()
		logging.SendLog(logging.LogRecord{
			Level:   "critical",
			Message: "Failed to configure registrar, " + err.Error(),
		})
		return
	}

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

	up.ExecutionOutput += fmt.Sprintln("")
	up.ExecutionOutput += fmt.Sprintln("Done!")
}

func (u Up) Output() string {
	return up.ExecutionOutput
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
	if err != nil {
		return err
	}

	err = os.WriteFile(cliinit.ModuleTemplatePath("host", hostAlias), []byte(file), 0644)
	if err != nil {
		return err
	}

	return nil
}

func writeDNSRecordVarToRegistrarModule(hostAlias string, registrarAlias string, formattedSiteDomain string, recordVarName string) error {
	var registrarModuleTemplate hosts.RegistrarModuleTemplate
	templateData, readFileErr := os.ReadFile(cliinit.ModuleTemplatePath("registrar", registrarAlias))
	if readFileErr != nil {
		return readFileErr
	}

	err := json.Unmarshal(templateData, &registrarModuleTemplate)
	if err != nil {
		return err
	}

	registrarKey := "registrar_" + registrarAlias
	dnsRecordsKey := formattedSiteDomain + "_dns_records"
	registrarModule := registrarModuleTemplate.Module[registrarKey]
	newDnsRecordVar := "module.host_" + hostAlias + "." + recordVarName
	dnsRecords := hosts.DNSRecords{}

	if registrarModule.DNSRecords == nil {
		registrarModule.DNSRecords = make(map[string]hosts.DNSRecords)
	}

	if _, keyPresent := registrarModule.DNSRecords[dnsRecordsKey]; keyPresent {
		dnsRecords = registrarModule.DNSRecords[dnsRecordsKey]
	}

	for _, recordVar := range dnsRecords.RecordVars {
		if strings.Compare(recordVar, newDnsRecordVar) == 0 {
			return nil
		}
	}

	dnsRecords.RecordVars = append(dnsRecords.RecordVars, newDnsRecordVar)

	registrarModule.DNSRecords[dnsRecordsKey] = dnsRecords
	registrarModuleTemplate.Module[registrarKey] = registrarModule

	file, err := json.MarshalIndent(registrarModuleTemplate, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(cliinit.ModuleTemplatePath("registrar", registrarAlias), []byte(file), 0644)
	if err != nil {
		return err
	}

	err = updateDNSRecordsInRegistrarModule(registrarAlias, formattedSiteDomain)
	if err != nil {
		return err
	}

	return nil
}

func updateDNSRecordsInRegistrarModule(registrarAlias string, formattedSiteDomain string) error {
	var registrarModuleTemplate hosts.RegistrarModuleTemplate
	templateData, readFileErr := os.ReadFile(cliinit.ModuleTemplatePath("registrar", registrarAlias))
	if readFileErr != nil {
		return readFileErr
	}

	err := json.Unmarshal(templateData, &registrarModuleTemplate)
	if err != nil {
		return err
	}

	registrarKey := "registrar_" + registrarAlias
	dnsRecordsKey := formattedSiteDomain + "_dns_records"
	registrarModule := registrarModuleTemplate.Module[registrarKey]

	dnsRecords := registrarModule.DNSRecords[dnsRecordsKey]
	dnsRecords.Records = "${concat("
	for index, recordVar := range dnsRecords.RecordVars {
		if index == len(dnsRecords.RecordVars)-1 {
			dnsRecords.Records += recordVar

		} else {
			dnsRecords.Records += recordVar + ","
		}
	}
	dnsRecords.Records += ")}"

	registrarModule.DNSRecords[dnsRecordsKey] = dnsRecords
	registrarModuleTemplate.Module[registrarKey] = registrarModule

	file, err := json.MarshalIndent(registrarModuleTemplate, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(cliinit.ModuleTemplatePath("registrar", registrarAlias), []byte(file), 0644)
	if err != nil {
		return err
	}

	return nil
}
