package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Open(team string) (*sql.DB, error) {
	dbPath := fmt.Sprintf("stat-%s.sql", team)

	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func CalculateStats(dbConn *sql.DB, limit int) ([]PastStats, error) {
	// 	query := fmt.Sprintf(`
	// SELECT
	// 	distinct o.name,
	// 	(select avg(value) from (select value from dailies as i where i.name=o.name order by time desc limit %d)),
	// 	(select max(value) from (select value from dailies as i where i.name=o.name order by time desc limit %d))
	// from dailies as o;`, limit, limit)
	// 	ctx := context.TODO()
	//
	// 	rows, err := dbConn.QueryContext(ctx, query)
	query := `
SELECT
	distinct o.name,
	(select avg(value) from (select value from dailies as i where i.name=o.name order by time desc limit $1)),
	(select max(value) from (select value from dailies as i where i.name=o.name order by time desc limit $1))
from dailies as o;`
	ctx := context.TODO()

	rows, err := dbConn.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []PastStats
	for rows.Next() {
		var row PastStats
		err = rows.Scan(
			&row.Name,
			&row.Average,
			&row.Max,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}
	return data, nil
}

func InsertDaily(dbConn *sql.DB, daily []Dailies) error {
	now := time.Now().Format(time.RFC3339)
	queryArray := []string{}
	for _, elm := range daily {
		queryArray = append(queryArray, fmt.Sprintf("(\"%s\", \"%s\", %d)", elm.Name, now, elm.Time))
	}

	ctx := context.TODO()
	_, err := dbConn.ExecContext(ctx, fmt.Sprintf("INSERT INTO dailies (name, time, value) VALUES %s;", strings.Join(queryArray, ",\n")))
	if err != nil {
		return err
	}
	return nil
}
