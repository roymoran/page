package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"builtonpage.com/main/cliinit"
	"github.com/hashicorp/terraform-exec/tfexec"
)

type AmazonWebServices struct {
	Infrastructure string
}

var AwsDefinition TerraformTemplate = TerraformTemplate{
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

func (aws AmazonWebServices) ConfigureHost(alias string) (bool, error) {
	hostName := "aws"
	providerId := alias
	hostPath := filepath.Join(cliinit.TfInstallPath, hostName)
	definitionFilePath, stateFilePath := aws.InstallTerraformPlugin(providerId, hostPath)
	tf, _ := tfexec.NewTerraform(hostPath, cliinit.TfExecPath)
	tf.Apply(context.Background(), tfexec.State(stateFilePath))

	provider := cliinit.ProviderConfig{
		Id:               providerId,
		Type:             "host",
		Alias:            alias,
		HostName:         hostName,
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      stateFilePath,
	}

	addProviderErr := cliinit.AddProvider(provider)

	return true, addProviderErr
}

func (aws AmazonWebServices) HostProviderDefinition() []byte {
	file, _ := json.MarshalIndent(AwsDefinition, "", " ")
	return file
}

func (aws AmazonWebServices) InstallTerraformPlugin(providerId string, hostPath string) (defPath string, statePath string) {
	definitionPath := filepath.Join(hostPath, providerId+".tf.json")
	stateDefinitionPath := filepath.Join(hostPath, providerId+".tfstate")

	_ = ioutil.WriteFile(definitionPath, aws.HostProviderDefinition(), 0644)
	tf, err := tfexec.NewTerraform(hostPath, cliinit.TfExecPath)
	if err != nil {
		log.Fatalln("error creating NewTerraform", hostPath, cliinit.TfInstallPath)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))

	if err != nil {
		fmt.Println(tf.Output(context.Background()))
		log.Fatalln("error initializing tf directory", hostPath, cliinit.TfInstallPath, err)
	}

	tf.Apply(context.Background(), tfexec.State(stateDefinitionPath))

	return definitionPath, stateDefinitionPath
}
