package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"builtonpage.com/main/cliinit"
)

type AmazonWebServices struct {
	HostName string
}

var hostName string = "aws"

var awsProviderTemplate ProviderTemplate = ProviderTemplate{
	Terraform: RequiredProviders{
		RequiredProvider: map[string]Provider{
			hostName: {
				Source:  "hashicorp/aws",
				Version: "3.25.0",
			},
		},
	},
}

var awsProviderConfigTemplate ProviderConfigTemplate = ProviderConfigTemplate{
	Provider: map[string]ProviderConfig{
		hostName: {
			Profile: "default",
			Region:  "us-east-2",
		},
	},
}

func (aws AmazonWebServices) ConfigureHost(alias string) error {
	if !baseInfraConfigured() {
		err := configureBaseInfra(alias)
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

	//
	fmt.Println("configured aws host")
	return nil
}

// AddHost creates a new ProviderConfig and writes it
// to the existing config.json
func (aws AmazonWebServices) AddHost(alias string, definitionFilePath string, stateFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             hostName,
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      stateFilePath,
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
	file, _ := json.MarshalIndent(awsProviderConfigTemplate, "", " ")
	return file
}

// baseInfraTemplate returns a byte slice that represents the base
// infrastructure to be deployed on the aws host
func baseInfraTemplate() []byte {
	randstr := randSeq(10)
	bucketName := "pagecli-" + randstr
	var awsBaseInfraDefinition BaseInfraTemplate = BaseInfraTemplate{
		Resource: map[string]interface{}{
			"aws_s3_bucket": map[string]interface{}{
				randstr: map[string]interface{}{
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
func siteTemplate() []byte {
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

func baseInfraConfigured() bool {
	return false
}

func configureBaseInfra(alias string) error {
	baseInfraFile := filepath.Join(cliinit.HostPath(hostName), alias+"_base.tf.json")
	err := ioutil.WriteFile(baseInfraFile, baseInfraTemplate(), 0644)

	if err != nil {
		fmt.Println("error configureBaseInfra writing provider.tf.json for host", hostName)
		return err
	}

	err = TfApply(cliinit.HostPath(hostName))
	if err != nil {
		if strings.Contains(err.Error(), "NoCredentialProviders") {
			return fmt.Errorf("Error: missing %v host credentials for %v.\nMake sure you have installed the %v cli and it is configured with your aws credentials", hostName, alias, hostName)
		}
	}

	return nil
}

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
