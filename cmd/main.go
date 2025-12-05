package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"daily-timer/internal"
	"daily-timer/internal/dbrepo"
	"daily-timer/internal/dbrepo/csv"
	"daily-timer/internal/dbrepo/sqlite"
	"daily-timer/internal/ui"
)

const (
	Version string = "1.2.0"
)

func main() {
	var configPath string
	var fileMode string
	var showVersion bool
	var fileOperations dbrepo.FileOperations

	flag.StringVar(&configPath, "c", "config.json", "Path to configuration file.")
	flag.StringVar(&fileMode, "m", "sqlite", "Stat file interface, possible values: [sqlite, csv]. Defaults to sqlite.")
	flag.BoolVar(&showVersion, "v", false, "Show program's version number and exit.")
	flag.Parse()

	if showVersion {
		fmt.Printf("daily-timer %s\n", Version)
		return
	}

	switch fileMode {
	case "sqlite":
		fileOperations = &sqlite.FileOperations{}
	case "csv":
		fileOperations = &csv.FileOperations{}
	default:
		log.Panicf("unexpected file interface, %s", fileMode)
	}

	// load configurations
	configs, err := internal.LoadConfigurations(configPath)
	if err != nil {
		log.Panic(err)
	}
	// Initialize "DB" and get Stats
	team := strings.TrimSuffix(filepath.Base(configPath), filepath.Ext(configPath))
	err = fileOperations.Connect(team)
	if err != nil {
		log.Panic(err)
	}
	stats, err := fileOperations.GetStats(configs.Participants, configs.Status.LastDailies)
	if err != nil {
		log.Panic(err)
	}

	// defering writing current session to DB
	defer fileOperations.InsertDailies(&stats, configs.AddTemp)

	// Initialize ui
	appUI := ui.NewAppUI(*configs, &stats)

	// Start ui
	err = appUI.Start(Version)
	if err != nil {
		log.Panic(err)
	}
}
