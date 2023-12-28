package hosts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-exec/tfexec"
	"pagecli.com/main/cliinit"
	"pagecli.com/main/progress"
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

	// Stream Terraform's output to stdout in debug mode
	if os.Getenv("PAGE_CLI_DEBUG") == "true" {
		tf.SetStdout(os.Stdout)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// TfApply runs the terraform apply command in the specified
// applyPath and emits a progress message.
func TfApply(infrastructureCheckMessage string, progressMessageSequence []string, timeout time.Duration) error {
	var err error = nil
	const maxRetries = 3

	tf, err := tfexec.NewTerraform(cliinit.ProvidersPath, cliinit.TfExecPath)
	if err != nil {
		return err
	}

	if os.Getenv("PAGE_CLI_DEBUG") == "true" {
		// Stream Terraform's output to stdout in debug mode
		tf.SetStdout(os.Stdout)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		ticker := time.NewTicker(200 * time.Millisecond)
		go func() {
			for {
				select {
				case <-ctx.Done():
					ticker.Stop()
					return
				case <-ticker.C:
					// Cycle through the progress animation
					next := progress.ProgressSequence[0]
					message := infrastructureCheckMessage + progressMessageSequence[0] + next
					progress.ProgressSequence = append(progress.ProgressSequence[1:], progress.ProgressSequence[0]) // Rotate the slice
					fmt.Printf("\r\033[2K%s", message)                                                              // Clear the line, then overwrite it                                                       // Overwrite the current line with the new progress indicator
				}
			}
		}()

		err = tf.Apply(ctx)
		if err == nil {
			// If apply succeeds, exit the loop
			cancel()
			break
		}

		cancel()                    // Cancel the current context
		time.Sleep(time.Second * 1) // Wait a second before retrying
	}

	var progressResultIndicator, progressResultMessage string
	if err != nil {
		progressResultMessage = progress.ProgressFailedMessage
		progressResultIndicator = progress.ProgressFailed
	} else {
		progressResultMessage = progressMessageSequence[1]
		progressResultIndicator = progress.ProgressComplete
	}

	// Add a short delay to allow the goroutine to stop the ticker before
	// printing the final progress indicator
	time.Sleep(100 * time.Millisecond)

	message := infrastructureCheckMessage + progressResultMessage + progressResultIndicator
	fmt.Printf("\r\033[2K%s\n", message) // Clear the line, then overwrite it and move to the next line

	return err
}

// TfApplyWithTarget runs the terraform apply command on
// the specified target and emits a progress message
func TfApplyWithTarget(infrastructureCheckMessage string, progressMessageSequence []string, timeout time.Duration, targets []string) error {
	var err error = nil
	const maxRetries = 3

	tf, err := tfexec.NewTerraform(cliinit.ProvidersPath, cliinit.TfExecPath)
	if err != nil {
		return err
	}

	if os.Getenv("PAGE_CLI_DEBUG") == "true" {
		// Stream Terraform's output to stdout in debug mode
		tf.SetStdout(os.Stdout)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		ticker := time.NewTicker(200 * time.Millisecond)
		go func() {
			for {
				select {
				case <-ctx.Done():
					ticker.Stop()
					return
				case <-ticker.C:
					// Cycle through the progress animation
					next := progress.ProgressSequence[0]
					message := infrastructureCheckMessage + progressMessageSequence[0] + next
					progress.ProgressSequence = append(progress.ProgressSequence[1:], progress.ProgressSequence[0]) // Rotate the slice
					fmt.Printf("\r\033[2K%s", message)                                                              // Clear the line, then overwrite it                                                       // Overwrite the current line with the new progress indicator
				}
			}
		}()

		var applyOptions []tfexec.ApplyOption
		for _, target := range targets {
			applyOptions = append(applyOptions, tfexec.Target(target))
		}

		err = tf.Apply(ctx, applyOptions...)
		if err == nil {
			// If apply succeeds, exit the loop
			cancel()
			break
		}

		cancel() // Cancel the current context
		// TODO: Can a failed terraform apply be detrimental to
		// the page up process? Will a failed apply here cause
		// future applies to fail? What's the best approach here?
		// Should we use terraform destroy to clean up any partially created resource?
		time.Sleep(time.Second * 1) // Wait a second before retrying
	}

	var progressResultIndicator, progressResultMessage string
	if err != nil {
		progressResultMessage = progress.ProgressFailedMessage
		progressResultIndicator = progress.ProgressFailed
	} else {
		progressResultMessage = progressMessageSequence[1]
		progressResultIndicator = progress.ProgressComplete
	}

	// Add a short delay to allow the goroutine to stop the ticker before
	// printing the final progress indicator
	time.Sleep(100 * time.Millisecond)

	message := infrastructureCheckMessage + progressResultMessage + progressResultIndicator
	fmt.Printf("\r\033[2K%s\n", message) // Clear the line, then overwrite it and move to the next line

	return err
}

// TfOutput returns the output for the variable
// given the 'name'
func TfSingleOutput[V string | int](key string, identifier string) (V, error) {
	var value V
	tf, err := tfexec.NewTerraform(cliinit.ProvidersPath, cliinit.TfExecPath)
	if err != nil {
		fmt.Println(err)
		return value, err
	}

	outMeta, err := tf.Output(context.Background())

	if err != nil {
		// TODO: Log TF error
		fmt.Println(err)
		return value, err
	}

	var mapValues map[string]V = make(map[string]V)
	// json decode values
	err = json.Unmarshal(outMeta[key].Value, &mapValues)
	value = mapValues[identifier]

	return value, nil
}
