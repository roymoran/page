package constants

type Consts struct {
	appVersion string
}

var consts Consts = Consts{
	appVersion: "v0.1.0",
}

func Constants() Consts {
	return consts
}

func AppVersion() string {
	return consts.appVersion
}
