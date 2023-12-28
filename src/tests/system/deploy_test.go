package unit_tests

import (
	"strings"
	"testing"

	"pagecli.com/main/commands"
)

func TestPlaceholder(t *testing.T) {
	args := []string{"up"}
	expectedSubstring := "page cli v"
	output := make(chan string)
	go commands.Handle(args, output)
	if !strings.Contains(<-output, expectedSubstring) {
		t.Errorf("Expected output to contain '%v' instead got '%v'", expectedSubstring, output)
	}
}
