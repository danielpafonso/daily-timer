package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"daily-timer/internal"
	"daily-timer/internal/sqlite"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "c", "config.json", "Path to configuration file")
	flag.Parse()

	configs, err := internal.LoadConfigurations(configPath)
	if err != nil {
		log.Panic(err)
	}
	// Initialize DB and get Stats
	team := strings.TrimSuffix(filepath.Base(configPath), filepath.Ext(configPath))
	dbConn, err := sqlite.Open(team)
	if err != nil {
		log.Panic(err)
	}
	stats, err := internal.GetStats(dbConn, configs.Participants, configs.Status.LastDailies)
	if err != nil {
		log.Panic(err)
	}
	internal.AddNewPeople(stats, configs.Participants)

	// Write current session to DB
	defer internal.InsertDaily(dbConn, stats)

	// Initialize ui
	appUI := internal.NewAppUI(*configs)
	// Start ui
	err = appUI.Start()
	if err != nil {
		log.Panic(err)
	}
}
