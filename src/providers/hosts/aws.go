package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"builtonpage.com/main/cliinit"
	"github.com/hashicorp/terraform-exec/tfexec"
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
}

func (aws AmazonWebServices) Deploy() bool {
	return true
}

func (aws AmazonWebServices) ConfigureHost(alias string, definitionFilePath string, stateFilePath string) (bool, error) {
	fmt.Println("entered aws ConfigureHost")
	hostName := "aws"
	hostPath := filepath.Join(cliinit.TfInstallPath, hostName)
	tf, _ := tfexec.NewTerraform(hostPath, cliinit.TfExecPath)
	tf.Apply(context.Background(), tfexec.State(stateFilePath))

	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		HostName:         hostName,
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      stateFilePath,
	}

	addProviderErr := cliinit.AddProvider(provider)
	fmt.Println("finished aws ConfigureHost")
	return true, addProviderErr
}

func (aws AmazonWebServices) HostProviderDefinition() []byte {
	file, _ := json.MarshalIndent(AwsProviderDefinition, "", " ")
	return file
}
