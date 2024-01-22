/*
Package definition defines methods for interacting wtih
and managing the yaml definition file for this cli
tool. Example of page definition file can be found
in from root directory at docs/README.md
*/
package definition

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var defaultTemplate = `# This configuration file defines your site to be deployed
# The comments displayed before each property should be informative
# enough to get your site deployed. If you run into a problem
# you can create an issue at https://github.com/roymoran/page.

# page config template version
version: "0"
# specify a supported host name or an alias
# supported hosts can be found with 'page conf host list'
host: "aws"
# specify a supported registrar name or an alias
# supported registrars can be found with 'page conf registrar list'
registrar: "namecheap"
# specify the domain name for your site. The registrar
# you specified above must own the domain name.
# Only specify a top-level domain name.
domain: "example.com"
# template - a path to static assets to be published, 
# either a local file path or public git url.
template: "https://gitlab.com/page-templates/placeholders/comingsoon.git"
`

type TemplateSource int

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type PageDefinition struct {
	Template  string
	Registrar string
	Host      string
	Domain    string
	Version   string
}

// TemplateSource is the source of
// the static site template.
const (
	GitURL TemplateSource = iota
	FilePath
)

type PageDefinitionConfig struct {
	TemplateSource TemplateSource

	// validation indicators for the fields in the
	// PageDefinition struct
	ValidTemplate  bool
	ValidRegistrar bool
	ValidHost      bool
	ValidDomain    bool
	// if any of the fields are invalid, lets indicate why
	InvalidFields map[string]string
}

// WriteDefinitionFile writes the yaml file
// with default configurations at the location
// specified by the path. If no path is passed
// the current directory where the commnd was
// executed is assumed. A return of true signals
// that the file was writtent succesfully, otherwise
// false is returned.
func WriteDefinitionFile() error {
	writeErr := os.WriteFile("page.yml", []byte(defaultTemplate), 0644)
	if writeErr != nil {
		log.Fatal("Failed to write page.yml", writeErr)
		return writeErr
	}

	return nil
}

// ReadDefinitionFile loads the yaml file into
// the program state. It tries to read the file
// from the current path given the filename. If
// no file is found in the current path it returns
// a file not found error.
func ReadDefinitionFile() (PageDefinition, error) {
	t := PageDefinition{}
	data, err := os.ReadFile("page.yml")
	if err != nil {
		return t, fmt.Errorf("unable to find page.yml in current directory")
	}

	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return t, fmt.Errorf("error parsing page.yml. " + err.Error())
	}

	return t, nil
}

// ProccessDefinitionFile reads the values from the
// PageDefinition struct and determines things about
// the values provided. For example, if the template
// is a git url or a local file path.
func ProccessDefinitionFile(pd *PageDefinition) (PageDefinitionConfig, error) {
	pageDefinitionConfig := PageDefinitionConfig{
		ValidTemplate:  true,
		ValidRegistrar: true,
		ValidHost:      true,
		ValidDomain:    true,
		InvalidFields:  make(map[string]string),
	}

	// if any of the fields are empty mark the field as invalid
	if pd.Template == "" {
		pageDefinitionConfig.ValidTemplate = false
		pageDefinitionConfig.InvalidFields["template"] = "template field is empty"
	}
	if pd.Registrar == "" {
		pageDefinitionConfig.ValidRegistrar = false
		pageDefinitionConfig.InvalidFields["registrar"] = "registrar field is empty"
	}
	if pd.Host == "" {
		pageDefinitionConfig.ValidHost = false
		pageDefinitionConfig.InvalidFields["host"] = "host field is empty"
	}
	if pd.Domain == "" {
		pageDefinitionConfig.ValidDomain = false
		pageDefinitionConfig.InvalidFields["domain"] = "domain field is empty"
	}

	// if any of the fields are invalid, return the config and error
	if !pageDefinitionConfig.ValidTemplate || !pageDefinitionConfig.ValidRegistrar || !pageDefinitionConfig.ValidHost || !pageDefinitionConfig.ValidDomain {
		return pageDefinitionConfig, templateFieldErrors(pageDefinitionConfig.InvalidFields)
	}

	// check if the template is a valid source (git url or file path)
	if isGitURL(pd.Template) {
		pageDefinitionConfig.TemplateSource = GitURL
	} else {
		// assume file path now check if it is a valid path
		// on the system
		_, err := os.Stat(pd.Template)
		if err != nil {
			pageDefinitionConfig.ValidTemplate = false
			pageDefinitionConfig.InvalidFields["template"] = "template field is not a valid git url or file path"
			return pageDefinitionConfig, templateFieldErrors(pageDefinitionConfig.InvalidFields)
		}

		pageDefinitionConfig.TemplateSource = FilePath
	}

	rootDomain, domainErr := GetRootDomain(pd.Domain)
	if domainErr != nil {
		pageDefinitionConfig.ValidDomain = false
		pageDefinitionConfig.InvalidFields["domain"] = "domain field is not a valid domain name"
		return pageDefinitionConfig, templateFieldErrors(pageDefinitionConfig.InvalidFields)
	}

	// update the domain field with the root domain
	pd.Domain = rootDomain

	return pageDefinitionConfig, nil
}

// getRootDomain takes a URL and attempts to return the root domain.
func GetRootDomain(inputURL string) (string, error) {

	if inputURL == "" {
		return "", fmt.Errorf("empty input url")
	}

	// Parse the URL input.
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	// Extract the hostname from the URL.
	hostname := u.Hostname()

	// if host name is empty, we may have a url without a scheme
	// so lets try to parse removing any paths then continue to split
	// the hostname by dots
	if hostname == "" {
		u, err := url.Parse(strings.Split(inputURL, "/")[0])
		if err != nil {
			return "", err
		}
		hostname = u.Path
	}

	// Split the hostname by dots.
	parts := strings.Split(hostname, ".")

	// Check if we have at least a domain and a TLD.
	if len(parts) >= 2 {
		length := len(parts)
		// Construct the root domain by taking the last two parts.
		rootDomain := parts[length-2] + "." + parts[length-1]
		return rootDomain, nil
	}

	return "", fmt.Errorf("unable to extract root domain from: %s", inputURL)
}

func isGitURL(template string) bool {
	// Checking for HTTPS, HTTP, and Git protocol URLs
	if strings.HasPrefix(template, "https://") ||
		strings.HasPrefix(template, "http://") ||
		strings.HasPrefix(template, "git://") {
		return true
	}

	// Checking for SSH URLs
	if strings.HasPrefix(template, "ssh://git@") ||
		strings.HasPrefix(template, "git@") {
		return true
	}

	return false
}

func templateFieldErrors(invalidFields map[string]string) error {
	var errStr strings.Builder
	errStr.WriteString("invalid page definition file:\n")

	for field, err := range invalidFields {
		errStr.WriteString(fmt.Sprintf("%s: %s\n", field, err))
	}

	return fmt.Errorf(errStr.String())
}
