package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type StationInfo struct {
	Max        float64
	Min        float64
	SumTemps   float64
	CountTemps int
}

var stationStats = make(map[string]StationInfo, 42_000)

func main() {
	// exepct arg 1 to be the file path
	if len(os.Args) < 2 {
		logrus.Fatal("Please provide a file path")
	}
	filePath := os.Args[1]

	start := time.Now()

	fh, err := os.Open(filePath)
	if err != nil {
		logrus.WithError(err).Fatal("Error opening file")
	}
	defer func() { _ = fh.Close() }()

	scnr := bufio.NewScanner(fh)
	for scnr.Scan() {
		line := scnr.Bytes()
		parts := bytes.Split(line, []byte(";"))
		name := string(parts[0])
		temp, _ := strconv.ParseFloat(string(parts[1]), 64)

		if station, ok := stationStats[name]; ok {
			station.CountTemps++
			station.SumTemps += temp
			if temp > station.Max {
				station.Max = temp
			}
			if temp < station.Min {
				station.Min = temp
			}
			stationStats[name] = station
			continue
		}
		stationStats[name] = StationInfo{
			Max:        temp,
			Min:        temp,
			SumTemps:   temp,
			CountTemps: 1,
		}
	}

	const padLen = -40
	for name, stats := range stationStats {
		fmt.Printf("%*s min %6.1f\t max %6.1f\t avg %6.1f\n",
			padLen, name, stats.Min, stats.Max, stats.SumTemps/float64(stats.CountTemps))
	}

	logrus.Infof("Execution time: %v", time.Since(start))
}
