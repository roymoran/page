/*
Package definition defines methods for interacting wtih
and managing the yaml definition file for this cli
tool. Example of page definition file can be found
in from root directory at docs/README.md
*/
package definition

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var defaultTemplate = `# version - page config template version
version: "0"
# specify a supported host name or an alias
host: "aws"
# specify a supported registrar name or an alias
registrar: "namecheap"
# specify the domain name for your site. The registrar
# specified above must own the domain name. Only specify
# a top-level domain name. 
domain: "example.com"
# template - a url of a git repo containing static assets
# to be hosted. url should be accessible from machine running
# 'page up'
template: "https://github.com/roymoran/index"
`

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type PageDefinition struct {
	Template  string
	Registrar string
	Host      string
	Domain    string
	Version   string
}

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
func ReadDefinitionFile() (PageDefinition, error) {
	t := PageDefinition{}
	data, err := ioutil.ReadFile("page.yml")
	if err != nil {
		return t, fmt.Errorf("unable to find page.yml in current directory")
	}

	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return t, fmt.Errorf("error parsing page.yml. " + err.Error())
	}

	return t, nil
}
