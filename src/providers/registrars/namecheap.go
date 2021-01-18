package providers

type Namecheap struct {
}

func (n Namecheap) RegisterDomain() bool {
	return true
}

func (n Namecheap) ConfigureDns() bool {
	return true
}
