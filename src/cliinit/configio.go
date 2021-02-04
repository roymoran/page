package cliinit

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func AddProvider(provider ProviderConfig) error {
	config, readErr := ReadConfigFile()
	if readErr != nil {
		log.Fatal("error reading cli config file", readErr)
		return readErr
	}

	providers := append(config.Providers, provider)
	config.Providers = providers
	writeErr := WriteConfigFile(config)

	if writeErr != nil {
		log.Fatal("error reading cli config file", writeErr)
		return writeErr
	}

	return nil
}

func ReadConfigFile() (PageConfig, error) {
	var config PageConfig
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("error reading cli config file", err)
		return PageConfig{}, err
	}
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal("error parsing config file", err)
		return PageConfig{}, err
	}

	return config, nil
}

func WriteConfigFile(config PageConfig) error {
	file, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		log.Fatal("error parsing config file", err)
		return err
	}
	// TODO: perm os.FileMode? 0644?
	err = ioutil.WriteFile(configPath, []byte(file), 0644)

	if err != nil {
		log.Fatal("error writing config file", err)
		return err
	}

	return nil
}
