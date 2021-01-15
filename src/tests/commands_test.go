package tests

import (
	"strings"
	"testing"

	"builtonpage.com/main/command"
)

func TestNoneCommand(t *testing.T) {
	args := []string{}
	expectedSubstring := "page version v"
	var output string = command.Handle(args)
	if !strings.Contains(output, expectedSubstring) {
		t.Errorf("Expected output to contain '%v' instead got '%v'", expectedSubstring, output)
	}
}

// func TestInitCommand(t *testing.T) {
// 	args := []string{"init"}
// 	command.Handle(args)
// 	t.Error("error")
// }
