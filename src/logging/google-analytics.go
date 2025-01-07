package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/denisbrodbeck/machineid"
	"pagecli.com/main/constants"
)

type GAVars struct {
	GAEndpoint               string
	GADebugEndpoint          string
	MeasurementID            string
	GoogleAnalyticsApiSecret string
	DefaultEngagementTime    int64
	SessionExpirationMin     int64
}

var gaVars GAVars = GAVars{
	GAEndpoint:               "https://www.google-analytics.com/mp/collect",
	GADebugEndpoint:          "https://www.google-analytics.com/debug/mp/collect",
	MeasurementID:            constants.MeasurementID(),
	GoogleAnalyticsApiSecret: constants.GoogleAnalyticsApiSecretKey(),
	DefaultEngagementTime:    100,
	SessionExpirationMin:     30,
}

type SessionData struct {
	SessionID string `json:"session_id"`
	Timestamp int64  `json:"timestamp"`
}

type EventParams map[string]interface{}

type Event struct {
	Name   string      `json:"name"`
	Params EventParams `json:"params"`
}

type AnalyticsEventPayload struct {
	ClientID string  `json:"client_id"`
	Events   []Event `json:"events"`
}

type Analytics struct {
	Debug    bool
	mu       sync.Mutex
	clientID string
	session  SessionData
}

func NewAnalytics() *Analytics {
	return &Analytics{
		Debug: !constants.IsProductionValue(),
	}
}

func (a *Analytics) getOrCreateClientID() string {
	a.mu.Lock()
	defer a.mu.Unlock()

	id, err := machineid.ProtectedID(constants.AppName())
	if err != nil {
		id = "none"
	}

	a.clientID = id

	return a.clientID
}

func (a *Analytics) getOrCreateSessionID() string {
	a.mu.Lock()
	defer a.mu.Unlock()

	currentTime := time.Now().UnixMilli()

	if a.session.Timestamp != 0 {
		duration := (currentTime - a.session.Timestamp) / 60000
		if duration > gaVars.SessionExpirationMin {
			a.session = SessionData{}
		} else {
			a.session.Timestamp = currentTime
			return a.session.SessionID
		}
	}

	a.session = SessionData{
		SessionID: fmt.Sprintf("%d", currentTime),
		Timestamp: currentTime,
	}
	return a.session.SessionID
}

func (a *Analytics) FireEvent(name string, params EventParams) {
	/*
	 * This function is responsible for sending an event to Google Analytics.
	 * Example usage:
	 * analytics.FireEvent("testseven", logging.EventParams{
	 * "error_type": "error type",
	 *	"error_msg":  "error message",
	 *	"error_code": "error code",
	 * })
	 */
	if params == nil {
		params = EventParams{}
	}

	if _, exists := params["session_id"]; !exists {
		params["session_id"] = a.getOrCreateSessionID()
	}

	if _, exists := params["engagement_time_msec"]; !exists {
		params["engagement_time_msec"] = gaVars.DefaultEngagementTime
	}

	if _, exists := params["name"]; !exists {
		params["name"] = constants.AppName()
	}

	if _, exists := params["version"]; !exists {
		params["version"] = constants.AppVersion()
	}

	if _, exists := params["runtime"]; !exists {
		params["runtime_GOOS"] = runtime.GOOS
		params["runtime_version"] = runtime.Version()
		params["runtime_arch"] = runtime.GOARCH
		params["runtime_num_cpu"] = runtime.NumCPU()
		params["runtime_num_goroutine"] = runtime.NumGoroutine()
	}

	payload := AnalyticsEventPayload{
		ClientID: a.getOrCreateClientID(),
		Events: []Event{
			{
				Name:   name,
				Params: params,
			},
		},
	}

	payloadBytes, _ := json.Marshal(payload)

	endpoint := gaVars.GAEndpoint
	if a.Debug {
		endpoint = gaVars.GADebugEndpoint
	}

	resp, _ := http.Post(fmt.Sprintf("%s?measurement_id=%s&api_secret=%s", endpoint, gaVars.MeasurementID, gaVars.GoogleAnalyticsApiSecret),
		"application/json", bytes.NewBuffer(payloadBytes))

	defer resp.Body.Close()

	if a.Debug {
		log.Printf("Debug response: %v", resp)
	}
}
