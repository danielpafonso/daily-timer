package internal

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DailyTimes struct {
	Name string
	Date time.Time
	Time int
}

type Stats struct {
	Current int
	Average int
	Max     int
}

func Open(team string) (*sql.DB, error) {
	dbPath := fmt.Sprintf("stat-%s.sql", team)

	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func GetStats(dbConn *sql.DB, limitDailies int) (map[string]Stats, error) {
	query := fmt.Sprintf(`
SELECT
	distinct o.name,
	(select avg(value) from (select value from dailies as i where i.name=o.name order by time desc limit %d)),
	(select max(value) from (select value from dailies as i where i.name=o.name order by time desc limit %d))
from dailies as o;`, limitDailies, limitDailies)
	ctx := context.TODO()

	rows, err := dbConn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]Stats)

	for rows.Next() {
		var row Stats
		var name string
		err = rows.Scan(
			&name,
			&row.Average,
			&row.Max,
		)
		if err != nil {
			return nil, err
		}
		stats[name] = row
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
