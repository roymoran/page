package hosts

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pagecli.com/main/cliinit"
	"pagecli.com/main/definition"
	"pagecli.com/main/progress"
	"pagecli.com/main/terraformutils"
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
				Version: "5.26.0",
			},
			"external": {
				Source:  "hashicorp/external",
				Version: "2.3.2",
			},
		},
	},
}

// ConfigureAuth reads user input to request
// the accessKey and secretKey that will be
// stored with this host provider. These
// credentials are used to deploy infrastructure
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

func (aws AmazonWebServices) ConfigureHost(hostAlias string, templatePath string, page definition.PageDefinition) error {
	if baseInfraFile := filepath.Join(cliinit.ProviderAliasPath(aws.HostName, hostAlias), "base.tf.json"); !terraformutils.ResourcesConfigured(baseInfraFile) {
		randstr := randSeq(12)
		bucketName := "pagecli" + randstr

		err := os.WriteFile(baseInfraFile, baseInfraTemplate(bucketName), 0644)

		if err != nil {
			fmt.Println("error configureBaseInfra writing base.tf.json for host", aws.HostName)
			return err
		}

		err = TfApply(progress.HostCheck, progress.HostProvisioningSequence, progress.StandardTimeout)
		if err != nil {
			os.Remove(baseInfraFile)
		}

		return nil
	}

	// host already configured assumed by the presence of base.tf.json
	// TODO: send the right targets for this host
	var moduleIdentifier string = "module.host_" + hostAlias + "."
	TfApplyWithTarget(progress.HostCheck, progress.ValidatingSequence, progress.StandardTimeout, []string{moduleIdentifier + "aws_s3_bucket.pages_storage", moduleIdentifier + "aws_s3_bucket_policy.pages_storage_policy", moduleIdentifier + "aws_s3_bucket_public_access_block.pages_storage_pab", moduleIdentifier + "aws_s3_bucket_website_configuration.pages_storage_website_configuration", moduleIdentifier + "data.aws_iam_policy_document.pages_storage_policy_document"})
	return nil
}

func (aws AmazonWebServices) ConfigureWebsite(hostAlias string, templatePath string, page definition.PageDefinition) error {
	if siteInfraFile := filepath.Join(cliinit.ProviderAliasPath(aws.HostName, hostAlias), strings.Replace(page.Domain, ".", "_", -1)+"_site.tf.json"); !terraformutils.ResourcesConfigured(siteInfraFile) {
		err := os.WriteFile(siteInfraFile, siteTemplate(page.Domain, templatePath, hostAlias), 0644)

		if err != nil {
			fmt.Println("error ConfigureWebsite writing site infra file", templatePath)
			return err
		}

		err = TfApply(progress.WebsiteFilesCheck, progress.HostWebsiteFilesUploadingSequence, 10*time.Minute)
		if err != nil {
			os.Remove(siteInfraFile)
			if strings.Contains(err.Error(), "NoCredentialProviders") {
				return fmt.Errorf("error: missing credentials for %v host", aws.HostName)
			} else if strings.Contains(err.Error(), "InvalidClientTokenId") {
				return fmt.Errorf("error: invalid access_key for %v host", aws.HostName)
			} else if strings.Contains(err.Error(), "SignatureDoesNotMatch") {
				return fmt.Errorf("error: invalid secret_key for %v host", aws.HostName)
			} else {
				// TODO: Log this
				return err
			}
		}

		return nil
	}

	formattedDomain := strings.Replace(page.Domain, ".", "_", -1)

	var moduleIdentifier string = "module.host_" + hostAlias + "."
	var certificateIdentifier string = moduleIdentifier + "aws_acm_certificate." + formattedDomain + "_cert"
	var cdnIdentifier string = moduleIdentifier + "aws_cloudfront_distribution.pagecli_com_s3_cdn"
	var siteFilesIdentifier string = moduleIdentifier + "aws_s3_object.pagecli_com_site_files"
	var tlsPrivateKeyIdentifier string = moduleIdentifier + "tls_private_key.pagecli_com_tls_private_key"
	var tlsSelfSignedCertIdentifier string = moduleIdentifier + "tls_self_signed_cert.pagecli_com_tls_self_signed_cert"
	TfApplyWithTarget(progress.WebsiteFilesCheck, progress.ValidatingSequence, progress.StandardTimeout, []string{certificateIdentifier, cdnIdentifier, siteFilesIdentifier, tlsPrivateKeyIdentifier, tlsSelfSignedCertIdentifier})

	return nil
}

