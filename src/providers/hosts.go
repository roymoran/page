package providers

import hosts "builtonpage.com/main/providers/hosts"

var supportedHosts = map[string]IHost{
	"page": hosts.PageHost{},
}

type IHost interface {
	Deploy() bool
}
