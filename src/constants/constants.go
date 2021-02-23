package constants

type Consts struct {
	appVersion    string
	acmeServerUrl string
}

var consts Consts = Consts{
	appVersion:    "v0.1.0-alpha.1",
	acmeServerUrl: "https://acme-v02.api.letsencrypt.org/directory",
}

func Constants() Consts {
	return consts
}

func AppVersion() string {
	return consts.appVersion
}

func AcmeServerUrl() string {
	return consts.acmeServerUrl
}
