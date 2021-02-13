package providers

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
	providers "builtonpage.com/main/providers/hosts"
)

type IHost interface {
	ConfigureAuth() error
	ConfigureHost(alias string, templatePath string, page definition.PageDefinition) error
	AddHost(alias string, definitionPath string) error
	ProviderTemplate() []byte
	ProviderConfigTemplate() []byte
}

func (hp HostProvider) Add(name string, channel chan string) error {
	// TODO Check if alias for host has already been added. if so return with
	// error
	var alias string = AssignAliasName("host")
	hostProvider := SupportedProviders.Providers["host"].(HostProvider)
	host := hostProvider.Supported[name]
	host.ConfigureAuth()
	providerTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "provider.tf.json")
	// This doesn't work with multiple aliases since
	// provider config file is created only once on host dir configuration
	providerConfigTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "providerconfig.tf.json")

	moduleTemplatePath := filepath.Join(cliinit.HostPath(name), alias+".tf.json")
	if !HostDirectoryConfigured(cliinit.ProviderAliasPath(name, alias)) {
		channel <- fmt.Sprint("Configuring ", name, " host...")
		err := InstallTerraformProvider(name, alias, cliinit.HostPath(name), cliinit.ProviderAliasPath(name, alias), host, providerTemplatePath, providerConfigTemplatePath, moduleTemplatePath)
		if err != nil {
			return err
		}
	}

	channel <- fmt.Sprint("Saving ", name, " host configuration...")
	host.AddHost(alias, providerTemplatePath)

	return nil
}

func (hp HostProvider) List(name string, channel chan string) error {
	for _, hostName := range SupportedHosts {
		channel <- fmt.Sprint(hostName)
	}
	return nil
}

func HostDirectoryConfigured(hostPath string) bool {
	exists := true
	_, err := os.Stat(hostPath)
	if err != nil {
		return !exists
	}
	return exists
}

func InstallTerraformProvider(name string, alias string, hostPath string, hostAliasPath string, host IHost, providerTemplatePath string, providerConfigTemplatePath string, moduleTemplatePath string) error {
	hostDirErr := os.MkdirAll(hostAliasPath, os.ModePerm)
	if hostDirErr != nil {
		os.Remove(hostAliasPath)
		log.Fatalln("error creating host config directory for", hostAliasPath, hostDirErr)
		return hostDirErr
	}

	moduleTemplatePathErr := ioutil.WriteFile(moduleTemplatePath, moduleTemplate(alias), 0644)
	providerTemplatePathErr := ioutil.WriteFile(providerTemplatePath, host.ProviderTemplate(), 0644)
	providerConfigTemplatePathErr := ioutil.WriteFile(providerConfigTemplatePath, host.ProviderConfigTemplate(), 0644)

	if moduleTemplatePathErr != nil || providerTemplatePathErr != nil || providerConfigTemplatePathErr != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(hostAliasPath)
		fmt.Println("failed ioutil.WriteFile for provider template")
		return fmt.Errorf("failed ioutil.WriteFile for provider template")
	}

	err := providers.TfInit(hostPath)
	if err != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(hostAliasPath)
		fmt.Println("failed init on new terraform directory", hostPath)
		return err
	}

	return nil
}

func moduleTemplate(alias string) []byte {

	var awsProviderTemplate providers.ModuleTemplate = providers.ModuleTemplate{
		Module: map[string]interface{}{
			alias: map[string]interface{}{
				"source": "./" + alias,
			},
		},
		Output: map[string]interface{}{
			alias + "_bucket": map[string]interface{}{
				"value": "${module.alias1.bucket}",
			},
			alias + "_bucket_regional_domain_name": map[string]interface{}{
				"value": "${module.alias1.bucket_regional_domain_name}",
			},
		},
	}

	file, _ := json.MarshalIndent(awsProviderTemplate, "", " ")
	return file
}

// AcmeCertificateResourceTemplate returns the resource definition
// for generating an ssl certificate with ACME using the provided
// DNS, Domain, and DNS configuration (required credentials)
func AcmeCertificateResourceTemplate(dnsProvider string, dnsProviderConfig map[string]interface{}, siteDomain string) map[string]interface{} {
	wildcardSubdomain := "*." + siteDomain
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)
	return map[string]interface{}{
		"resource": map[string]interface{}{
			"acme_certificate": map[string]interface{}{
				formattedDomain + "_certificate": map[string]interface{}{
					"account_key_pem":           "${acme_registration.reg.account_key_pem}",
					"common_name":               siteDomain,
					"subject_alternative_names": []string{wildcardSubdomain},

					"dns_challenge": map[string]interface{}{
						"provider": dnsProvider,
						"config":   dnsProviderConfig,
					},
				},
			},
		},
	}
}
