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

	if dnsRecordInfraFile := filepath.Join(cliinit.ProviderAliasPath(n.RegistrarName, registrarAlias), strings.Replace(pageConfig.Domain, ".", "_", -1)+"_dns.tf.json"); !terraformutils.ResourcesConfigured(dnsRecordInfraFile) {
		registrarDnsRecordResourceTemplate := dnsResourceTemplate(pageConfig.Domain)
		dnsRecordResourceFile, _ := json.MarshalIndent(registrarDnsRecordResourceTemplate, "", " ")

		err := os.WriteFile(dnsRecordInfraFile, dnsRecordResourceFile, 0644)

		if err != nil {
			fmt.Println("error writing acme certificate resource template", err)
			return err
		}

		err = hosts.TfApply(progress.DomainCheck, progress.DomainUpdatingSequence, progress.StandardTimeout)

		if err != nil {
			os.Remove(dnsRecordInfraFile)
			return err
		}

		return nil
	}
	var moduleIdentifier string = "module.registrar_" + registrarAlias + "."
	var dnsRecordIdentifier string = moduleIdentifier + "namecheap_domain_records." + strings.Replace(pageConfig.Domain, ".", "_", -1) + "_dns_records"
	err := hosts.TfApplyWithTarget(progress.DomainCheck, progress.ValidatingSequence, progress.StandardTimeout, []string{dnsRecordIdentifier})
	if err != nil {
		return err
	}
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
	err := hosts.TfApplyWithTarget(progress.CertificateCheck, progress.ValidatingSequence, progress.StandardTimeout, []string{moduleIdentifier + "acme_certificate." + strings.Replace(pageConfig.Domain, ".", "_", -1) + "_certificate"})
	if err != nil {
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

func dnsResourceTemplate(siteDomain string) map[string]interface{} {
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)

	var dnsRecordResourceTemplate map[string]interface{} = map[string]interface{}{
		"resource": map[string]interface{}{
			"namecheap_domain_records": map[string]interface{}{
				formattedDomain + "_dns_records": map[string]interface{}{
					"domain": siteDomain,
					"dynamic": map[string]interface{}{
						"record": map[string]interface{}{
							"for_each": "${lookup(var.dns_records, " + fmt.Sprintf(`"`) + formattedDomain + "_dns_records" + fmt.Sprintf(`"`) + ").records}",
							"content": map[string]interface{}{
								"hostname": "${replace(record.value.host, \"" + siteDomain + ".\", \"\")}",
								"type":     "${record.value.type}",
								"address":  "${record.value.value}",
							},
						},
					},
				},
			},
		},
	}

	return dnsRecordResourceTemplate
}
