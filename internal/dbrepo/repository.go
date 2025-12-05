package dbrepo

import "daily-timer/internal"

type FileOperations interface {
	Connect(connectionString string) error
	GetStats(participants []string, limitDailies int) ([]internal.Stats, error)
	InsertDailies(stats *[]internal.Stats, writeTemp bool) error
}
