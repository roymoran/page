package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"pagecli.com/main/constants"
)

type LogRecord struct {
	Project string     `json:"project"`
	Name    string     `json:"name"`
	Level   string     `json:"level"`
	Message string     `json:"message"`
	Context AppContext `json:"context"`
}

type AppContext struct {
	Version     string            `json:"version"`
	OS          string            `json:"os"`
	Arch        string            `json:"arch"`
	NumCPU      int               `json:"num_cpu"`
	GoVersion   string            `json:"go_version"`
	Environment map[string]string `json:"environment"`
	Caller      string            `json:"caller"`
	StackTrace  []string          `json:"stack_trace"`
}

func SendLog(logRecord LogRecord) {
	/**
	 * Send a log record to the logging server.
	 *
	 * @param logRecord - A struct containing the log fields ('name', 'level', 'message').
	 * @return error - An error if the log record could not be sent.
	 * @return nil - If the log record was successfully sent.
	 * @example
	 * logRecord := LogRecord{
	 *     Level: "info",
	 *     Message: "Starting application",
	 * }
	 * err := SendLog(logRecord)
	 * if err != nil {
	 *     fmt.Println("Failed to send log: ", err)
	 * }
	 */

	var trace []string = CollectStackTrace()

	go func() {
		if !constants.IsProductionValue() {
			return
		}
		// set default project name
		logRecord.Project = "page"
		logRecord.Name = constants.AppName()
		logRecord.Context = collectContext()
		logRecord.Context.StackTrace = trace

		// set logging server URL, username, and password
		url := constants.LoggingServerURLValue()
		username := constants.LoggingServerUsername()
		password := constants.LoggingServerPasswordValue()

		// Convert log record to JSON
		jsonData, err := json.Marshal(logRecord)
		if err != nil {
			fmt.Printf("Failed to marshal log record: %v\n", err)
			return
		}

		// Create HTTP request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Failed to create request: %v\n", err)
			return
		}

		// Set headers
		req.SetBasicAuth(username, password)
		req.Header.Set("Content-Type", "application/json")

		// Create HTTP client
		client := &http.Client{}

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to send log request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Check response status (optional, since you don't care about the response)
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to send log: status code %d\n", resp.StatusCode)
		}
	}()
}

// CollectContext gathers debug context for the log.
func collectContext() AppContext {
	// Retrieve runtime caller information
	_, file, line, ok := runtime.Caller(2)
	caller := "unknown"
	if ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}

	// Filter environment variables to include only relevant ones
	environment := make(map[string]string)
	for _, key := range []string{"APP_TIER", "PRODUCTION", "GOOGLE_ANALYTICS_API_SECRET", "LOGGING_SERVER_URL"} {
		if value, exists := os.LookupEnv(key); exists {
			environment[key] = value
		}
	}

	return AppContext{
		Version:     constants.AppVersion(),
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		NumCPU:      runtime.NumCPU(),
		GoVersion:   runtime.Version(),
		Environment: environment,
		Caller:      caller,
	}
}

// CollectStackTrace collects a formatted stack trace for debugging.
func CollectStackTrace() []string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // Skip 3 frames (runtime.Callers, CollectStackTrace, and SendLog)
	frames := runtime.CallersFrames(pcs[:n])

	var stackTrace []string
	for {
		frame, more := frames.Next()
		stackTrace = append(stackTrace, fmt.Sprintf("%s\n\t%s:%d", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return stackTrace
}
