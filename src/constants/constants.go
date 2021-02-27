package constants

type Consts struct {
	appName       string
	appVersion    string
	acmeServerURL string
	analyticsID   string
}

var consts Consts = Consts{
	appName:       "page cli",
	appVersion:    "v0.1.0-alpha.3",
	acmeServerURL: "https://acme-v02.api.letsencrypt.org/directory",
	analyticsID:   "UA-189047059-2",
}

func Constants() Consts {
	return consts
}
func AppName() string {
	return consts.appName
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
