/*
Each of these tests should complete in a few seconds
or less. These tests should not depend on external resources
that are long-running tasks like standing up infrastructure on
a cloud provider. But they may depend on external resources
that are short-lived like the file system or network.
*/
package integration_tests

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"pagecli.com/main/cliinit"
	"pagecli.com/main/providers"
)

// Test that the terraform binary is installed
// successfully. This is typically done on the
// intial run of the CLI.
func TestTerraformInstalled(t *testing.T) {
	setup()
	defer cleanup()

	// CliInit() should already have been called
	// through setup()

	// now look for terraform binary
	if !fileExists(cliinit.TfExecPath) {
		t.Errorf("Expected terraform binary at '%v' instead got nothing", cliinit.TfExecPath)
	}
}

// Test that the terraform namecheap registrar provider
// is installed in the expected path. This test emulates
// the execution of the following command:
// $ page registrar add namecheap
// Note: Since the command requires input from the user,
// the input is sent to stdin via a temporary file defined
// as part of the test setup.
func TestTerraformNamecheapRegistrarProviderInstalled(t *testing.T) {
	setup()
	defer cleanup()

	input := []byte("nc\nrmoran20\napiKey1\n")
	tmpfile, err := os.CreateTemp("", "gotest")

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(input); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile

	output := createOutputChannel()
	defer close(output)
	provider, _ := providers.SupportedProviders.Providers["registrar"]
	providerInstallErr := provider.Add("namecheap", output)
	namecheapProviderPath := filepath.Join(cliinit.ProvidersPath, ".terraform", "providers", "registry.terraform.io", "namecheap", "namecheap")

	if providerInstallErr != nil {
		// If providerInstallErr reports EOF, then
		// the test is failing because the user input
		// is not sufficiently provided
		t.Errorf("provider.Add('namecheap', output) failed with '%v'", providerInstallErr)
	}

	if !fileExists(namecheapProviderPath) {
		t.Errorf("Expected namecheap provider directory at '%v'", namecheapProviderPath)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
}

func TestTerraformAwsHostProviderInstalled(t *testing.T) {
	setup()
	defer cleanup()

	input := []byte("awsmain\nkey1\nkey2\n")
	tmpfile, err := os.CreateTemp("", "gotest")

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(input); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile

	output := createOutputChannel()
	defer close(output)
	provider, _ := providers.SupportedProviders.Providers["host"]
	providerInstallErr := provider.Add("aws", output)
	awsProviderPath := filepath.Join(cliinit.ProvidersPath, ".terraform", "providers", "registry.terraform.io", "hashicorp", "aws")

	if providerInstallErr != nil {
		// If providerInstallErr reports EOF, then
		// the test is failing because the user input
		// is not sufficiently provided
		t.Errorf("provider.Add('aws', output) failed with '%v'", providerInstallErr)
	}

	if !fileExists(awsProviderPath) {
		t.Errorf("Expected aws provider directory at '%v'", awsProviderPath)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
}

func setup() {
	// setup by initializing the cli
	// which creates the page cli directory
	// and downloads terraform binary
	cliinit.CliInit()
}

func cleanup() {
	// remove page cli directory for
	// subsequent tests
	os.RemoveAll(cliinit.PageCliPath)
}

func createOutputChannel() chan string {
	output := make(chan string)

	// Start a goroutine to read from the output channel
	go func() {
		for msg := range output {
			fmt.Println(msg)
		}
	}()

	return output
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if err == nil {
		return true // File exists
	}
	if os.IsNotExist(err) {
		return false // File does not exist
	}
	return false // Error occurred while checking file existence
}
