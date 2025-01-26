package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"plugin"
	"strings"

	"daily-timer/internal"
	"daily-timer/internal/ui"
	"daily-timer/plugins"
)

const (
	Version string = "1.1.1-rc"
)

func main() {
	var configPath string
	var fileMode string
	var showVersion bool

	flag.StringVar(&configPath, "c", "config.json", "Path to configuration file.")
	flag.StringVar(&fileMode, "m", "sqlite", "Stat file interface, possible values: [sqlite, csv]. Defaults to sqlite.")
	flag.BoolVar(&showVersion, "v", false, "Show program's version number and exit.")
	flag.Parse()

	if showVersion {
		fmt.Printf("daily-timer %s\n", Version)
		return
	}

	var path string
	if fileMode == "sqlite" {
		path = "sqlite.so"
	} else if fileMode == "csv" {
		path = "csv.so"
	} else {
		log.Panicf("unexpected file interface, %s", fileMode)
		return
	}
	// load plugin for file operations
	pluginBinary, err := plugin.Open(path)
	if err != nil {
		log.Panic(err)
	}
	symbolPlugin, err := pluginBinary.Lookup("FileOperations")
	if err != nil {
		log.Panic(err)
	}
	fileOperations, ok := symbolPlugin.(plugins.FileOperations)
	if !ok {
		log.Panic("unexpected type from module symbol")
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
