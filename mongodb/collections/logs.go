package collections

import (
	"time"
)

type LogEntry struct {
	Timestamp time.Time              `bson:"timestamp"`          // The time of the log event
	Message   string                 `bson:"message"`            // The log message
	Level     string                 `bson:"level"`              // Log level (INFO, ERROR, etc.)
	Metadata  map[string]interface{} `bson:"metadata,omitempty"` // Metadata about the log event
}
