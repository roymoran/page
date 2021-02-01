package tests

import (
	"strings"
	"testing"

	"builtonpage.com/main/commands"
)

func TestNoneCommand(t *testing.T) {
	args := []string{}
	expectedSubstring := "page version v"
	output := make(chan string)
	go commands.Handle(args, output)
	if !strings.Contains(<-output, expectedSubstring) {
		t.Errorf("Expected output to contain '%v' instead got '%v'", expectedSubstring, output)
	}
}

// func TestInitCommand(t *testing.T) {
// 	args := []string{"init"}
// 	command.Handle(args)
// 	t.Error("error")
// }
