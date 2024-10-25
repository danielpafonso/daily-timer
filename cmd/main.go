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
	// Other initializations
	team := strings.TrimSuffix(filepath.Base(configPath), filepath.Ext(configPath))
	dbConn, err := sqlite.Open(team)
	if err != nil {
		log.Panic(err)
	}
	stats, err := internal.GetStats(dbConn, 2)
	if err != nil {
		log.Panic(err)
	}
	for k, v := range stats {
		log.Println(k, v)
	}

	defer internal.InsertDaily(dbConn, stats)

	// Initialize ui
	appUI := internal.NewAppUI(*configs)
	// Start ui
	err = appUI.Start()
	if err != nil {
		log.Panic(err)
	}
}
