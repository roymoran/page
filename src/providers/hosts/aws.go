package providers

import (
	"encoding/json"
	"fmt"

	"builtonpage.com/main/cliinit"
)

type AmazonWebServices struct {
	Infrastructure string
}

var AwsProviderDefinition TerraformTemplate = TerraformTemplate{
	Terraform: RequiredProviders{
		RequiredProvider: map[string]Provider{
			"aws": {
				Source:  "hashicorp/aws",
				Version: "3.25.0",
			},
		},
	},
	Provider: map[string]ProviderConfig{
		"aws": {
			Profile: "default",
			Region:  "us-east-2",
		},
	},
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

func (aws AmazonWebServices) ConfigureHost() bool {
	fmt.Println("configured aws host")
	return true
}

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

func (aws AmazonWebServices) HostProviderDefinition() []byte {
	file, _ := json.MarshalIndent(AwsProviderDefinition, "", " ")
	return file
}
