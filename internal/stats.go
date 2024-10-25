package internal

import (
	"database/sql"
	"time"

	"daily-timer/internal/sqlite"
)

type Stats struct {
	Current int
	Average int
	Max     int
	Active  bool
}

// AddNewPeople check if there is new persons not in database and add them to stats
func AddNewPeople(stats map[string]Stats, people []string) {
	for _, person := range people {
		if _, ok := stats[person]; !ok {
			stats[person] = Stats{
				Current: 0,
				Average: 0,
				Max:     0,
				Active:  true,
			}
		}
	}
}

// GetStats retrives calculated states from db
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
			Active:  true,
		}
	}

	return stats, nil
}

// InsertDaily writes current daily session to db
func InsertDaily(dbConn *sql.DB, stats map[string]Stats) error {
	now := time.Now()
	insertData := make([]sqlite.Dailies, 0)
	for k, v := range stats {
		if v.Active {
			insertData = append(insertData, sqlite.Dailies{
				Name: k,
				Date: now,
				Time: v.Current,
			})
		}
	}
	err := sqlite.InsertDaily(dbConn, insertData)
	if err != nil {
		return err
	}
	return nil
}
