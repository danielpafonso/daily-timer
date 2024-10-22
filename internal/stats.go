package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"daily-timer/internal/sqlite"
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
	now := time.Now().Format(time.RFC3339)
	queryArray := []string{}
	for k, v := range stats {
		queryArray = append(queryArray, fmt.Sprintf("(\"%s\", \"%s\", %d)", k, now, v.Current))
	}
	// queryArray = append(queryArray, ";")

	ctx := context.TODO()
	_, err := dbConn.ExecContext(ctx, fmt.Sprintf("INSERT INTO dailies (name, time, value) VALUES %s;", strings.Join(queryArray, ",\n")))
	fmt.Println(strings.Join(queryArray, ",\n"))
	if err != nil {
		return err
	}
	return nil
}
