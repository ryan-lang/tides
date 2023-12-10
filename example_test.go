package tides_test

import (
	"fmt"
	"time"

	"github.com/ryan-lang/tides"
)

func Example_loadHarmonicsAndPredict() {

	// Load harmonics from file
	har, err := tides.LoadHarmonicsFromFile("./data", "9447130")
	if err != nil {
		panic(err)
	}

	// Create a new prediction for a date range
	start := time.Date(2023, 4, 10, 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour * 1)
	prediction := har.NewRangePrediction(start, end, tides.WithInterval(time.Minute*10))

	// Get the prediction results
	results := prediction.Predict()
	for _, result := range results {
		fmt.Printf("%f @ %s\n", result.Level, result.Time)
	}

	// Output:
	// -0.617102 @ 2023-04-10 00:00:00 +0000 UTC
	// -0.480268 @ 2023-04-10 00:10:00 +0000 UTC
	// -0.344586 @ 2023-04-10 00:20:00 +0000 UTC
	// -0.210775 @ 2023-04-10 00:30:00 +0000 UTC
	// -0.079523 @ 2023-04-10 00:40:00 +0000 UTC
	// 0.048514 @ 2023-04-10 00:50:00 +0000 UTC
}
