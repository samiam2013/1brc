package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func main() {
	// read in the names from weather_stations.csv to a map
	fh, err := os.Open("weather_stations.csv")
	if err != nil {
		logrus.WithError(err).Fatal("could not open weather_stations.csv. is it in the current directory?")
	}
	defer func() { _ = fh.Close() }()

	wsNames := make(map[int]string, 42_000)
	scanner := bufio.NewScanner(fh)
	for i := 0; scanner.Scan(); i++ {
		name := scanner.Text()
		if strings.HasPrefix(name, "#") {
			continue
		}
		wsNames[i] = name
	}
	wsCount := len(wsNames)
	logrus.Infof("read %d weather station names", len(wsNames))

	// make a file 1brc.txt
	fh, err = os.Create("1brc.txt")
	if err != nil {
		logrus.WithError(err).Fatal("could not create 1brc.txt")
	}
	defer func() { _ = fh.Close() }()

	// randomly pick a weather station
	// 	generate a normally distributed temperature in Celcius
	//	between -99.9 and 99.9
	for lTW := 0; lTW < 1_000_000_000; lTW++ {
		line := wsNames[rand.Intn(wsCount)] + ";" + fmt.Sprintf("%.1f\n", rand.NormFloat64()*20+19)
		if _, err := fh.WriteString(line); err != nil {
			logrus.WithError(err).Fatal("could not write to 1brc.txt")
		}
		if lTW%10_00_000 == 0 {
			logrus.Infof("%.1f%% done", float64(lTW)/10_000_000)
		}
	}
}
