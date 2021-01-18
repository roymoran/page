package providers

type Page struct {
}

func (p Page) RegisterDomain() bool {
	return true
}

func (p Page) ConfigureDns() bool {
	return true
}
