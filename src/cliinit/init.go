/*
Package cliinit provides types and methods to
perform intial configuration of cli including
creating files/directories to support cli
functionality.
*/
package cliinit

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfinstall"
)

var installDir, _ = os.UserHomeDir()
var pageCliPath string = filepath.Join(installDir, ".pagecli")

// TfInstallPath returns the path to the
// directory containing terraform binary
var TfInstallPath string = filepath.Join(pageCliPath, "tf")

// TfExecPath returns the path to the
// terraform binary
var TfExecPath string = filepath.Join(TfInstallPath, "terraform")

// ConfigPath returns the path to the
// page cli config.json file. Which contains
// configuration details for this cli tool
var configPath string = filepath.Join(pageCliPath, "config.json")

// ProvidersPath returns the path to the
// 'provider' directory a nested directory
// inside .pagecli which is the root to all
// hosts/registrars
var ProvidersPath string = filepath.Join(TfInstallPath, "providers")

// HostPath returns the path to a specific host directory
// with the given 'hostName' which contains terraform
// configuration files
var HostPath func(hostName string) string = func(hostName string) string { return filepath.Join(ProvidersPath, hostName) }

// ProviderAliasPath returns the path to a specific alias
// directory for a provider
var ProviderAliasPath func(providerName string, alias string) string = func(providerName string, alias string) string {
	return filepath.Join(ProvidersPath, providerName, alias)
}
var exactTfVersion string = "0.14.5"

var initialPageConfig PageConfig = PageConfig{
	TfPath:       TfInstallPath,
	TfExecPath:   "",
	TFVersion:    exactTfVersion,
	Providers:    []ProviderConfig{},
	ConfigStatus: false,
}

// CliInit creates the required directories
// and installs required executables for the
// cli
func CliInit() {
	dirErr := os.MkdirAll(TfInstallPath, os.ModePerm)
	if dirErr != nil {
		log.Fatal("CliInit error. Error creating tf install path.", dirErr)
	}

	configError := WriteConfigFile(initialPageConfig)

	if configError != nil {
		log.Fatal("CliInit error. Error creating config.json.", configError)
	}

	execPath, installErr := InstallTerraform()

	if installErr != nil {
		log.Fatal("CliInit error. Error installing terraform.", installErr)
	}

	initialPageConfig.ConfigStatus = true
	initialPageConfig.TfExecPath = execPath
	configError = WriteConfigFile(initialPageConfig)

	if configError != nil {
		log.Fatal("CliInit error. Error setting InitialConfig to true.", configError)
	}
}

// CliInitialized checks whether the
// directories and required executables
// have been installed for the cli to
// function properly
func CliInitialized() bool {
	initialized := false
	configData, fileErr := ioutil.ReadFile(configPath)

	if fileErr != nil {
		return initialized
	}

	var config PageConfig
	unmarshalErr := json.Unmarshal(configData, &config)

	if unmarshalErr != nil {
		return initialized
	}

	return config.ConfigStatus
}

// InstallTerraform installs the terraform binary
// with the version specified by 'exactTfVersion'
// in the directory specified by 'TfInstallPath'
func InstallTerraform() (string, error) {
	execPath, installErr := tfinstall.Find(context.Background(), tfinstall.ExactVersion(exactTfVersion, TfInstallPath))

	if installErr != nil {
		log.Fatal("InstallTerraform error.", installErr)
		return "", installErr
	}

	return execPath, nil
}
