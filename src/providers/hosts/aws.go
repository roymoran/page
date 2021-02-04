package providers

import (
	"encoding/json"
	"fmt"

	"builtonpage.com/main/cliinit"
)

type AmazonWebServices struct {
}

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

var awsProviderConfigTemplate ProviderConfigTemplate = ProviderConfigTemplate{
	Provider: map[string]ProviderConfig{
		"aws": {
			Profile: "default",
			Region:  "us-east-2",
		},
	},
}

func (aws AmazonWebServices) ConfigureHost() bool {
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

	//
	fmt.Println("configured aws host")
	return true
}

// AddHost creates a new ProviderConfig and writes it
// to the existing config.json
func (aws AmazonWebServices) AddHost(alias string, definitionFilePath string, stateFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             "aws",
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      stateFilePath,
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

// HostProviderTemplate returns a byte slice that represents
// a template for creating an aws host
func (aws AmazonWebServices) ProviderTemplate() []byte {
	file, _ := json.MarshalIndent(awsProviderTemplate, "", " ")
	return file
}

// HostProviderConfigTemplate returns a byte slice that represents
// configuration settings for the aws provider.
func (aws AmazonWebServices) ProviderConfigTemplate() []byte {
	file, _ := json.MarshalIndent(awsProviderConfigTemplate, "", " ")
	return file
}

// baseInfraTemplate returns a byte slice that represents the base
// infrastructure to be deployed on the aws host
func (aws AmazonWebServices) baseInfraTemplate(terraformResourceName string) []byte {
	bucketName := "pagecli-2827005964"
	var awsBaseInfraDefinition BaseInfraTemplate = BaseInfraTemplate{
		Resource: map[string]interface{}{
			"aws_s3_bucket": map[string]interface{}{
				terraformResourceName: map[string]interface{}{
					"bucket": bucketName,
					"website": map[string]interface{}{
						"index_document": "index.html",
						"error_document": "index.html",
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
func (aws AmazonWebServices) siteTemplate() []byte {
	var awsSiteDefinition SiteTemplate = SiteTemplate{
		Resource: map[string]interface{}{
			"aws_s3_bucket": map[string]interface{}{
				"b": map[string]interface{}{
					"bucket": "pagecli-2827005964",
					"website": map[string]interface{}{
						"index_document": "index.html",
						"error_document": "index.html",
					},
				},
			},
		},
	}

	file, _ := json.MarshalIndent(awsSiteDefinition, "", " ")
	return file
}
