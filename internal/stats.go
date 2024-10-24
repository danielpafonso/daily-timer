package internal

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"daily-timer/internal/sqlite"

	"golang.org/x/tools/go/analysis/passes/appends"
)

type Stats struct {
	Current int
	Average int
	Max     int
}

func GetStats(dbConn *sql.DB, limitDailies int) (map[string]Stats, error) {
	sqlStats, err := sqlite.CalculateStats(dbConn, limitDailies)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]Stats)
	for _, stat := range sqlStats {
		stats[stat.Name] = Stats{
			Current: 0,
			Average: stat.Average,
			Max:     stat.Max,
		}
	}

	return stats, nil
}

func InsertDaily(dbConn *sql.DB, stats map[string]Stats) error {
	now := time.Now()
	insertData := make([]sqlite.Dailies, 0)
	for k, v := range stats {
		insertData = append(insertData, sqlite.Dailies{
			Name: k,
			Date: now,
			Time: v.Current,
		})
	}
	err := sqlite.InsertDaily(dbConn, insertData)
	if err != nil {
		return err
	}
	return nil
}
