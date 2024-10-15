package internal

import (
	"encoding/json"
	"os"
)

type Configurations struct {
	Time         int      `json:"time"`
	Warning      int      `json:"warning"`
	Participants []string `json:"participants"`
	Random       bool     `json:"randomOrder"`
	Stopwatch    bool     `json:"stopwatch"`
	Status       struct {
		Display     bool `json:"display"`
		LastDailies int  `json:"lastDailies"`
	} `json:"stats"`
}

func LoadConfigurations(filepath string) (*Configurations, error) {
	var config Configurations
	// read file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	// unmashall data
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
