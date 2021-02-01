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

var InstallDir, _ = os.UserHomeDir()
var PageCliPath string = filepath.Join(InstallDir, ".pagecli")
var TfInstallPath string = filepath.Join(PageCliPath, "tf")
var TfExecPath string = filepath.Join(TfInstallPath, "terraform")
var ConfigPath string = filepath.Join(PageCliPath, "config.json")
var ExactTfVersion string = "0.14.5"

var initialPageConfig PageConfig = PageConfig{
	TfPath:       TfInstallPath,
	TfExecPath:   "",
	TFVersion:    ExactTfVersion,
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
	configData, fileErr := ioutil.ReadFile(ConfigPath)

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

func InstallTerraform() (string, error) {
	execPath, installErr := tfinstall.Find(context.Background(), tfinstall.ExactVersion(ExactTfVersion, TfInstallPath))

	if installErr != nil {
		log.Fatal("InstallTerraform error.", installErr)
		return "", installErr
	}

	return execPath, nil
}
