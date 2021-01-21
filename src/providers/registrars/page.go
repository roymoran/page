package providers

import "fmt"

type Page struct {
}

func (p Page) RegisterDomain() bool {
	return true
}

func (p Page) ConfigureDns() bool {
	return true
}

func (p Page) ConfigureRegistrar() bool {
	fmt.Println("configured page registrar")
	return true
}
