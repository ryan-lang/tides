package main

import (
	"flag"
	"log"
	"time"

	"github.com/ryan-lang/tides"
)

func main() {

	// cli flags
	stationIdFlag := flag.String("station-id", "united-states-wa-seattle", "as defined by the `id` field in the station json")
	dataDirFlag := flag.String("data", "./data", "directory of the station json files")
	predictionStartFlag := flag.String("start", time.Now().Format(time.RFC822), "start time of the prediction")
	predictionEndFlag := flag.String("end", time.Now().Format(time.RFC822), "end time of the prediction")
	flag.Parse()

	// load station data
	stationStore := tides.LoadStationData(*dataDirFlag)

	// lookup the provided station
	station := stationStore.GetStation(*stationIdFlag)
	if station == nil {
		log.Fatalf("Tide station not found: %s\nMake sure an associated json file exists in the data directory.\n", *stationIdFlag)
	}

	// parse the prediction start and end times
	predictionStart, err := time.Parse(time.RFC822, *predictionStartFlag)
	if err != nil {
		log.Fatalf("Error parsing prediction start time: %s\n", err)
	}
	predictionEnd, err := time.Parse(time.RFC822, *predictionEndFlag)
	if err != nil {
		log.Fatalf("Error parsing prediction end time: %s\n", err)
	}

	// get the tide prediction
	predictions := station.GetTidePredictions(predictionStart, predictionEnd)
	for _, prediction := range predictions {
		log.Printf("%s: %f\n", prediction.Time, prediction.Height)
	}
}
