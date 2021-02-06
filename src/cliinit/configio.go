package cliinit

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

func AddProvider(provider ProviderConfig) error {
	config, readErr := ReadConfigFile()
	if readErr != nil {
		log.Fatal("error reading cli config file", readErr)
		return readErr
	}

	ShouldChangeDefaultField(config.Providers, &provider)
	providers := append(config.Providers, provider)
	config.Providers = providers
	writeErr := WriteConfigFile(config)

	if writeErr != nil {
		log.Fatal("error writing cli config file", writeErr)
		return writeErr
	}

	return nil
}

// FindHostByAlias returns the name of the host given
// an alias
func FindHostByAlias(alias string) (string, error) {
	pageConfig, _ := ReadConfigFile()

	for _, provider := range pageConfig.Providers {
		if provider.Type != "host" {
			continue
		}

		if provider.Alias == alias {
			return provider.Name, nil
		}
	}

	return "", errors.New("err")
}

// FindDefaultAliasForHost returns the alias for the default
// host provider
func FindDefaultAliasForHost(hostName string) (string, error) {
	pageConfig, _ := ReadConfigFile()

	for _, provider := range pageConfig.Providers {
		if provider.Type != "host" {
			continue
		}

		if provider.Name == hostName && provider.Default {
			return provider.Alias, nil
		}
	}

	return "", errors.New("err")
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

// IsDefault changes the 'Default' field of ProviderConfig
// to false if there already exists a default provider
// for the host
func ShouldChangeDefaultField(providers []ProviderConfig, provider *ProviderConfig) {
	for _, p := range providers {
		if p.Name == provider.Name && p.Default {
			// there already exists a default provider
			// for host
			provider.Default = false
		}
	}
}
