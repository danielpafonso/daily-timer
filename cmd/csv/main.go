package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"plugin"
	"strings"

	"daily-timer/internal"
	"daily-timer/plugins"
)

const (
	Version string = "1.0.1"
)

func main() {
	// test plugin
	path := "csv.so"
	// load plugin file
	plugin, err := plugin.Open(path)
	if err != nil {
		log.Panic(err)
	}
	// lookup symbol in plugin
	symbolPlugin, err := plugin.Lookup("FileOperations")
	if err != nil {
		log.Panic(err)
	}

	// reflect symbol to confirm correctness
	fileOperations, ok := symbolPlugin.(plugins.FileOperations)
	if !ok {
		log.Panic("unexpected type from module symbol")
	}

	// run function
	fileOperations.Load()

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
	stats, err := ReadStats(team, configs.Participants, configs.Status.LastDailies)
	if err != nil {
		log.Panic(err)
	}

	// defering writing current session to DB
	defer WriteDaily(team, &stats, configs.AddTemp)

	// Initialize ui
	appUI := internal.NewAppUI(*configs, &stats)
	// Start ui
	err = appUI.Start(Version)
	if err != nil {
		log.Panic(err)
	}
}
