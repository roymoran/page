package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
	providers "builtonpage.com/main/providers/hosts"
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

func (n Namecheap) ConfigureRegistrar(registrarAlias string, hostAlias string, pageConfig definition.PageDefinition) error {
	fmt.Println("configured namecheap registrar")
	var certificateFilePath string = filepath.Join(cliinit.ProviderAliasPath(n.RegistrarName, registrarAlias), strings.Replace(pageConfig.Domain, ".", "_", -1)+"_certificate.tf.json")
	// TODO: AcmeCertificateResourceTemplate should exist where it
	// is accessible by any registrar implementation
	acmeCertificateResourceTemplate := AcmeCertificateResourceTemplate(n.RegistrarName, pageConfig.Domain)
	acmeCertificateResourceFile, _ := json.MarshalIndent(acmeCertificateResourceTemplate, "", " ")

	err := ioutil.WriteFile(certificateFilePath, acmeCertificateResourceFile, 0644)

	writeCertificateToHostModule(hostAlias, registrarAlias, pageConfig.Domain)

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
func AcmeCertificateResourceTemplate(dnsProvider string, siteDomain string) map[string]interface{} {
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
							"NAMECHEAP_API_USER": "rmoran20",
							"NAMECHEAP_API_KEY":  "decc1ef18be94299b37e786d95c40dc7",
						},
					},
				},
			},
		},
		"output": map[string]interface{}{
			formattedDomain + "_certificate": map[string]interface{}{
				"value": map[string]interface{}{
					"certificate_pem":   "${acme_certificate." + formattedDomain + "_certificate.certificate_pem}",
					"private_key_pem":   "${acme_certificate." + formattedDomain + "_certificate.private_key_pem}",
					"certificate_chain": "${acme_certificate." + formattedDomain + "_certificate.certificate_pem}${acme_certificate." + formattedDomain + "_certificate.issuer_pem}",
				},
			},
		},
	}
}

func jsonCertificateProperties(siteDomain string, registrarAlias string) map[string]interface{} {
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)
	return map[string]interface{}{
		formattedDomain + "_certificate": map[string]interface{}{
			"certificate_pem":   "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.certificate_pem}",
			"private_key_pem":   "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.private_key_pem}",
			"certificate_chain": "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.certificate_chain}",
		},
	}
}

func writeCertificateToHostModule(hostAlias string, registrarAlias string, siteDomain string) {
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)
	var hostModuleTemplate providers.HostModuleTemplate
	templateData, _ := ioutil.ReadFile(cliinit.ModuleTemplatePath("host", hostAlias))
	_ = json.Unmarshal(templateData, &hostModuleTemplate)
	newCert := providers.Certificate{
		CertificateChain: "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.certificate_chain}",
		PrivateKeyPem:    "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.private_key_pem}",
		CertificatePem:   "${module.registrar_" + registrarAlias + "." + formattedDomain + "_certificate.certificate_pem}",
	}

	hostModuleTemplate.Module["host_"+hostAlias].Certificates[formattedDomain+"_certificate"] = newCert

	file, _ := json.MarshalIndent(hostModuleTemplate, "", " ")
	_ = ioutil.WriteFile(cliinit.ModuleTemplatePath("host", hostAlias), []byte(file), 0644)
}
