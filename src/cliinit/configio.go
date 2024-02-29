package cliinit

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

func AddProvider(provider ProviderConfig) error {
	config, readErr := readConfigFile()
	if readErr != nil {
		log.Fatal("error reading cli config file", readErr)
		return readErr
	}

	shouldChangeDefaultField(config.Providers, &provider)
	providers := append(config.Providers, provider)
	config.Providers = providers
	writeErr := writeConfigFile(config)

	if writeErr != nil {
		log.Fatal("error writing cli config file", writeErr)
		return writeErr
	}

	return nil
}

// FindAllHostAliases returns the name of the host given
// an alias
func FindAllHostAliases() ([]string, error) {
	pageConfig, readConfigErr := readConfigFile()
	aliases := []string{}
	defaultMessage := ""

	if readConfigErr != nil {
		return aliases, readConfigErr
	}

	for _, provider := range pageConfig.Providers {
		if provider.Type != "host" {
			continue
		}

		if provider.Default {
			defaultMessage = "*default"
		}

		aliases = append(aliases, provider.Alias+" "+"("+provider.Name+") "+defaultMessage)
		defaultMessage = ""
	}

	return aliases, nil
}

// FindAllRegistrarAliases returns the name of the host given
// an alias
func FindAllRegistrarAliases() ([]string, error) {
	pageConfig, readConfigErr := readConfigFile()
	aliases := []string{}
	defaultMessage := ""

	if readConfigErr != nil {
		return aliases, readConfigErr
	}

	for _, provider := range pageConfig.Providers {
		if provider.Type != "registrar" {
			continue
		}

		if provider.Default {
			defaultMessage = "*default"
		}

		aliases = append(aliases, provider.Alias+" "+"("+provider.Name+") "+defaultMessage)
		defaultMessage = ""
	}

	return aliases, nil
}

// FindHostByAlias returns the name of the host given
// an alias
func FindHostByAlias(alias string) (string, error) {
	pageConfig, _ := readConfigFile()

	for _, provider := range pageConfig.Providers {
		if provider.Type != "host" {
			continue
		}

		if provider.Alias == alias {
			return provider.Name, nil
		}
	}

	return "", errors.New("no host found for alias")
}

// FindRegistrarByAlias returns the name of the host given
// an alias
func FindRegistrarByAlias(alias string) (string, error) {
	pageConfig, _ := readConfigFile()

	for _, provider := range pageConfig.Providers {
		if provider.Type != "registrar" {
			continue
		}

		if provider.Alias == alias {
			return provider.Name, nil
		}
	}

	return "", errors.New("no registrar found for alias")
}

// FindDefaultAliasForHost returns the alias for the default
// host provider
func FindDefaultAliasForHost(hostName string) (string, error) {
	pageConfig, _ := readConfigFile()

	for _, provider := range pageConfig.Providers {
		if provider.Type != "host" {
			continue
		}

		if provider.Name == hostName && provider.Default {
			return provider.Alias, nil
		}
	}

	return "", errors.New("no host found for " + hostName)
}

// FindRegistrarCredentials returns the credentials for a registrar
func FindRegistrarCredentials(alias string) (Credentials, error) {
	pageConfig, _ := readConfigFile()
	credentials := Credentials{}

	for _, provider := range pageConfig.Providers {
		if provider.Type != "registrar" {
			continue
		}

		if provider.Alias == alias {
			return provider.Credentials, nil
		}
	}

	return credentials, errors.New("")
}

// FindDefaultAliasForRegistrar returns the alias for the default
// registrar provider
func FindDefaultAliasForRegistrar(registrarName string) (string, error) {
	pageConfig, _ := readConfigFile()

	for _, provider := range pageConfig.Providers {
		if provider.Type != "registrar" {
			continue
		}

		if provider.Name == registrarName && provider.Default {
			return provider.Alias, nil
		}
	}

	return "", errors.New("no registrar found for " + registrarName)
}

func readConfigFile() (PageConfig, error) {
	var config PageConfig
	configData, err := os.ReadFile(ConfigPath)

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

func writeConfigFile(config PageConfig) error {
	file, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		log.Fatal("error parsing config file", err)
		return err
	}
	// TODO: perm os.FileMode? 0644?
	err = os.WriteFile(ConfigPath, []byte(file), 0644)

	if err != nil {
		log.Fatal("error writing config file", err)
		return err
	}

	return nil
}

// IsDefault changes the 'Default' field of ProviderConfig
// to false if there already exists a default provider
// for the host
func shouldChangeDefaultField(providers []ProviderConfig, provider *ProviderConfig) {
	for _, p := range providers {
		if p.Name == provider.Name && p.Default {
			// there already exists a default provider
			// for host
			provider.Default = false
		}
	}
}
