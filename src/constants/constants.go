package constants

type Consts struct {
	appName        string
	appVersion     string
	acmeServerURL  string
	analyticsID    string
	upgradeMessage string
	appAuthors     []string
}

var consts Consts = Consts{
	appName:        "page cli",
	appVersion:     "v0.3.1",
	acmeServerURL:  "https://acme-v02.api.letsencrypt.org/directory",
	analyticsID:    "UA-189047059-2",
	upgradeMessage: "This feature is not available in the current version of the tool. Purchase information available at https://pagecli.com#pricing.\n\nA single license is available at https://buy.stripe.com/bIYbL12XA83P7pC3cc\n",
	appAuthors:     []string{"Roy Moran (roy.moran@icloud.com)"},
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

func AppUpgradeMessage() string {
	return consts.upgradeMessage
}

func AcmeServerURL() string {
	return consts.acmeServerURL
}

func AnalyticsID() string {
	return consts.analyticsID
}
