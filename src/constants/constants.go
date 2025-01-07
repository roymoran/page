package constants

import (
	"os"
)

type AppTierType string

const (
	Free AppTierType = "OSS"
	Demo AppTierType = "Demo"
	Paid AppTierType = "Paid"
)

type Consts struct {
	appName               string
	appVersion            string
	acmeServerURL         string
	loggingServerUsername string
	measurementID         string
	upgradeMessage        string
	appAuthors            []string
}

type AppVars struct {
	googleAnalyticsApiSecret string
	loggingServerPassword    string
	production               bool
	appTier                  string
	loggingServerURL         string
}

var consts Consts = Consts{
	appName:               "page",
	appVersion:            "0.3.3",
	acmeServerURL:         "https://acme-v02.api.letsencrypt.org/directory",
	loggingServerUsername: "loggingserver",
	measurementID:         "G-BRG7DX159G",
	upgradeMessage:        "This feature is not available in the current version of the tool. Purchase information available at https://pagecli.com#pricing.\n\nA single license is available at https://buy.stripe.com/bIYbL12XA83P7pC3cc\n",
	appAuthors:            []string{"Roy Moran (roy.moran@icloud.com)"},
}

// set at build time via ldflags with
// defaults assigned here
var GoogleAnalyticsApiSecret string = ""
var LoggingServerPassword string = ""
var LoggingServerURL string = "https://logs.nasapps.net/"
var IsProduction string = "true"
var AppTier string = "OSS"

func init() {
	GoogleAnalyticsApiSecret = getEnv("GOOGLE_ANALYTICS_API_SECRET", GoogleAnalyticsApiSecret)
	LoggingServerPassword = getEnv("BASIC_AUTH_PASSWORD", LoggingServerPassword)
	LoggingServerURL = getEnv("LOGGING_SERVER_URL", LoggingServerURL)
	IsProduction = getEnv("PRODUCTION", IsProduction)
	AppTier = getEnv("APP_TIER", AppTier)
}

func Constants() Consts {
	return consts
}

func AppVariables() AppVars {
	return AppVars{
		googleAnalyticsApiSecret: GoogleAnalyticsApiSecret,
		loggingServerPassword:    LoggingServerPassword,
		production:               IsProductionValue(),
		appTier:                  AppTier,
		loggingServerURL:         LoggingServerURL,
	}
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

func MeasurementID() string {
	return consts.measurementID
}

func LoggingServerURLValue() string {
	return LoggingServerURL
}

func LoggingServerUsername() string {
	return consts.loggingServerUsername
}

func GoogleAnalyticsApiSecretKey() string {
	return GoogleAnalyticsApiSecret
}

func LoggingServerPasswordValue() string {
	return LoggingServerPassword
}

func IsProductionValue() bool {
	return IsProduction == "true"
}

func AppTierValue() string {
	return string(AppTier)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
