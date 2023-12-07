package tides

import (
	"fmt"
	"time"

	"github.com/ryan-lang/tides/constituents"
	"github.com/ryan-lang/tides/harmonics"
)

type (
	TidePrediction struct {
		Time   time.Time
		Height float64
	}
)

func GetTidePredictions(station *Station, start, end time.Time) (*harmonics.PredictionValue, error) {

	// make harmonics object
	predHarmonics := harmonics.Harmonics{
		Constituents: station.HarmonicConstituents,
	}

	// TODO
	offset := 0.0
	if offset != 0 {
		predHarmonics.Constituents = append(predHarmonics.Constituents, harmonics.HarmonicConstituent{
			Name:      "Z0",
			Model:     &constituents.CONSTITUENT_Z0,
			Phase:     0,
			Amplitude: offset,
		})
	}

	// set the time span
	predHarmonics.SetTimeSpan(start, end)

	prediction, err := predHarmonics.NewPrediction()
	if err != nil {
		return nil, fmt.Errorf("error creating tide prediction: %s", err)
	}

	return &prediction.GetTimelinePrediction()[0], nil
}
