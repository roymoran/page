package registrars

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pagecli.com/main/cliinit"
	"pagecli.com/main/definition"
	"pagecli.com/main/progress"
	"pagecli.com/main/providers/hosts"
	"pagecli.com/main/terraformutils"
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

func (n Namecheap) ConfigureDNS(registrarAlias string, pageConfig definition.PageDefinition) error {

	if cnameInfraFile := filepath.Join(cliinit.ProviderAliasPath(n.RegistrarName, registrarAlias), strings.Replace(pageConfig.Domain, ".", "_", -1)+"_cname.tf.json"); !terraformutils.ResourcesConfigured(cnameInfraFile) {
		registrarCnameResourceTemplate := cnameResourceTemplate(pageConfig.Domain)
		cnameResourceFile, _ := json.MarshalIndent(registrarCnameResourceTemplate, "", " ")

		err := os.WriteFile(cnameInfraFile, cnameResourceFile, 0644)

		if err != nil {
			fmt.Println("error writing acme certificate resource template", err)
			return err
		}

		err = hosts.TfApply(progress.DomainCheck, progress.DomainUpdatingSequence, progress.StandardTimeout)

		if err != nil {
			os.Remove(cnameInfraFile)
			return err
		}

		return nil
	}
	var moduleIdentifier string = "module.registrar_" + registrarAlias + "."
	var cnameIdentifier string = moduleIdentifier + "namecheap_domain_records." + strings.Replace(pageConfig.Domain, ".", "_", -1) + "_cname"
	hosts.TfApplyWithTarget(progress.DomainCheck, progress.ValidatingSequence, progress.StandardTimeout, []string{cnameIdentifier})

	return nil
}

func (n Namecheap) ConfigureCertificate(registrarAlias string, pageConfig definition.PageDefinition) error {
	credentials, readCredentialsErr := cliinit.FindRegistrarCredentials(registrarAlias)

	if readCredentialsErr != nil {
		return readCredentialsErr
	}

	acmeCertificateResourceTemplate := AcmeCertificateResourceTemplate(n.RegistrarName, pageConfig.Domain, credentials)
	acmeCertificateResourceFile, _ := json.MarshalIndent(acmeCertificateResourceTemplate, "", " ")

	if certificateInfraFile := filepath.Join(cliinit.ProviderAliasPath(n.RegistrarName, registrarAlias), strings.Replace(pageConfig.Domain, ".", "_", -1)+"_certificate.tf.json"); !terraformutils.ResourcesConfigured(certificateInfraFile) {
		err := os.WriteFile(certificateInfraFile, acmeCertificateResourceFile, 0644)

		if err != nil {
			fmt.Println("error writing acme certificate resource template", err)
			return err
		}

		err = hosts.TfApply(progress.CertificateCheck, progress.CertificateGeneratingSequence, 2*time.Minute)

		if err != nil {
			os.Remove(certificateInfraFile)
			return err
		}

		return nil
	}

	// TODO:
	// at this point we should assume that the certificate has been configured so we should check
	// for renewal and renew if necessary. If certificate is not up for renewal, lets provide
	// some output that indicates to the user when the certificate will be up for renewal
	var moduleIdentifier string = "module.registrar_" + registrarAlias + "."
	hosts.TfApplyWithTarget(progress.CertificateCheck, progress.ValidatingSequence, progress.StandardTimeout, []string{moduleIdentifier + "acme_certificate." + strings.Replace(pageConfig.Domain, ".", "_", -1) + "_certificate"})

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
	return "namecheap", hosts.Provider{Version: "2.1.0", Source: "namecheap/namecheap"}
}

func (n Namecheap) ProviderConfig(username string, password string) map[string]string {
	return map[string]string{
		"user_name":   username,
		"api_user":    username,
		"api_key":     password,
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
					"certificate_pem":       "${acme_certificate." + formattedDomain + "_certificate.certificate_pem}",
					"private_key_pem":       "${acme_certificate." + formattedDomain + "_certificate.private_key_pem}",
					"certificate_chain":     "${acme_certificate." + formattedDomain + "_certificate.certificate_pem}${acme_certificate." + formattedDomain + "_certificate.issuer_pem}",
					"certificate_not_after": "${acme_certificate." + formattedDomain + "_certificate.certificate_not_after}",
					"min_days_remaining":    "${acme_certificate." + formattedDomain + "_certificate.min_days_remaining}",
				},
			},
		},
	}
}

// TODO: Add ability to add multiple records to a domain here
func cnameResourceTemplate(siteDomain string) map[string]interface{} {
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)

	var cnameResourceTemplate map[string]interface{} = map[string]interface{}{
		"resource": map[string]interface{}{
			"namecheap_domain_records": map[string]interface{}{
				formattedDomain + "_cname": map[string]interface{}{
					"domain": siteDomain,
					"record": map[string]interface{}{
						"hostname": "@",
						"type":     "CNAME",
						"address":  "${lookup(var.domains, " + fmt.Sprintf(`"`) + formattedDomain + "_domain" + fmt.Sprintf(`"`) + ").domain}.",
					},
					// TODO: Add ability to add multiple records to a domain here
				},
			},
		},
	}

	return cnameResourceTemplate
}
