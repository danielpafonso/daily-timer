package plugins

import "daily-timer/internal"

// Functions that plugins must implement
type FileOperations interface {
	Connect(connectionString string) error
	GetStats(participants []string, limitDailies int) ([]internal.Stats, error)
	InsertDailies(stats *[]internal.Stats, writeTemp bool) error
}
