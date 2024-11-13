package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"daily-timer/internal"
	"daily-timer/internal/sqlite"
)

const (
	Version string = "0.0.2"
)

func main() {
	var configPath string
	var showVersion bool

	flag.StringVar(&configPath, "c", "config.json", "Path to configuration file")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("daily-timer %s\n", Version)
		return
	}

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
	stats, err := GetStats(dbConn, configs.Participants, configs.Status.LastDailies)
	if err != nil {
		log.Panic(err)
	}

	// defering writing current session to DB
	defer InsertDaily(dbConn, stats)

	// Initialize ui
	appUI := internal.NewAppUI(*configs, &stats)
	// Start ui
	err = appUI.Start(Version)
	if err != nil {
		log.Panic(err)
	}
}