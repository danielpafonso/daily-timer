package sqlite

import (
	"database/sql"
	"time"

	"daily-timer/internal"
)

type FileOperations struct {
	dbConn *sql.DB
}

func (f *FileOperations) Connect(team string) error {
	db, err := Open(team)
	if err != nil {
		return err
	}
	f.dbConn = db
	return nil
}

func (f *FileOperations) GetStats(participants []string, limitDailies int) ([]internal.Stats, error) {
	sqlStats, err := CalculateStats(f.dbConn, participants, limitDailies)
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

func (f *FileOperations) InsertDailies(stats *[]internal.Stats, writeTemp bool) error {
	now := time.Now()
	insertData := make([]Dailies, 0)
	for _, stat := range *stats {
		if stat.Temp {
			if writeTemp {
				insertData = append(insertData, Dailies{
					Name: stat.Name,
					Date: now,
					Time: stat.Current,
				})
			}
		} else if stat.Active {
			insertData = append(insertData, Dailies{
				Name: stat.Name,
				Date: now,
				Time: stat.Current,
			})
		}
	}
	err := InsertDaily(f.dbConn, insertData)
	if err != nil {
		return err
	}
	return nil
}
