package providers

type IProvider interface {
	Add() bool
	Remove() bool
	List() bool
}

var provider = map[string]IProvider{
	"host":      HostProvider{},
	"registrar": RegistrarProvider{},
}
