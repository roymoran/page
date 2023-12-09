package unit_tests

import (
	"testing"

	"pagecli.com/main/definition"
)

func TestGetRootDomain(t *testing.T) {
	tests := []struct {
		name         string
		inputURL     string
		expected     string
		expectingErr bool
	}{
		{"Test with www and https", "https://www.pagecli.com", "pagecli.com", false},
		{"Test with http and no www", "http://pagecli.com", "pagecli.com", false},
		{"Test with path", "www.pagecli.com/somepath", "pagecli.com", false},
		{"Test with multiple paths", "www.pagecli.com/somepath/somepath2", "pagecli.com", false},
		{"Test with subdomains", "sub1.sub2.pagecli.com", "pagecli.com", false},
		{"Test with no scheme", "pagecli.com", "pagecli.com", false},
		{"Test with empty input", "", "", true},
		{"Test with invalid URL", "http:///www.example.com", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := definition.GetRootDomain(tc.inputURL)
			if tc.expectingErr {
				if err == nil {
					t.Errorf("Expected an error for input '%s' but did not get one", tc.inputURL)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for input '%s' but got: %v", tc.inputURL, err)
				}
				if actual != tc.expected {
					t.Errorf("Expected root domain '%s' for input '%s', but got '%s'", tc.expected, tc.inputURL, actual)
				}
			}
		})
	}
}
