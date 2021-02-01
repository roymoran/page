/*
Package main sole purpose is to read args and relay
them to command handler for further processing
*/
package main

import (
	"fmt"
	"os"

	"builtonpage.com/main/commands"
)

func main() {
	// TODO: Allow for intermediate output as command
	// is processing, mainly for long-running commands
	output := commands.Handle(os.Args)
	fmt.Print(output)
}
