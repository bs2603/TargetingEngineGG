package app

import (
	"log"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

type AccessLog struct {
	Timestamp  string      `json:"timestamp"`   // ISO8601 UTC
	StatusCode int         `json:"status"`      // HTTP status code
	Error      string      `json:"error"`       // Error
	Request    string      `json:"request"`     // app=x&country=y
	Response   interface{} `json:"response"`    // Final output object
	DurationMS int64       `json:"duration_ms"` // How long the handler took
}
