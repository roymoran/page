package providers

type IHost interface {
	Deploy() bool
}

func (hp HostProvider) Add() bool {
	return true
}

func (hp HostProvider) Remove() bool {
	return true
}

func (hp HostProvider) List() bool {
	return true
}
