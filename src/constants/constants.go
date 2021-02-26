package constants

type Consts struct {
	appName       string
	appVersion    string
	acmeServerURL string
	analyticsID   string
}

var consts Consts = Consts{
	appName:       "Page CLI",
	appVersion:    "v0.1.0-alpha.1",
	acmeServerURL: "https://acme-v02.api.letsencrypt.org/directory",
	analyticsID:   "UA-189047059-2",
}

func Constants() Consts {
	return consts
}
func AppName() string {
	return consts.appVersion
}
func AppVersion() string {
	return consts.appVersion
}

func AcmeServerURL() string {
	return consts.acmeServerURL
}

func AnalyticsID() string {
	return consts.analyticsID
}
