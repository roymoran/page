package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
)

type Namecheap struct {
	RegistrarName string
}

var apiKey string

func (n Namecheap) ConfigureAuth() error {
	fmt.Print("Namecheap API Key: ")
	_, err := fmt.Scanln(&apiKey)
	if err != nil {
		return err
	}
	return nil
}

func (n Namecheap) ConfigureRegistrar(alias string, pageConfig definition.PageDefinition) error {
	fmt.Println("configured namecheap registrar")
	// TODO: "aws" hardcoded for testing purpose, change
	var certificateFilePath string = filepath.Join(cliinit.ProviderAliasPath("aws", alias), strings.Replace(pageConfig.Domain, ".", "_", -1)+"_certificate.tf.json")
	// TODO: Use dns config values collected on initial adding of registrar
	acmeCertificateResourceTemplate := AcmeCertificateResourceTemplate(n.RegistrarName, "rmoran20", "decc1ef18be94299b37e786d95c40dc7", pageConfig.Domain)
	acmeCertificateResourceFile, _ := json.MarshalIndent(acmeCertificateResourceTemplate, "", " ")
	err := ioutil.WriteFile(certificateFilePath, acmeCertificateResourceFile, 0644)

	if err != nil {
		fmt.Println("error writing acme certificate resource template", err)
		return err
	}

	return nil
}

func (n Namecheap) ConfigureDns() bool {
	return true
}

func (n Namecheap) AddRegistrar(alias string) error {
	provider := cliinit.ProviderConfig{
		Type:             "registrar",
		Alias:            alias,
		Name:             "namecheap",
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: "",
		TfStatePath:      "",
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

// AcmeCertificateResourceTemplate returns the resource definition
// for generating an ssl certificate with ACME using the provided
// DNS, Domain, and DNS configuration (required credentials)
func AcmeCertificateResourceTemplate(dnsProvider string, namecheapApiUser string, namecheapApiKey string, siteDomain string) map[string]interface{} {
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
						"config": map[string]interface{}{
							"NAMECHEAP_API_USER": namecheapApiUser,
							"NAMECHEAP_API_KEY":  namecheapApiKey,
						},
					},
				},
			},
		},
	}
}
