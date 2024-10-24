package main

import (
	"flag"
	"log"

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
	// Other initializations

	// Initialize ui
	appUI := internal.NewAppUI(*configs)
	// Start ui
	err = appUI.Start()
	if err != nil {
		log.Panic(err)
	}
}
