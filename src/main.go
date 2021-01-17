/*
Package main sole purpose is to read args and relay
to arguments to command handler for further processing
*/
package main

import (
	"fmt"
	"os"

	"builtonpage.com/main/commands"
)

func main() {
	output := commands.Handle(os.Args[1:])
	fmt.Print(output)
}
