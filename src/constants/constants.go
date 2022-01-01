package constants

type Consts struct {
	appName       string
	appVersion    string
	acmeServerURL string
	analyticsID   string
	appAuthors    []string
}

var consts Consts = Consts{
	appName:       "page cli",
	appVersion:    "v0.1.0-alpha.10",
	acmeServerURL: "https://acme-v02.api.letsencrypt.org/directory",
	analyticsID:   "UA-189047059-2",
	appAuthors:    []string{"Roy Moran (https://github.com/roymoran)"},
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

func AppAuthors() []string {
	return consts.appAuthors
}

func AcmeServerURL() string {
	return consts.acmeServerURL
}

func AnalyticsID() string {
	return consts.analyticsID
}
