package terraformutils

import "os"

func ResourcesConfigured(infraFile string) bool {
	exists := true
	_, err := os.Stat(infraFile)
	if err != nil {
		return !exists
	}
	return exists
}
