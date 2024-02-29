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
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"pagecli.com/main/logging"
)

var installDir, _ = os.UserHomeDir()
var PageCliPath string = filepath.Join(installDir, ".pagecli")

// Path to website assets
var SitePath string = filepath.Join(PageCliPath, "sites")
var SiteFilesPath string = filepath.Join(SitePath, "files")
var SiteCertsPath string = filepath.Join(SitePath, "certs")

// TfInstallPath returns the path to the
// directory containing terraform binary
var TfInstallPath string = filepath.Join(PageCliPath, "tf")

// TfExecPath returns the path to the
// terraform binary
var TfExecPath string = filepath.Join(TfInstallPath, "terraform")

// ConfigPath returns the path to the
// page cli config.json file. Which contains
// configuration details for this cli tool
var ConfigPath string = filepath.Join(PageCliPath, "config.json")

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

var ModuleTemplatePath func(providerType string, alias string) string = func(providerType string, alias string) string {
	return filepath.Join(ProvidersPath, providerType+"_"+alias+".tf.json")
}
var exactTfVersion string = "1.6.4"

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
	logMessage := ""
	dirErr := os.MkdirAll(TfInstallPath, os.ModePerm)
	if dirErr != nil {
		logMessage = fmt.Sprint("CliInit error. Error creating tf install path.", dirErr)

		logging.LogException(logMessage, true)
		log.Fatal(logMessage)
	}

	configError := writeConfigFile(initialPageConfig)

	if configError != nil {
		logMessage = fmt.Sprint("CliInit error. Error creating config.json.", configError)

		logging.LogException(logMessage, true)
		log.Fatal(logMessage)
	}

	execPath, installErr := InstallTerraform()

	if installErr != nil {
		logMessage = fmt.Sprint("CliInit error. Error installing terraform.", installErr)

		logging.LogException(logMessage, true)
		log.Fatal(logMessage)
	}

	initialPageConfig.ConfigStatus = true
	initialPageConfig.TfExecPath = execPath
	configError = writeConfigFile(initialPageConfig)

	if configError != nil {
		logMessage = fmt.Sprint("CliInit error. Error setting InitialConfig to true.", configError)

		logging.LogException(logMessage, true)
		log.Fatal(logMessage)
	}
}

// CliInitialized checks whether the
// directories and required executables
// have been installed for the cli to
// function properly
func CliInitialized() bool {
	initialized := false
	configData, fileErr := os.ReadFile(ConfigPath)

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
	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion("1.6.4")),
		InstallDir: TfInstallPath,
	}

	execPath, installErr := installer.Install(context.Background())

	if installErr != nil {
		log.Fatal("InstallTerraform error.", installErr)
		return "", installErr
	}

	return execPath, nil
}
