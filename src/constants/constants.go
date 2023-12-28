package constants

type Consts struct {
	appName       string
	appVersion    string
	appTier       string
	acmeServerURL string
	analyticsID   string
	upgradeInfo   string
	appAuthors    []string
}

var consts Consts = Consts{
	appName:       "page cli",
	appVersion:    "v0.1.0-alpha.13",
	acmeServerURL: "https://acme-v02.api.letsencrypt.org/directory",
	analyticsID:   "UA-189047059-2",
	appTier:       "free version",
	upgradeInfo:   "upgrade available https://buy.stripe.com/bIYbL12XA83P7pC3cc",
	appAuthors:    []string{"Roy Moran (roy.moran@icloud.com)"},
}

func Constants() Consts {
	return consts
}

func AppName() string {
	return consts.appName
}

func AppTier() string {
	return consts.appTier
}

func AppVersion() string {
	return consts.appVersion
}

func AppAuthors() []string {
	return consts.appAuthors
}

func AppUpgradeInfo() string {
	return consts.upgradeInfo
}

func AcmeServerURL() string {
	return consts.acmeServerURL
}

func AnalyticsID() string {
	return consts.analyticsID
}
