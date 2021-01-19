package providers

import "fmt"

type IRegistrar interface {
	RegisterDomain() bool
	ConfigureDns() bool
}

func (rp RegistrarProvider) Add() bool {
	fmt.Println("registrar add")
	return true
}

func (rp RegistrarProvider) List() bool {
	fmt.Println("registrar list")
	return true
}
