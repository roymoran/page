package logging

import (
	"runtime"

	"github.com/denisbrodbeck/machineid"
	ga "github.com/jpillora/go-ogle-analytics"
	"pagecli.com/main/constants"
)

var client, _ = ga.NewClient(constants.AnalyticsID())

// LogEvent records a google analytics event
func LogEvent(category string, action string, label string, value int) {
	id, err := machineid.ProtectedID(constants.AppName())
	if err != nil {
		id = "none"
	}
	client.ClientID(id)
	client.UserAgentOverride("")
	client.ApplicationName(constants.AppName())
	client.ApplicationVersion(constants.AppVersion())
	client.ApplicationInstallerID(runtime.GOOS)
	client.Send(ga.NewEvent(category, action).Label(label))
}

// LogException records a google analytics exception/crash
func LogException(description string, fatal bool) {
	id, err := machineid.ProtectedID(constants.AppName())
	if err != nil {
		id = "none"
	}
	client.ClientID(id)
	client.UserAgentOverride("")
	client.ApplicationName(constants.AppName())
	client.ApplicationVersion(constants.AppVersion())
	client.ApplicationInstallerID(runtime.GOOS)
	client.Send(ga.NewException().Description(description).IsExceptionFatal(fatal))
}
