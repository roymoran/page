/*
Package definition defines methods for interacting wtih
and managing the yaml definition file for this cli
tool. Example of page definition file can be found
in from root directory at docs/README.md
*/
package definition

import (
	"io/ioutil"
	"log"
)

var defaultTemplate = `# version - Page config template version
version: "v0"
# template - uniform resource locator where 
# page template is located and accessible
template: "https://github.com/roymoran/page"
# unexpanded domain uses default registrar info/config
domain: "example.com"
# unexpanded host uses default host info/config
host: "page"
`

// WriteDefinitionFile writes the yaml file
// with default configurations at the location
// specified by the path. If no path is passed
// the current directory where the commnd was
// executed is assumed. A return of true signals
// that the file was writtent succesfully, otherwise
// false is returned.
func WriteDefinitionFile() bool {
	status := true

	writeErr := ioutil.WriteFile("page.yml", []byte(defaultTemplate), 0644)
	if writeErr != nil {
		log.Fatal("Failed to init page.yml")
		log.Fatal(writeErr)
		status = false
	}

	return status
}

// ReadDefinitionFile loads the yaml file into
// the program state. It tries to read the file
// from the current path given the filename. If
// no file is found in the current path it returns
// a file not found error.
func ReadDefinitionFile() {
	// TODO: Implement
}
