package harmonics

import (
	"context"
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	noaaTides "github.com/ryan-lang/noaa-tidesandcurrents/client/dataApi"
	"github.com/stretchr/testify/assert"
)

const VAL_TOLERANCE = 0.01
const TIME_TOLERANCE = time.Minute * 1
const NOAA_VAL_TOLERANCE = 0.1               // TODO: why so poor?
const NOAA_TIME_TOLERANCE = time.Minute * 10 // TODO: why so poor?

func TestGetTimelinePrediction(t *testing.T) {
	har, err := LoadFromFile("../data/9447130.json")
	if err != nil {
		t.Error(err)
		return
	}

	start := time.Date(2023, 4, 10, 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	prediction := har.NewRangePrediction(start, end, WithInterval(time.Minute*10))

	results := prediction.Predict()
	expected := []float64{-0.61741382, -0.48061049, -0.34495767, -0.21117667, -0.07995447, 0.04805437}

	// Check length of results
	assert.Equal(t, 6, len(results))

	for i, result := range results {
		offBy := math.Abs(expected[i] - result.Level)
		fmt.Printf("result: %f at %s (off by %f)\n", result.Level, result.Time, offBy)
		assert.LessOrEqual(t, offBy, VAL_TOLERANCE, fmt.Sprintf("got %f, expected %f", result.Level, expected[i]))
	}
}

func TestGetHighLowPrediction(t *testing.T) {
	har, err := LoadFromFile("../data/9447130.json")
	if err != nil {
		t.Error(err)
		return
	}

	start := time.Date(2023, 4, 10, 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour * 24)

	prediction := har.NewRangePrediction(start, end)

	results := prediction.PredictExtrema()
	expectedLevel := []float64{1.272675070057166, -0.05582603086323429, 1.2097844518743732, -2.3734349778548514}
	expectedType := []string{"H", "L", "H", "L"}
	expectedTime := []time.Time{
		time.Date(2023, 4, 10, 3, 48, 42, 0, time.UTC),
		time.Date(2023, 4, 10, 9, 14, 44, 0, time.UTC),
		time.Date(2023, 4, 10, 14, 30, 3, 0, time.UTC),
		time.Date(2023, 4, 10, 21, 39, 7, 0, time.UTC),
	}

	// Check length of results
	assert.Equal(t, len(expectedLevel), len(results))

	for i, result := range results {
		valOffBy := math.Abs(expectedLevel[i] - result.Level)
		timeOffBy := math.Abs(expectedTime[i].Sub(result.Time).Minutes())
		fmt.Printf("result: %s of %f at %s (off by %f, %fm)\n", result.Type, result.Level, result.Time, valOffBy, timeOffBy)
		assert.LessOrEqual(t, valOffBy, VAL_TOLERANCE, fmt.Sprintf("got %f, expected %f", result.Level, expectedLevel[i]))
		assert.LessOrEqual(t, timeOffBy, TIME_TOLERANCE.Minutes(), fmt.Sprintf("got %s, expected %s", result.Time, expectedTime[i]))
		assert.Equal(t, expectedType[i], result.Type)
	}
}

func TestCompareWithNoaaHighLow(t *testing.T) {
	testStations := []string{"9447130", "9413450", "9411340"}

	ctx := context.Background()
	noaaClient := noaaTides.NewClient(true, "tides")

	now := time.Now().Add(time.Hour * 24 * -45)
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour * 24)

	for _, testStationID := range testStations {

		har, err := LoadFromFile("../data/" + testStationID + ".json")
		if err != nil {
			t.Error(err)
		}

		prediction := har.NewRangePrediction(start, end)

		localResults := prediction.PredictExtrema()
		remoteResults, err := noaaClient.TidePredictions(ctx, &noaaTides.TidePredictionsRequest{
			StationID: testStationID,
			Date: &noaaTides.DateParamBeginAndEnd{
				BeginDate: start,
				EndDate:   end,
			},
			Interval: noaaTides.INTERVAL_PARAM_HILO,
			Datum:    "MTL",
		})
		if err != nil {
			t.Error(err)
		}

		// Check length of results
		assert.Equal(t, len(remoteResults.Predictions), len(localResults))

		for i, result := range localResults {
			remoteTime := remoteResults.Predictions[i].Time
			remoteVal := remoteResults.Predictions[i].Value
			localTime := result.Time.In(time.UTC)
			localVal := result.Level
			fmt.Printf("noaa %f @ %s\t\t local %f @ %s\n", remoteVal, remoteTime, localVal, localTime)
			assert.LessOrEqual(t, math.Abs(remoteTime.Sub(localTime).Minutes()), NOAA_TIME_TOLERANCE.Minutes(), "minutes off")
			assert.LessOrEqual(t, math.Abs(remoteResults.Predictions[i].Value-result.Level), NOAA_VAL_TOLERANCE, fmt.Sprintf("index: %d", i))
		}
	}
}

func TestCompareWithNoaaTimeline(t *testing.T) {
	testStations := []string{"9447130", "9413450", "9411340"}

	ctx := context.Background()
	noaaClient := noaaTides.NewClient(true, "tides")

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour * 24)

	for _, testStationID := range testStations {

		har, err := LoadFromFile("../data/" + testStationID + ".json")
		if err != nil {
			t.Error(err)
		}

		prediction := har.NewRangePrediction(start, end)
		localResults := prediction.Predict()

		remotePredictions, err := noaaClient.TidePredictions(ctx, &noaaTides.TidePredictionsRequest{
			StationID: testStationID,
			Date: &noaaTides.DateParamBeginAndEnd{
				BeginDate: start,
				EndDate:   end,
			},
			Interval: noaaTides.INTERVAL_PARAM_1M,
			Datum:    "MTL",
		})
		if err != nil {
			t.Error(err)
		}

		// trim remote results
		remotePredictionsTrimmed := make([]noaaTides.TidePrediction, 0)
		for _, prediction := range remotePredictions.Predictions {
			if (prediction.Time.After(start) || prediction.Time.Equal(start)) && prediction.Time.Before(end) {
				remotePredictionsTrimmed = append(remotePredictionsTrimmed, prediction)
			}
		}

		// Check length of results
		assert.Equal(t, len(remotePredictionsTrimmed), len(localResults))

		// open file for writing
		f, err := os.Create("noaa-compare-" + testStationID + ".csv")
		if err != nil {
			t.Error(err)
		}
		f.WriteString("index,remote,local\n")

		// assert the results
		for i, result := range localResults {
			remoteTime := remotePredictionsTrimmed[i].Time
			remoteVal := remotePredictionsTrimmed[i].Value
			localTime := result.Time
			localVal := result.Level
			f.WriteString(fmt.Sprintf("%d", i) + ",")
			f.WriteString(fmt.Sprintf("%f", remoteVal) + ",")
			f.WriteString(fmt.Sprintf("%f", localVal))
			f.WriteString("\n")
			fmt.Printf("noaa %f @ %s\t\t local %f @ %s\n", remoteVal, remoteTime, localVal, localTime)
			assert.LessOrEqual(t, math.Abs(remotePredictionsTrimmed[i].Value-result.Level), NOAA_VAL_TOLERANCE, fmt.Sprintf("index: %d", i))
		}
	}
}
