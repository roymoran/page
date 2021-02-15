package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
)

type AmazonWebServices struct {
	HostName string
}

var accessKey string
var secretKey string

var awsProviderTemplate ProviderTemplate = ProviderTemplate{
	Terraform: RequiredProviders{
		RequiredProvider: map[string]Provider{
			"aws": {
				Source:  "hashicorp/aws",
				Version: "3.25.0",
			},
		},
	},
}

// ConfigureAuth reads user input to request
// the accessKey and secretKey that will be
// stored with this host provider. These
// credentails are used to deploy infrastructure
func (aws AmazonWebServices) ConfigureAuth() error {
	fmt.Print("IAM Access Key: ")
	_, err := fmt.Scanln(&accessKey)
	if err != nil {
		return err
	}
	fmt.Print("IAM Secret Key: ")
	_, err = fmt.Scanln(&secretKey)
	if err != nil {
		return err
	}

	return nil
}

func (aws AmazonWebServices) ConfigureHost(hostAlias string, registrarAlias string, templatePath string, page definition.PageDefinition) error {
	// set up base infra for site to be hosted
	// if not already created
	if baseInfraFile := filepath.Join(cliinit.ProviderAliasPath(aws.HostName, hostAlias), "base.tf.json"); !baseInfraConfigured(baseInfraFile) {
		fmt.Println("creating s3 storage")
		randstr := randSeq(12)
		bucketName := "pagecli" + randstr

		err := ioutil.WriteFile(baseInfraFile, baseInfraTemplate(bucketName), 0644)

		if err != nil {
			fmt.Println("error configureBaseInfra writing base.tf.json for host", aws.HostName)
			return err
		}
	}
	fmt.Println("creating site")
	// TODO: Add case if site is already live and active?
	// maybe show list of sites that are currently live
	// via cli command
	var siteFile string = filepath.Join(cliinit.ProviderAliasPath(aws.HostName, hostAlias), strings.Replace(page.Domain, ".", "_", -1)+".tf.json")
	err := aws.createSite(siteFile, page, templatePath, registrarAlias)
	if err != nil {
		return err
	}
	// TODO: Does the site bucket exist? domain.com_s3bucketobject.tf.json?
	// if no then create it as follows ->
	// TODO: Create s3 bucket object resource to upload site.
	// ensure bucket value is set to output of tf state when
	// bucket was created. Maybe can be done by finding bucket
	// with alias tag?
	// resource "aws_s3_bucket_object" "domain.com" each site
	// under the same host alias should be hosted on the same
	// bucket.
	// TODO: Otherwise if bucket domain.com_s3bucketobject.tf.json exists
	// we can be sure that resource exists on aws.

	// Now create distribution on Cloudfront with output from
	//

	fmt.Println("configured aws host")
	return nil
}

