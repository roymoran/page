// entry point for cli
package main

import (
	"os"

	"builtonpage.com/main/command"
)

const appVersion = "v0.1.0"

func main() {
	// read args
	// sanitize args
	// pass args to commands handler, and execute command handler action
	// send response to standard out
	var _, args = os.Args[0], os.Args[1:]
	command.Handle(args)
}
