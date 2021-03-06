package logging

import (
	"runtime"

	"builtonpage.com/main/constants"
	"github.com/denisbrodbeck/machineid"
	ga "github.com/jpillora/go-ogle-analytics"
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
	err = client.Send(ga.NewEvent(category, action).Label(label))
	if err != nil {
		return
	}
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
	err = client.Send(ga.NewException().Description(description).IsExceptionFatal(fatal))
	if err != nil {
		return
	}
}
