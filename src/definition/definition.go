/*
Package definiton defines methods for interacting wtih
and managing the yaml definition file for this cli
tool. Example of page definition file can be found
in from root directory at docs/README.md
*/
package definiton

// WriteDefinitionFile writes the yaml file
// with default configurations at the location
// specified by the path. If no path is passed
// the current directory where the commnd was
// executed is assumed.
func WriteDefinitionFile() {
	// TODO: Implement
}

// ReadDefinitionFile loads the yaml file into
// the program state. It tries to read the file
// from the current path given the filename. If
// no file is found in the current path it returns
// a file not found error.
func ReadDefinitionFile() {
	// TODO: Implement
}
