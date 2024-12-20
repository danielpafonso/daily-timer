package main

import (
	"database/sql"
	"time"

	"daily-timer/internal"
	"daily-timer/internal/sqlite"
)

// GetStats retrives calculated states from db
func GetStats(dbConn *sql.DB, participants []string, limitDailies int) ([]internal.Stats, error) {
	sqlStats, err := sqlite.CalculateStats(dbConn, participants, limitDailies)
	if err != nil {
		return nil, err
	}

	// map persons on database
	persons := make(map[string]int)
	outputStats := make([]internal.Stats, 0)
	for _, x := range sqlStats {
		persons[x.Name] = 0
		outputStats = append(outputStats, internal.Stats{
			Name:    x.Name,
			Current: 0,
			Average: x.Average,
			Max:     x.Max,
			Active:  true,
		})
	}
	// check for new persons
	for _, name := range participants {
		if _, ok := persons[name]; !ok {
			outputStats = append(outputStats, internal.Stats{
				Name:    name,
				Average: 0,
				Max:     0,
				Active:  true,
			})
		}

	}
	return outputStats, nil
}

// InsertDaily writes current daily session to db
func InsertDaily(dbConn *sql.DB, stats *[]internal.Stats, writeTemp bool) error {
	now := time.Now()
	insertData := make([]sqlite.Dailies, 0)
	for _, stat := range *stats {
		if stat.Temp {
			if writeTemp {
				insertData = append(insertData, sqlite.Dailies{
					Name: stat.Name,
					Date: now,
					Time: stat.Current,
				})
			}
		} else if stat.Active {
			insertData = append(insertData, sqlite.Dailies{
				Name: stat.Name,
				Date: now,
				Time: stat.Current,
			})
		}
	}
	err := sqlite.InsertDaily(dbConn, insertData)
	if err != nil {
		return err
	}
	return nil
}
