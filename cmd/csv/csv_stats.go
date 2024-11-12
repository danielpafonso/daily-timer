package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"daily-timer/internal"
)

type pastData struct {
	idx   int
	count int
	data  []int
}

// ReadStats Read statistics file and calculates stats
func ReadStats(team string, participants []string, limitDailies int) ([]internal.Stats, error) {
	statFile := fmt.Sprintf("stat-%s.csv", team)

	if _, ok := os.Stat(statFile); ok != nil || limitDailies == 0 {
		outputStats := make([]internal.Stats, 0)
		for _, name := range participants {
			outputStats = append(outputStats, internal.Stats{
				Name:   name,
				Active: true,
			})
		}
		return outputStats, nil
	}
	// initialite read and output maps
	outputMap := make(map[string]*internal.Stats)
	pastMap := make(map[string]*pastData)
	for _, name := range participants {
		outputMap[name] = &internal.Stats{
			Name:   name,
			Active: true,
		}
		pastMap[name] = &pastData{
			idx:   0,
			count: 0,
			data:  make([]int, limitDailies),
		}
	}
	// read file
	fileData, err := os.ReadFile(statFile)
	if err != nil {
		return nil, err
	}
	skipedHeader := false
	for _, line := range strings.Split(strings.TrimSpace(string(fileData)), "\n") {
		if !skipedHeader {
			skipedHeader = true
			continue
		}
		// date,name,value
		data := strings.Split(line, ",")
		value, _ := strconv.Atoi(data[2])
		past := pastMap[data[1]]

		past.data[past.idx] = value
		past.count += 1
		past.idx += 1
		if past.idx == limitDailies {
			past.idx = 0
		}
	}
	// calculate max and average
	for name, past := range pastMap {
		limit := min(past.count, limitDailies)
		sum := 0
		maxValue := 0
		for i := 0; i < limit; i++ {
			sum += past.data[i]
			maxValue = max(maxValue, past.data[i])
		}
		outputMap[name].Max = maxValue
		outputMap[name].Average = sum / limit
	}

	// create otuput array
	outputStats := make([]internal.Stats, 0)
	for _, v := range outputMap {
		outputStats = append(outputStats, *v)
	}
	return outputStats, nil
}

// WriteDaily writes current daily session to file
func WriteDaily(team string, stats []internal.Stats) error {
	var file *os.File
	var err error
	statFile := fmt.Sprintf("stat-%s.csv", team)
	// check if file don't exists, if so writes header
	if _, ok := os.Stat(statFile); ok != nil {
		file, err = os.Create(statFile)
		if err != nil {
			return err
		}
		file.WriteString("date,name,value\n")
	} else {
		file, err = os.OpenFile(statFile, os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
		}
	}

	now := time.Now().UTC()
	for _, stat := range stats {
		if stat.Active {
			file.WriteString(fmt.Sprintf(
				"%s,%s,%d\n",
				now.Format("2006-01-02 15:04:05.999"),
				stat.Name,
				stat.Current,
			))
		}
	}
	return nil
}
