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
	output := make(chan string)
	go commands.Handle(os.Args, output)
	for message := range output {
		fmt.Println(message)
	}
}
