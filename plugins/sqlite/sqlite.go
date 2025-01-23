package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func checkTable(dbConn *sql.DB) error {
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name='dailies'"
	ctx := context.TODO()

	row := dbConn.QueryRowContext(ctx, query)
	var table string
	row.Scan(&table)
	if table == "" {
		// create table
		query = "CREATE TABLE dailies (name TEXT NOT NULL, time TEXT NOT NULL, value INTEGER NOT NULL);"
		ctx = context.TODO()
		dbConn.ExecContext(ctx, query)
	}
	return nil
}

// Open creates connection to DB
func Open(team string) (*sql.DB, error) {
	dbPath := fmt.Sprintf("stat-%s.sqlite", team)

	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	checkTable(dbConn)
	return dbConn, nil
}

// CalculateStats reads daily tables and returns average and max for participants
func CalculateStats(dbConn *sql.DB, persons []string, limit int) ([]PastStats, error) {
	query := fmt.Sprintf(`
SELECT
	DISTINCT o.name,
	(SELECT CAST(ROUND(AVG(value)) AS INTEGER) FROM (SELECT value FROM DAILIES AS i WHERE i.name=o.name ORDER BY time DESC LIMIT %d)),
	(SELECT MAX(value) FROM (SELECT value FROM dailies AS i WHERE i.name=o.name ORDER BY time DESC LIMIT %d))
FROM dailies AS o
WHERE o.name IN ("%s");`, limit, limit, strings.Join(persons, "\", \""))
	ctx := context.TODO()

	rows, err := dbConn.QueryContext(ctx, query)
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

// InsertDaily writes current daily session to DB
func InsertDaily(dbConn *sql.DB, daily []Dailies) error {
	queryArray := []string{}
	for _, elm := range daily {
		queryArray = append(
			queryArray,
			fmt.Sprintf("(\"%s\", \"%s\", %d)", elm.Name, elm.Date.UTC().Format("2006-01-02 15:04:05.999"), elm.Time),
		)
	}

	ctx := context.TODO()
	_, err := dbConn.ExecContext(ctx, fmt.Sprintf("INSERT INTO dailies (name, time, value) VALUES %s;", strings.Join(queryArray, ",\n")))
	if err != nil {
		return err
	}
	return nil
}
