package providers

type IRegistrar interface {
	RegisterDomain() bool
	ConfigureDns() bool
}
