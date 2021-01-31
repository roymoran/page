package providers

import "fmt"

type Azure struct {
}

var AzureTerraformProvider = `
terraform {
	required_providers {
		azurerm = {
			source = "hashicorp/azurerm"
			version = "=2.44.0"
		}
	}
}
`

func (a Azure) Deploy() bool {
	return true
}

func (a Azure) ConfigureHost(alias string, definitionFilePath string, stateFilePath string) (bool, error) {
	fmt.Println("configured azure host")
	return true, nil
}

func (a Azure) HostProviderDefinition() []byte {
	return []byte{}
}
