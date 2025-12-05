package sqlite

import "time"

type Dailies struct {
	Name string
	Date time.Time
	Time int
}

type PastStats struct {
	Name    string
	Average int
	Max     int
}
