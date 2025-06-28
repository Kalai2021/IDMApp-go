package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// FluentLogger represents a logger that sends logs to Fluentd
type FluentLogger struct {
	endpoint string
	enabled  bool
	client   *http.Client
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

var (
	fluentLogger *FluentLogger
)

// InitFluentLogger initializes the Fluentd logger
func InitFluentLogger() {
	endpoint := os.Getenv("FLUENT_ENDPOINT")
	enabled := os.Getenv("FLUENT_ENABLED") == "true"

	if endpoint == "" {
		endpoint = "http://localhost:24224"
	}

	fluentLogger = &FluentLogger{
		endpoint: endpoint,
		enabled:  enabled,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	if enabled {
		log.Printf("Fluentd logger initialized with endpoint: %s", endpoint)
	} else {
		log.Printf("Fluentd logger disabled")
	}
}

// GetFluentLogger returns the singleton Fluentd logger instance
func GetFluentLogger() *FluentLogger {
	if fluentLogger == nil {
		InitFluentLogger()
	}
	return fluentLogger
}

// sendLog sends a log entry to Fluentd
func (f *FluentLogger) sendLog(level, message string, data map[string]interface{}) {
	if !f.enabled {
		return
	}

	logEntry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Service:   "idmapp-backend",
		Message:   message,
		Data:      data,
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}

	req, err := http.NewRequest("POST", f.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		log.Printf("Error sending log to Fluentd: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Fluentd returned status: %d", resp.StatusCode)
	}
}

// Info logs an info level message
func (f *FluentLogger) Info(message string, data map[string]interface{}) {
	f.sendLog("info", message, data)
	log.Printf("[INFO] %s", message)
}

// Error logs an error level message
func (f *FluentLogger) Error(message string, data map[string]interface{}) {
	f.sendLog("error", message, data)
	log.Printf("[ERROR] %s", message)
}

// Warn logs a warning level message
func (f *FluentLogger) Warn(message string, data map[string]interface{}) {
	f.sendLog("warn", message, data)
	log.Printf("[WARN] %s", message)
}

// Debug logs a debug level message
func (f *FluentLogger) Debug(message string, data map[string]interface{}) {
	f.sendLog("debug", message, data)
	log.Printf("[DEBUG] %s", message)
}

// LogRequest logs HTTP request information
func (f *FluentLogger) LogRequest(method, path, remoteAddr, userAgent string, statusCode int, duration time.Duration, userID string) {
	data := map[string]interface{}{
		"method":      method,
		"path":        path,
		"remote_addr": remoteAddr,
		"user_agent":  userAgent,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
		"user_id":     userID,
	}

	f.Info(fmt.Sprintf("%s %s - %d (%dms)", method, path, statusCode, duration.Milliseconds()), data)
}

// LogDatabase logs database operation information
func (f *FluentLogger) LogDatabase(operation, table string, duration time.Duration, rowsAffected int64, err error) {
	data := map[string]interface{}{
		"operation":     operation,
		"table":         table,
		"duration_ms":   duration.Milliseconds(),
		"rows_affected": rowsAffected,
	}

	if err != nil {
		data["error"] = err.Error()
		f.Error(fmt.Sprintf("Database %s on %s failed: %v", operation, table, err), data)
	} else {
		f.Info(fmt.Sprintf("Database %s on %s completed (%dms, %d rows)", operation, table, duration.Milliseconds(), rowsAffected), data)
	}
}

// LogAuth logs authentication events
func (f *FluentLogger) LogAuth(event, userID, sessionID, ipAddress string, success bool, details map[string]interface{}) {
	data := map[string]interface{}{
		"event":      event,
		"user_id":    userID,
		"session_id": sessionID,
		"ip_address": ipAddress,
		"success":    success,
	}

	// Merge additional details
	for k, v := range details {
		data[k] = v
	}

	if success {
		f.Info(fmt.Sprintf("Auth %s successful for user %s", event, userID), data)
	} else {
		f.Error(fmt.Sprintf("Auth %s failed for user %s", event, userID), data)
	}
}
