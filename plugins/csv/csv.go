package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"daily-timer/internal"
)

type funcs struct {
	filePath string
}

type pastData struct {
	idx   int
	count int
	data  []int
}

func (f *funcs) Connect(connectString string) error {
	f.filePath = fmt.Sprintf("stat-%s.csv", connectString)
	return nil
}

func (f *funcs) GetStats(participants []string, limitDailies int) ([]internal.Stats, error) {
	// short-circuit when limitDailies is zero, or stat file does not exist
	if _, ok := os.Stat(f.filePath); ok != nil || limitDailies == 0 {
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
	fileData, err := os.ReadFile(f.filePath)
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
		past, ok := pastMap[data[1]]
		if !ok {
			//person in stats but not in config
			continue
		}
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
		if limit == 0 {
			continue
		}
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

func (f *funcs) InsertDailies(stats *[]internal.Stats, writeTemp bool) error {
	var file *os.File
	var err error
	// check if file don't exists, if so writes header
	if _, ok := os.Stat(f.filePath); ok != nil {
		file, err = os.Create(f.filePath)
		if err != nil {
			return err
		}
		file.WriteString("date,name,value\n")
	} else {
		file, err = os.OpenFile(f.filePath, os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			return err
		}
	}

	now := time.Now().UTC()
	for _, stat := range *stats {
		if stat.Temp {
			if writeTemp {
				file.WriteString(fmt.Sprintf(
					"%s,%s,%d\n",
					now.Format("2006-01-02 15:04:05.999"),
					stat.Name,
					stat.Current,
				))
			}
		} else if stat.Active {
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

// Export symbols
var FileOperations funcs
