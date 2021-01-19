package providers

import hosts "builtonpage.com/main/providers/hosts"

var supportedHosts = map[string]IHost{
	"page": hosts.PageHost{},
}

type IHost interface {
	Deploy() bool
}
type HostProvider struct {
	SupportedHosts map[string]IHost
}

var hostProvider HostProvider = HostProvider{
	SupportedHosts: map[string]IHost{
		"page": hosts.PageHost{},
	},
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
