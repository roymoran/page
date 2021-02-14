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

var acmeResourceTemplate map[string]interface{} = map[string]interface{}{
	"resource": map[string]interface{}{
		"tls_private_key": map[string]interface{}{
			"private_key": map[string]interface{}{
				"algorithm": "RSA",
			},
		},

		"acme_registration": map[string]interface{}{
			"reg": map[string]interface{}{
				"account_key_pem": "${tls_private_key.private_key.private_key_pem}",
				// TODO: Read from user
				"email_address": "romoran1@outlook.com",
			},
		},
	},
}

type IRegistrar interface {
	ConfigureAuth() error
	ConfigureRegistrar(definition.PageDefinition) bool
	ConfigureDns() bool
	AddRegistrar(string) error
}

func (rp RegistrarProvider) Add(name string, channel chan string) error {
	var alias string = AssignAliasName("registrar")

	registrarProvider := SupportedProviders.Providers["registrar"].(RegistrarProvider)
	registrar := registrarProvider.Supported[name]
	registrar.ConfigureAuth()

	providerTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "provider.tf.json")
	// This doesn't work with multiple aliases since
	// provider config file is created only once on host dir configuration
	providerConfigTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "providerconfig.tf.json")

	moduleTemplatePath := filepath.Join(cliinit.ProvidersPath, name+"_"+alias+".tf.json")

	if !AliasDirectoryConfigured(cliinit.ProviderAliasPath(name, alias)) {
		channel <- fmt.Sprint("Configuring ", name, " registrar...")
		err := InstallAcmeTerraformProvider(name, alias, cliinit.ProviderAliasPath(name, alias), providerTemplatePath, providerConfigTemplatePath, moduleTemplatePath)
		if err != nil {
			return err
		}
	}

	// acmeResourceTemplatePath := filepath.Join(cliinit.ProviderAliasPath(name, alias), "acme.tf.json")
	registrar.AddRegistrar(alias)

	return nil
}

func (rp RegistrarProvider) List(name string, channel chan string) error {
	for _, registrarName := range SupportedRegistrars {
		channel <- fmt.Sprint(registrarName)
	}
	return nil
}

func InstallAcmeTerraformProvider(name string, alias string, providerAliasPath string, providerTemplatePath string, providerConfigTemplatePath string, moduleTemplatePath string) error {
	hostDirErr := os.MkdirAll(providerAliasPath, os.ModePerm)
	if hostDirErr != nil {
		os.Remove(providerAliasPath)
		log.Fatalln("error creating host config directory for", providerAliasPath, hostDirErr)
		return hostDirErr
	}

	moduleTemplatePathErr := ioutil.WriteFile(moduleTemplatePath, registrarModuleTemplate(name, alias), 0644)
	providerTemplatePathErr := ioutil.WriteFile(providerTemplatePath, acmeProviderTemplate(), 0644)
	providerConfigTemplatePathErr := ioutil.WriteFile(providerConfigTemplatePath, acmeProviderConfigTemplate(), 0644)

	if moduleTemplatePathErr != nil || providerTemplatePathErr != nil || providerConfigTemplatePathErr != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(providerAliasPath)
		fmt.Println("failed ioutil.WriteFile for provider template")
		return fmt.Errorf("failed ioutil.WriteFile for provider template")
	}

	err := providers.TfInit(cliinit.ProvidersPath)
	if err != nil {
		os.Remove(moduleTemplatePath)
		os.RemoveAll(providerAliasPath)
		fmt.Println("failed init on new terraform directory", cliinit.ProvidersPath)
		return err
	}

	return nil
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

func registrarModuleTemplate(providerName string, alias string) []byte {

	var awsProviderTemplate providers.ModuleTemplate = providers.ModuleTemplate{
		Module: map[string]interface{}{
			"registrar_" + alias: map[string]interface{}{
				"source": "./" + providerName + "/" + alias,
			},
		},
	}

	file, _ := json.MarshalIndent(awsProviderTemplate, "", " ")
	return file
}

func acmeProviderTemplate() []byte {
	var providerTemplate providers.ProviderTemplate = providers.ProviderTemplate{
		Terraform: providers.RequiredProviders{
			RequiredProvider: map[string]providers.Provider{
				"acme": {
					Source:  "vancluever/acme",
					Version: "2.0.0",
				},
			},
		},
	}

	file, _ := json.MarshalIndent(providerTemplate, "", " ")
	return file
}

func acmeProviderConfigTemplate() []byte {
	var providerConfigTemplate providers.ProviderConfigTemplate = providers.ProviderConfigTemplate{
		Provider: map[string]interface{}{
			"acme": map[string]interface{}{
				"server_url": "https://acme-staging-v02.api.letsencrypt.org/directory",
			},
		},
	}

	file, _ := json.MarshalIndent(providerConfigTemplate, "", " ")
	return file
}
