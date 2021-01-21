package providers

import "fmt"

type Namecheap struct {
}

func (n Namecheap) RegisterDomain() bool {
	return true
}

func (n Namecheap) ConfigureDns() bool {
	return true
}

func (n Namecheap) ConfigureRegistrar() bool {
	fmt.Println("configured namecheap registrar")
	return true
}
