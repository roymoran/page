package providers

import registrars "builtonpage.com/main/providers/registrars"

var supportedRegistrars = map[string]IRegistrar{
	"namecheap": registrars.Namecheap{},
	"page":      registrars.Page{},
}

type IRegistrar interface {
	RegisterDomain() bool
	ConfigureDns() bool
}
