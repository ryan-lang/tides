// Package tides provides a library for calculating tide predictions using provided harmonic data,
// and optionally datums and offset data.
//
// Supports both harmonic tide stations (reference stations), and subordinate tide stations
// (derived from reference stations).
//
// A CLI is also provided for downloading NOAA station data, and calculating tide predictions.
package tides

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/ryan-lang/tides/astronomy"
)

const (
	METERS_TO_FEET   = 3.28084
	PREDICTION_DATUM = "MTL"
)

type (
	Prediction struct {
		Start           time.Time
		End             time.Time
		Interval        time.Duration
		Harmonics       *Harmonics
		Datum           string
		Units           string
		extendedStart   time.Time
		extendedEnd     time.Time
		extendedResults []*PredictionValue // holds an expanded result set for working on
		extremaResults  []*PredictionValue // holds the extrema results to save a filtering step
	}
	PredictionOpt   func(*Prediction)
	PredictionValue struct {
		Time        time.Time
		Level       float64
		Type        string // I = intermediate, H = high, L = low
		lastExtrema *PredictionValue
		nextExtrema *PredictionValue
		// used to store uncorrected time/level prior to offsets being applied
		uncTime  time.Time
		uncLevel float64
	}
)

// Sets the datum on the Prediction
func WithDatum(datum string) PredictionOpt {
	return func(p *Prediction) {
		p.Datum = datum
	}
}

// Sets the units on the Prediction
func WithUnits(units string) PredictionOpt {
	return func(p *Prediction) {
		p.Units = units
	}
}

// Sets the interval on the Prediction
func WithInterval(interval time.Duration) PredictionOpt {
	return func(p *Prediction) {
		p.Interval = interval
	}
}

