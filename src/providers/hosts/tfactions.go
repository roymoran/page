package hosts

import (
	"context"
	"fmt"

	"builtonpage.com/main/cliinit"
	"github.com/hashicorp/terraform-exec/tfexec"
)

// TfInit initializes a new terraform directory
// in the specified initPath and installs plugins
// specified in provider.tf.json
func TfInit(initPath string) error {
	tf, err := tfexec.NewTerraform(initPath, cliinit.TfExecPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// TfApply runs the terraform apply command in the specified
// applyPath
func TfApply(applyPath string) error {
	tf, err := tfexec.NewTerraform(applyPath, cliinit.TfExecPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = tf.Apply(context.Background())
	if err != nil {
		// TODO: Log TF error
		fmt.Println(err)
		return err
	}

	return nil
}

// TfOutput returns the output for the variable
// given the 'name'
func TfOutput(path string, identifier string) (string, error) {
	tf, err := tfexec.NewTerraform(path, cliinit.TfExecPath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	outMeta, err := tf.Output(context.Background())

	if err != nil {
		// TODO: Log TF error
		fmt.Println(err)
		return "", err
	}

	return string(outMeta[identifier].Value), nil
}

// TFResourceNameGenerator generates a string which is used
// as the name for a resource defined in a terraform template
// improper naming of a resource inside a terraform template
// will throw an error.
func TFResourceNameGenerator(n int) string {
	// A name must start with a letter or underscore and may contain only letters,
	// digits, underscores, and dashes.
	return ""
}