// AddHost creates a new ProviderConfig and writes it
// to the existing config.json
func (aws AmazonWebServices) AddHost(alias string, definitionFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             aws.HostName,
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      "",
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

// ProviderTemplate returns a byte slice that represents
// a template for creating an aws host
func (aws AmazonWebServices) ProviderTemplate() []byte {
	file, _ := json.MarshalIndent(awsProviderTemplate, "", " ")
	return file
}

// ProviderConfigTemplate returns a byte slice that represents
// configuration settings for the aws provider.
func (aws AmazonWebServices) ProviderConfigTemplate() []byte {
	file, _ := json.MarshalIndent(aws.providerConfigTemplate(accessKey, secretKey), "", " ")
	return file
}

// providerConfigTemplate returns a ProviderConfigTemplate struct
// which contains info about the provider configuration including
// authentication fields.
func (aws AmazonWebServices) providerConfigTemplate(accessKey string, secretKey string) ProviderConfigTemplate {
	var awsProviderConfigTemplate ProviderConfigTemplate = ProviderConfigTemplate{
		Provider: map[string]interface{}{
			aws.HostName: map[string]interface{}{
				"region":     "us-east-1",
				"access_key": accessKey,
				"secret_key": secretKey,
			},
		},
	}
	return awsProviderConfigTemplate
}

// baseInfraTemplate returns a byte slice that represents the base
// infrastructure to be deployed on the aws host
func baseInfraTemplate(bucketName string) []byte {
	var awsBaseInfraDefinition BaseInfraTemplate = BaseInfraTemplate{
		Resource: map[string]interface{}{
			"aws_s3_bucket": map[string]interface{}{
				"pages_storage": map[string]interface{}{
					"bucket": bucketName,
					"website": map[string]interface{}{
						"index_document": "index.html",
						"error_document": "index.html",
					},
				},
			},
		},
		Output: map[string]interface{}{
			"bucket": map[string]interface{}{
				"value": "${aws_s3_bucket.pages_storage.bucket}",
			},
			"bucket_regional_domain_name": map[string]interface{}{
				"value": "${aws_s3_bucket.pages_storage.bucket_regional_domain_name}",
			},
		},
	}

	file, _ := json.MarshalIndent(awsBaseInfraDefinition, "", " ")
	return file
}

// siteTemplate returns a byte slice that represents a site
// on the aws host
func siteTemplate(siteDomain string, templatePath string, registrarAlias string) []byte {
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)
	var awsSiteDefinition SiteTemplate = SiteTemplate{
		Site: map[string]interface{}{
			"aws_s3_bucket_object": map[string]interface{}{
				formattedDomain + "_site_files": map[string]interface{}{
					"bucket":       "${aws_s3_bucket.pages_storage.bucket}",
					"key":          siteDomain + "/index.html",
					"source":       filepath.Join(templatePath, "index.html"),
					"acl":          "public-read",
					"content_type": "text/html",
					"depends_on":   []string{"aws_s3_bucket.pages_storage"},
				},
			},
			"aws_cloudfront_distribution": map[string]interface{}{
				formattedDomain + "s3_cdn": map[string]interface{}{
					"origin": map[string]interface{}{
						"domain_name": "${aws_s3_bucket.pages_storage.bucket_regional_domain_name}",
						"origin_path": "/" + siteDomain,
						"origin_id":   formattedDomain,
					},

					"enabled":             true,
					"is_ipv6_enabled":     true,
					"default_root_object": "index.html",
					"aliases":             []string{siteDomain, "*." + siteDomain},

					"default_cache_behavior": map[string]interface{}{
						"allowed_methods":  []string{"GET", "HEAD"},
						"cached_methods":   []string{"GET", "HEAD"},
						"target_origin_id": formattedDomain,

						"forwarded_values": map[string]interface{}{
							"query_string": true,

							"cookies": map[string]interface{}{
								"forward": "none",
							},
						},
						"viewer_protocol_policy": "redirect-to-https",
						"min_ttl":                0,
						"default_ttl":            3600,
						"max_ttl":                86400,
					},
					"restrictions": map[string]interface{}{
						"geo_restriction": map[string]interface{}{
							"restriction_type": "none",
						},
					},
					"viewer_certificate": map[string]interface{}{
						"acm_certificate_arn":      "${aws_acm_certificate." + formattedDomain + "_cert.arn}",
						"ssl_support_method":       "sni-only",
						"minimum_protocol_version": "TLSv1.2_2019",
					},
					"depends_on": []string{"aws_s3_bucket.pages_storage"},
				},
			},

			"tls_private_key": map[string]interface{}{
				formattedDomain + "_tls_private_key": map[string]interface{}{
					"algorithm": "RSA",
				},
			},

			"tls_self_signed_cert": map[string]interface{}{
				formattedDomain + "_tls_self_signed_cert": map[string]interface{}{
					"key_algorithm":   "RSA",
					"private_key_pem": "${acme_certificate." + formattedDomain + "_certificate.private_key_pem}",

					"subject": map[string]interface{}{
						"common_name":  siteDomain,
						"organization": "ACME Examples, Inc",
					},

					"validity_period_hours": 12,

					"allowed_uses": []string{
						"key_encipherment",
						"digital_signature",
						"server_auth",
					},
				},
			},

			"aws_acm_certificate": map[string]interface{}{
				formattedDomain + "_cert": map[string]interface{}{
					"certificate_body":  "${acme_certificate." + formattedDomain + "_certificate.certificate_pem}",
					"private_key":       "${acme_certificate." + formattedDomain + "_certificate.private_key_pem}",
					"certificate_chain": "${acme_certificate." + formattedDomain + "_certificate.certificate_pem}${acme_certificate." + formattedDomain + "_certificate.issuer_pem}",
				},
			},
		},
	}

	file, _ := json.MarshalIndent(awsSiteDefinition, "", " ")
	return file
}

func baseInfraConfigured(baseInfraFile string) bool {
	exists := true
	_, err := os.Stat(baseInfraFile)
	if err != nil {
		return !exists
	}
	return exists
}

func (aws AmazonWebServices) createSite(siteFile string, page definition.PageDefinition, templatePath string, registrarAlias string) error {
	err := ioutil.WriteFile(siteFile, siteTemplate(page.Domain, templatePath, registrarAlias), 0644)

	if err != nil {
		fmt.Println("error createSite writing site template", templatePath)
		return err
	}

	err = TfApply(cliinit.ProvidersPath)
	if err != nil {
		os.Remove(siteFile)
		if strings.Contains(err.Error(), "NoCredentialProviders") {
			return fmt.Errorf("error: missing credentials for %v host", aws.HostName)
		} else if strings.Contains(err.Error(), "InvalidClientTokenId") {
			return fmt.Errorf("error: invalid access_key for %v host", aws.HostName)
		} else if strings.Contains(err.Error(), "SignatureDoesNotMatch") {
			return fmt.Errorf("error: invalid secret_key for %v host", aws.HostName)
		} else {
			// unknown error
			// TODO: Log this
			return err
		}
	}

	return nil
}

// TODO implement method for importing generated ssl cert
func importCertificateTemplate() []byte {
	return []byte{}
}

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