// AddHost creates a new ProviderConfig and writes it
// to the existing config.json
func (aws AmazonWebServices) AddHost(alias string, definitionFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             aws.HostName,
		Credentials:      cliinit.Credentials{},
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
				},
			},
			"aws_s3_bucket_policy": map[string]interface{}{
				"pages_storage_policy": map[string]interface{}{
					"bucket": bucketName,
					"policy": "${data.aws_iam_policy_document.pages_storage_policy_document.json}",
				},
			},
			// TODO: Consider that a template may not
			// have index.html entry point.
			"aws_s3_bucket_website_configuration": map[string]interface{}{
				"pages_storage_website_configuration": map[string]interface{}{
					"bucket": bucketName,
					"index_document": map[string]interface{}{
						"suffix": "index.html",
					},
					"error_document": map[string]interface{}{
						"key": "index.html",
					},
				},
			},
			"aws_s3_bucket_public_access_block": map[string]interface{}{
				"pages_storage_pab": map[string]interface{}{
					"bucket": bucketName,

					"block_public_acls":       false,
					"block_public_policy":     false,
					"ignore_public_acls":      false,
					"restrict_public_buckets": false,
				},
			},
		},
		Data: map[string]interface{}{
			"aws_iam_policy_document": map[string]interface{}{
				"pages_storage_policy_document": map[string]interface{}{
					"statement": map[string]interface{}{
						"principals": map[string]interface{}{
							"type":        "*",
							"identifiers": []string{"*"},
						},
						"sid":       "PublicReadGetObject",
						"effect":    "Allow",
						"actions":   []string{"s3:GetObject"},
						"resources": []string{"${aws_s3_bucket.pages_storage.arn}/*"},
					},
				},
			},
		},
	}

	file, _ := json.MarshalIndent(awsBaseInfraDefinition, "", " ")
	return file
}

// siteTemplate returns a byte slice that represents a site
// on the aws host
func siteTemplate(siteDomain string, templatePath string, hostAlias string) []byte {
	formattedDomain := strings.Replace(siteDomain, ".", "_", -1)
	executablePath, _ := os.Executable()
	var awsSiteDefinition map[string]interface{} = map[string]interface{}{
		"resource": map[string]interface{}{
			"aws_s3_object": map[string]interface{}{
				formattedDomain + "_site_files": map[string]interface{}{
					"for_each": "${fileset(" + fmt.Sprintf(`"`) + templatePath + fmt.Sprintf(`"`) + "," + fmt.Sprintf(`"`) + "**/*" + fmt.Sprintf(`"`) + ")}",

					"bucket":       "${aws_s3_bucket.pages_storage.bucket}",
					"key":          siteDomain + "/${each.value}",
					"source":       filepath.Join(templatePath, "${each.value}"),
					"content_type": "${data.external.assign_content_type[each.value].result[" + fmt.Sprintf(`"`) + "mimetype" + fmt.Sprintf(`"`) + "]}",
					"etag":         "${filemd5(" + fmt.Sprintf(`"`) + filepath.Join(templatePath, "${each.value}") + fmt.Sprintf(`"`) + ")}",
					"depends_on":   []string{"aws_s3_bucket.pages_storage"},
				},
			},
			"aws_cloudfront_distribution": map[string]interface{}{
				formattedDomain + "_s3_cdn": map[string]interface{}{
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
					"private_key_pem": "${lookup(var.certificates, " + fmt.Sprintf(`"`) + formattedDomain + "_certificate" + fmt.Sprintf(`"`) + ").private_key_pem}",

					"subject": map[string]interface{}{
						"common_name":  siteDomain,
						"organization": siteDomain + ", Inc",
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
					"certificate_body":  "${lookup(var.certificates," + fmt.Sprintf(`"`) + formattedDomain + "_certificate" + fmt.Sprintf(`"`) + ").certificate_pem}",
					"private_key":       "${lookup(var.certificates," + fmt.Sprintf(`"`) + formattedDomain + "_certificate" + fmt.Sprintf(`"`) + ").private_key_pem}",
					"certificate_chain": "${lookup(var.certificates," + fmt.Sprintf(`"`) + formattedDomain + "_certificate" + fmt.Sprintf(`"`) + ").certificate_chain}",
				},
			},
		},
		"data": map[string]interface{}{
			// TODO: Can be eliminated if terraform contains
			// a built-in function for determining content type
			"external": map[string]interface{}{
				"assign_content_type": map[string]interface{}{
					"for_each": "${fileset(" + fmt.Sprintf(`"`) + templatePath + fmt.Sprintf(`"`) + "," + fmt.Sprintf(`"`) + "**/*" + fmt.Sprintf(`"`) + ")}",

					"program": []string{executablePath, "infra", "mimetype", "${each.value}"},
				},
			},
		},
		"output": map[string]interface{}{
			formattedDomain + "_domain": map[string]interface{}{
				"value": "${aws_cloudfront_distribution." + formattedDomain + "_s3_cdn.domain_name}",
			},
		},
	}

	file, _ := json.MarshalIndent(awsSiteDefinition, "", " ")
	return file
}

func randSeq(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
