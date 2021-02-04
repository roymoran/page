package providers

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
		fmt.Println(tf.Output(context.Background()))
		return err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))

	if err != nil {
		fmt.Println(tf.Output(context.Background()))
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
		fmt.Println(tf.Output(context.Background()))
		return err
	}

	err = tf.Apply(context.Background())

	if err != nil {
		fmt.Println(tf.Output(context.Background()))
		fmt.Println(err)
		return err
	}

	return nil
}
