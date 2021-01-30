/*
Package init provides types and methods to
perform intial configuration of cli including
creating files/directories to support cli
functionality.
*/
package cliinit

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfinstall"
)

var installDir, _ = os.UserHomeDir()
var pageCliPath string = filepath.Join(installDir, ".pagecli")
var tfInstallPath string = filepath.Join(pageCliPath, "tf")
var configPath string = filepath.Join(pageCliPath, "config.json")
var exactTfVersion string = "0.14.5"

var initialPageConfig PageConfigJson = PageConfigJson{
	TfPath:       tfInstallPath,
	TFVersion:    exactTfVersion,
	Providers:    []ProviderConfig{},
	ConfigStatus: false,
}

// CliInit creates the required directories
// and installs required executables for the
// cli
func CliInit() {
	dirErr := os.MkdirAll(tfInstallPath, os.ModePerm)
	if dirErr != nil {
		log.Fatal("CliInit error. Error creating tf install path.", dirErr)
	}

	file, _ := json.MarshalIndent(initialPageConfig, "", " ")
	configError := ioutil.WriteFile(configPath, file, 0644)

	if configError != nil {
		log.Fatal("CliInit error. Error creating config.json.", dirErr)
	}

	InstallTerraform()

	initialPageConfig.ConfigStatus = true
	file, _ = json.MarshalIndent(initialPageConfig, "", " ")
	configError = ioutil.WriteFile(configPath, file, 0644)

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
		log.Println(fileErr)
		return initialized
	}

	var config PageConfigJson
	unmarshalErr := json.Unmarshal(configData, &config)

	if unmarshalErr != nil {
		log.Println(unmarshalErr)
		return initialized
	}

	return config.ConfigStatus
}

func InstallTerraform() {
	execPath, installErr := tfinstall.Find(context.Background(), tfinstall.ExactVersion(exactTfVersion, tfInstallPath))
	fmt.Println("execPath", execPath)

	if installErr != nil {
		log.Fatal("InstallTerraform error.", installErr)
	}
}
