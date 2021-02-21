package registrars

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
	"builtonpage.com/main/providers/hosts"
)

type Namecheap struct {
	RegistrarName string
}

func (n Namecheap) ConfigureAuth() (cliinit.Credentials, error) {
	registrarCredentials := cliinit.Credentials{}
	var apiUsername string
	var apiKey string

	fmt.Print("Namecheap Username: ")
	_, err := fmt.Scanln(&apiUsername)
	if err != nil {
		return registrarCredentials, err
	}

	fmt.Print("Namecheap API Key: ")
	_, err = fmt.Scanln(&apiKey)
	if err != nil {
		return registrarCredentials, err
	}

	registrarCredentials.Username = apiUsername
	registrarCredentials.Password = apiKey
	return registrarCredentials, nil
}

func (n Namecheap) ConfigureRegistrar(registrarAlias string, pageConfig definition.PageDefinition) error {
	fmt.Println("configured namecheap registrar")
	var certificateFilePath string = filepath.Join(cliinit.ProviderAliasPath(n.RegistrarName, registrarAlias), strings.Replace(pageConfig.Domain, ".", "_", -1)+"_certificate.tf.json")
	var cnameFilePath string = filepath.Join(cliinit.ProviderAliasPath(n.RegistrarName, registrarAlias), strings.Replace(pageConfig.Domain, ".", "_", -1)+"_cname.tf.json")
	credentials, readCredentialsErr := cliinit.FindRegistrarCredentials(registrarAlias)

	if readCredentialsErr != nil {
		return readCredentialsErr
	}

	acmeCertificateResourceTemplate := AcmeCertificateResourceTemplate(n.RegistrarName, pageConfig.Domain, credentials)
	acmeCertificateResourceFile, _ := json.MarshalIndent(acmeCertificateResourceTemplate, "", " ")

	registrarCnameResourceTemplate := cnameResourceTemplate(pageConfig.Domain)
	cnameResourceFile, _ := json.MarshalIndent(registrarCnameResourceTemplate, "", " ")

	err := ioutil.WriteFile(certificateFilePath, acmeCertificateResourceFile, 0644)
	err = ioutil.WriteFile(cnameFilePath, cnameResourceFile, 0644)

	if err != nil {
		fmt.Println("error writing acme certificate resource template", err)
		return err
	}

	return nil
}

func (n Namecheap) AddRegistrar(alias string, credentials cliinit.Credentials) error {
	provider := cliinit.ProviderConfig{
		Type:             "registrar",
		Alias:            alias,
		Name:             "namecheap",
		Credentials:      credentials,
		Default:          true,
		TfDefinitionPath: "",
		TfStatePath:      "",
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

func (n Namecheap) ProviderDefinition() (string, hosts.Provider) {
	return "namecheap", hosts.Provider{Version: "1.7.0", Source: "robgmills/namecheap"}
}

func (n Namecheap) ProviderConfig(username string, password string) map[string]string {
	return map[string]string{
		"username":    username,
		"api_user":    username,
		"token":       password,
		"ip":          "127.0.0.1",
		"use_sandbox": "false",
	}
}

// AcmeCertificateResourceTemplate returns the resource definition
// for generating an ssl certificate with ACME using the provided
// DNS, Domain, and DNS configuration (required credentials)
func AcmeCertificateResourceTemplate(dnsProvider string, siteDomain string, credentials cliinit.Credentials) map[string]interface{} {
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
						"config": map[string]string{
							"NAMECHEAP_API_USER": credentials.Username,
							"NAMECHEAP_API_KEY":  credentials.Password,
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

func cnameResourceTemplate(siteDomain string) map[string]interface{} {
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)

	var cnameResourceTemplate map[string]interface{} = map[string]interface{}{
		"resource": map[string]interface{}{
			"namecheap_record": map[string]interface{}{
				formattedDomain + "_cname": map[string]interface{}{
					"domain":  siteDomain,
					"address": "${lookup(var.domains, " + fmt.Sprintf(`"`) + formattedDomain + "_domain" + fmt.Sprintf(`"`) + ").domain}.",
					"name":    "@",
					"type":    "CNAME",
				},
			},
		},
	}

	return cnameResourceTemplate
}
