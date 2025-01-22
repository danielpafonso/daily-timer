package main

import (
	"fmt"
	// "os"
	// "strconv"
	// "strings"
	// "time"

	"daily-timer/internal"
)

type funcs struct {
	filePath string
}

type pastData struct {
	idx   int
	count int
	data  []int
}

func (f *funcs) Connect(connectString string) error {
	f.filePath = fmt.Sprintf("stat-%s.csv", connectString)

	return nil
}

func (f *funcs) GetStats(participants []string, limitDailies int) ([]internal.Stats, error) {
	return nil, nil
}

func (f *funcs) InsertDailies(stats *[]internal.Stats, writeTemp bool) error {
	return nil
}

func (f *funcs) Load() {
	fmt.Println(">>  Execute function from plugin")
}

// Export symbols
var FileOperations funcs