// Calculates a prediction using the parameters provided in the Prediction
func (p *Prediction) Predict() []*PredictionValue {

	// resize start & end of bracket so that prior & next extrema are included
	// we are liberal here, because we will trim the results later
	p.extendedStart = p.Start.Add(-24 * time.Hour)
	p.extendedEnd = p.End.Add(24 * time.Hour)

	// step 1: calculate the tide results for our extended range; this should be wide enough
	// to include the prior and next extrema, but we haven't identified those points yet
	harmonicResults := harmonicResultsAtTime(p.Harmonics.Constituents, p.extendedStart)
	harmonicFactors := harmonicFactorsForRange(p.Harmonics.Constituents, p.extendedStart, p.extendedEnd, p.Interval)

	var i int
	for t := p.extendedStart; t.Before(p.extendedEnd); t = t.Add(p.Interval) {
		elapsedHours := t.Sub(p.extendedStart).Hours()
		level := p.getLevel(elapsedHours, harmonicResults, harmonicFactors[i])

		p.extendedResults = append(p.extendedResults, &PredictionValue{
			Time:  t,
			Level: level,
		})
		i++
	}

	// step 2: make a second pass over the extended results to pick
	// the extrema, and then trim our result set to only include the points between
	// the prior and next extrema. this is still a superset of the requested range.
	// after this point, the Type field of each PredictionValue will be set
	p.extremaResults = p.getExtrema(p.extendedResults, harmonicResults, harmonicFactors)
	var priorExtrema, nextExtrema *PredictionValue
	for i, ex := range p.extremaResults {
		if priorExtrema == nil && ex.Time.After(p.Start) {
			if i == 0 {
				// this should not happen; it means no extrema was found in prior 24h
				fmt.Printf("something went wrong: no prior extrema found for %s. aborting...\n", p.Start)
				return nil
			}
			priorExtrema = p.extremaResults[i-1]
		}
		if nextExtrema == nil && ex.Time.After(p.End) {
			nextExtrema = ex
		}
	}
	p.extendedResults = filterPredictions(p.extendedResults, priorExtrema.Time, nextExtrema.Time)

	// if this is a harmonic (reference) station, we are done
	if p.Harmonics.TidePredOffsets == nil {
		return filterPredictions(p.extendedResults, p.Start, p.End)
	}

	// for subordinate stations...

	// step 3: apply the offsets to the extrema
	for _, ex := range p.extremaResults {
		ex.uncTime = ex.Time
		ex.uncLevel = ex.Level

		switch ex.Type {
		case "H":
			ex.Time = ex.Time.Add(time.Duration(p.Harmonics.TidePredOffsets.TimeOffsetHighTide) * time.Minute)
			ex.Level *= p.Harmonics.TidePredOffsets.HeightOffsetHighTide
		case "L":
			ex.Time = ex.Time.Add(time.Duration(p.Harmonics.TidePredOffsets.TimeOffsetLowTide) * time.Minute)
			ex.Level *= p.Harmonics.TidePredOffsets.HeightOffsetLowTide
		}
	}

	// step 4: apply interpolated offsets to the intermediate points
	for _, result := range p.extendedResults {
		if result.Type == "H" || result.Type == "L" {
			continue
		}
		if result.nextExtrema == nil || result.lastExtrema == nil {
			continue
		}

		// calculate the proportion of time & value that our point represents between the uncorrected extrema
		uncInterpTime := (result.Time.Sub(result.lastExtrema.uncTime).Minutes() / result.nextExtrema.uncTime.Sub(result.lastExtrema.uncTime).Minutes())
		uncInterpLevel := (result.Level - result.lastExtrema.uncLevel) / (result.nextExtrema.uncLevel - result.lastExtrema.uncLevel)
		//fmt.Printf("uncorrected interpolation: %f/%f\n", uncInterpTime, uncInterpLevel)

		// uncInterpTime represents the proportion of time that our point represents between the uncorrected extrema
		// so we can take this proportion and apply it to the corrected extrema to get the interpolated time
		offsetTime := time.Duration(uncInterpTime*result.nextExtrema.Time.Sub(result.lastExtrema.Time).Minutes()) * time.Minute
		interpolatedTime := result.lastExtrema.Time.Add(offsetTime)
		//fmt.Printf("interpolated time: %s (%s/%s) %s\n", interpolatedTime, time.Duration(uncInterpTime*result.nextExtrema.uncTime.Sub(result.lastExtrema.uncTime).Minutes())*time.Minute, result.nextExtrema.uncTime.Sub(result.lastExtrema.uncTime), offsetTime)

		// uncInterpLevel represents the proportion of value that our point represents between the uncorrected extrema
		// so we can take this proportion and apply it to the corrected extrema to get the interpolated level
		offsetLevel := uncInterpLevel * (result.nextExtrema.Level - result.lastExtrema.Level)
		interpolatedLevel := result.lastExtrema.Level + offsetLevel
		//fmt.Printf("interpolated level: %f (%f/%f) %f\n", interpolatedLevel, uncInterpLevel, result.nextExtrema.Level-result.lastExtrema.Level, offsetLevel)

		result.uncLevel = result.Level
		result.uncTime = result.Time
		result.Level = interpolatedLevel
		result.Time = interpolatedTime
		break
	}

	return filterPredictions(p.extendedResults, p.Start, p.End)
}

// Calculates the extrema (highs & lows) using the parameters provided in the Prediction
func (p *Prediction) PredictExtrema() []*PredictionValue {
	p.Predict()
	return filterPredictions(p.extremaResults, p.Start, p.End)
}

// Same as PredictExtrema(), but only returns the lows
func (p *Prediction) PredictLows() []*PredictionValue {
	extrema := p.PredictExtrema()
	results := make([]*PredictionValue, 0)
	for _, ex := range extrema {
		if ex.Type == "L" {
			results = append(results, ex)
		}
	}
	return results
}

// Same as PredictExtrema(), but only returns the highs
func (p *Prediction) PredictHighs() []*PredictionValue {
	extrema := p.PredictExtrema()
	results := make([]*PredictionValue, 0)
	for _, ex := range extrema {
		if ex.Type == "H" {
			results = append(results, ex)
		}
	}
	return results
}

