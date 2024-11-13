package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"daily-timer/internal"
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
	stats, err := ReadStats(team, configs.Participants, configs.Status.LastDailies)
	if err != nil {
		log.Panic(err)
	}

	// defering writing current session to DB
	defer WriteDaily(team, stats)

	// Initialize ui
	appUI := internal.NewAppUI(*configs, &stats)
	// Start ui
	err = appUI.Start()
	if err != nil {
		log.Panic(err)
	}
}
