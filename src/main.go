/*
Package main sole purpose is to read args and relay
to arguments to command handler for further processing
*/
package main

import (
	"fmt"
	"os"

	"builtonpage.com/main/command"
)

func main() {
	// read args
	// sanitize args
	// pass args to commands handler, and execute command handler action
	// send response to standard out
	// var _, args = os.Args[0], os.Args[1:]
	output := command.Handle(os.Args[1:])
	fmt.Println(output)
}