func (p *Prediction) getLevel(t float64, harmonicResults harmonicResults, harmonicFactors harmonicFactors) float64 {
	amplitudes := make([]float64, 0)
	result := 0.0

	for _, constituent := range p.Harmonics.Constituents {
		_, amplitude, f, angle := calcConstituentParts(constituent, t, harmonicResults[constituent.Name], harmonicFactors[constituent.Name])
		amplitudes = append(amplitudes, amplitude*f*math.Cos(angle))
	}

	for _, item := range amplitudes {
		result += item
	}

	if p.Datum != "" && !strings.EqualFold(p.Datum, PREDICTION_DATUM) {
		datum, err := p.Harmonics.DatumConvert(PREDICTION_DATUM, p.Datum, result)
		if err != nil {
			log.Fatalf("Error converting datum: %s", err.Error())
		}
		result = datum
	}

	if p.Units == "ft" {
		result = result * METERS_TO_FEET
	}

	return result
}

func (p *Prediction) getExtrema(predictions []*PredictionValue, hResults harmonicResults, hFactors []harmonicFactors) (extrema []*PredictionValue) {
	var isFalling bool
	var lastExtrema *PredictionValue

	// can't work on less than 2 points
	if len(predictions) < 2 {
		return
	}

	// set initial value for isFalling, or else
	// the first point will be always considered an extrema
	if predictions[1].Level < predictions[0].Level {
		isFalling = true
	}

	for i := 0; i < len(predictions)-1; i++ {
		p := predictions[i]
		nextP := predictions[i+1]

		if nextP.Level < p.Level {
			if !isFalling {
				p.Type = "H"
				extrema = append(extrema, p)
				lastExtrema = p
			}
			isFalling = true
		} else {
			if isFalling {
				p.Type = "L"
				extrema = append(extrema, p)
				lastExtrema = p
			}
			isFalling = false
		}

		p.lastExtrema = lastExtrema
	}
	predictions[len(predictions)-1].lastExtrema = lastExtrema

	// we now have all the extrema; walk the list
	// again and set nextExtrema on each point
	var i, extremaCursor int
	for i < len(predictions) && extremaCursor < len(extrema) {
		nextExtrema := extrema[extremaCursor]
		p := predictions[i]

		// if we've passed the next extrema, advance to the next extrema
		// without advancing the prediction and continue
		if !p.Time.Before(nextExtrema.Time) {
			extremaCursor++
			continue
		}

		p.nextExtrema = nextExtrema
		i++
	}

	return
}

func modulus(a, b float64) float64 {
	result := math.Mod(a, b)
	if result < 0 {
		result += b
	}
	return result
}

func harmonicFactorsForRange(constituents []*HarmonicConstituent, start, end time.Time, interval time.Duration) []harmonicFactors {
	factors := make([]harmonicFactors, 0)
	for t := start; t.Before(end); t = t.Add(interval) {
		factors = append(factors, harmonicFactorsAtTime(constituents, t))
	}
	return factors
}

func (p *Prediction) calculateMinDelta(t time.Time) float64 {
	minDelta := math.MaxFloat64
	astro := &astronomy.Astro{Time: t}

	for _, c := range p.Harmonics.Constituents {
		speed := c.Model.Speed(astro)
		if speed != 0 {
			delta := 90.0 / speed
			if delta < minDelta {
				minDelta = delta
			}
		}
	}

	if minDelta == math.MaxFloat64 {
		minDelta = 0 // or some other default value or error handling
	}

	return minDelta
}

func filterPredictions(predictions []*PredictionValue, start, end time.Time) []*PredictionValue {
	filtered := make([]*PredictionValue, 0)
	for _, prediction := range predictions {
		if (prediction.Time.After(start) || prediction.Time.Equal(start)) && prediction.Time.Before(end) {
			filtered = append(filtered, prediction)
		}
	}
	return filtered
}

func calcConstituentParts(c *HarmonicConstituent, t float64, hResult harmonicResult, hFactor harmonicFactor) (speed, amplitude, f, angle float64) {
	amplitude = c.Amplitude
	phase := c.PhaseUTC * astronomy.DEG_TO_RAD
	f = hFactor.form
	speed = hResult.speed
	u := hFactor.node
	V0 := hResult.value
	angle = speed*t + (V0 + u) - phase
	return speed, amplitude, f, angle
}
