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
)

// TODO: Add logic so that there are character restrictions
// to resource names for aws resources like S3.
// ensure these are enforced accord
// S3 Bucket Name Rules: https://docs.aws.amazon.com/AmazonS3/latest/dev/BucketRestrictions.html#bucketnamingrules
type AmazonWebServices struct {
	HostName string
}

var hostName string = "aws"
var accessKey string
var secretKey string

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

// ConfigureAuth reads user input to request
// the accessKey and secretKey that will be
// stored with this host provider. These
// credentails are used to deploy infrastructure
func (aws AmazonWebServices) ConfigureAuth() error {
	fmt.Print("Enter your IAM Access Key: ")
	_, err := fmt.Scanln(&accessKey)
	if err != nil {
		return err
	}
	fmt.Print("Enter your IAM Secret Key: ")
	_, err = fmt.Scanln(&secretKey)
	if err != nil {
		return err
	}

	return nil
}

func (aws AmazonWebServices) ConfigureHost(alias string) error {
	var baseInfraFile string = filepath.Join(cliinit.HostPath(hostName), alias+"_base.tf.json")

	if !baseInfraConfigured(baseInfraFile) {
		err := configureBaseInfra(baseInfraFile)
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
	file, _ := json.MarshalIndent(providerConfigTemplate(accessKey, secretKey), "", " ")
	return file
}

// providerConfigTemplate returns a ProviderConfigTemplate struct
// which contains info about the provider configuration including
// authentication fields.
func providerConfigTemplate(accessKey string, secretKey string) ProviderConfigTemplate {
	var awsProviderConfigTemplate ProviderConfigTemplate = ProviderConfigTemplate{
		Provider: map[string]interface{}{
			hostName: map[string]interface{}{
				"region":     "us-east-2",
				"access_key": accessKey,
				"secret_key": secretKey,
			},
		},
	}
	return awsProviderConfigTemplate
}

// baseInfraTemplate returns a byte slice that represents the base
// infrastructure to be deployed on the aws host
func baseInfraTemplate() []byte {
	randstr := randSeq(12)
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

func baseInfraConfigured(baseInfraFile string) bool {
	exists := true
	_, err := os.Stat(baseInfraFile)
	if err != nil {
		return !exists
	}
	return exists
}

func configureBaseInfra(baseInfraFile string) error {
	err := ioutil.WriteFile(baseInfraFile, baseInfraTemplate(), 0644)

	if err != nil {
		fmt.Println("error configureBaseInfra writing provider.tf.json for host", hostName)
		return err
	}

	err = TfApply(cliinit.HostPath(hostName))
	if err != nil {
		os.Remove(baseInfraFile)
		if strings.Contains(err.Error(), "NoCredentialProviders") {
			return fmt.Errorf("error: missing credentials for %v host", hostName)
		} else if strings.Contains(err.Error(), "InvalidClientTokenId") {
			return fmt.Errorf("error: invalid access_key for %v host", hostName)
		} else if strings.Contains(err.Error(), "SignatureDoesNotMatch") {
			return fmt.Errorf("error: invalid secret_key for %v host", hostName)
		}
	}

	return nil
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
