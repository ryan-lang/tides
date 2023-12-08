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
	PREDICTION_DATUM = "MTL" // TODO is this always true?
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
	}
	PredictionOpt   func(*Prediction)
	PredictionValue struct {
		Time  time.Time
		Level float64
		Type  string // I = intermediate, H = high, L = low
	}
	timeline struct {
		Times []time.Time
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

	// TODO: resize start of bracket to include prior extrema
	// TODO: resize end of bracket to include next extrema
	p.extendedStart = p.Start
	p.extendedEnd = p.End

	// step 1: calculate the tide results for our extended range; this should be wide enough
	// to include the prior and next extrema, but we haven't identified those points yet
	harmonicResults := harmonicResultsAtTime(p.Harmonics.Constituents, p.extendedStart)
	harmonicFactors := harmonicFactorsForRange(p.Harmonics.Constituents, p.extendedStart, p.extendedEnd, p.Interval)

	for i, time := range p.timeline().Times {
		hoursSinceStart := time.Sub(p.Start).Seconds() / 3600
		level := p.getLevel(hoursSinceStart, harmonicResults, harmonicFactors[i])

		p.extendedResults = append(p.extendedResults, &PredictionValue{
			Time:  time,
			Level: level,
		})
	}

	// step 2: make a second pass over the extended results to calculate
	// the extrema, and then trim our result set to only include the points between
	// the prior and next extrema. this is still a superset of the requested range.
	// after this point, the Type field of each PredictionValue will be set
	extrema := p.getExtrema(0)
	var priorExtrema, nextExtrema *PredictionValue
	for i, ex := range extrema {
		if priorExtrema == nil && ex.Time.After(p.Start) {
			priorExtrema = extrema[i-1] // TODO guard against index out of bounds
		}
		if nextExtrema == nil && ex.Time.After(p.End) {
			nextExtrema = ex
		}
	}
	p.extendedStart = priorExtrema.Time
	p.extendedEnd = nextExtrema.Time.Add(p.Interval) // end is exclusive, so add one interval
	p.extendedResults = filterPredictions(p.extendedResults, p.extendedStart, p.extendedEnd)

	// if this is a harmonic (reference) station, we are done
	if p.Harmonics.TidePredOffsets == nil {
		return filterPredictions(p.extendedResults, p.Start, p.End)
	}

	// for subordinate stations...

	// step 3: apply the offsets to the extrema
	for _, ex := range extrema {
		switch ex.Type {
		case "H":
			ex.Time = ex.Time.Add(time.Duration(p.Harmonics.TidePredOffsets.TimeOffsetHighTide) * time.Minute)
			ex.Level += p.Harmonics.TidePredOffsets.HeightOffsetHighTide
		case "L":
			ex.Time = ex.Time.Add(time.Duration(p.Harmonics.TidePredOffsets.TimeOffsetLowTide) * time.Minute)
			ex.Level += p.Harmonics.TidePredOffsets.HeightOffsetLowTide
		}
	}

	// step 4: apply interpolated offsets to the intermediate points
	// TODO

	return filterPredictions(p.extendedResults, p.Start, p.End)
}

// Calculates the extrema (highs & lows) using the parameters provided in the Prediction
func (p *Prediction) PredictExtrema() []*PredictionValue {
	return p.getExtrema(0)
}

// Same as PredictExtrema(), but only returns the lows
func (p *Prediction) PredictLows() []*PredictionValue {
	extrema := p.getExtrema(0)
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
	extrema := p.getExtrema(0)
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

// given a list of predictions, sets the Type field by calulating derivates to determine extrema
// also return the filtered list of extrema, to save a filtering step later on
func (p *Prediction) getExtrema(partition float64) []*PredictionValue {
	extrema := make([]*PredictionValue, 0)

	if partition == 0 {
		partition = 2400.0
	}

	// we assume that the extrema are separated by at least delta hours
	delta := p.calculateMinDelta(p.Start)

	// We search for stationary points from offset hours before t0 to
	// ensure we find any which might occur very soon after t0.
	offset := 24.0
	intervalCount := int(math.Ceil((partition+offset)/delta)) + 1

	start := p.Start
	for start.Before(p.End) {
		end := start.Add(time.Duration(partition) * time.Hour)

		// get the harmonic results at the start of the partition, and the
		// factor results at the midpoint of the partition, assuming this makes a
		// good approximation of the factor results for the entire partition
		harmonicResults := harmonicResultsAtTime(p.Harmonics.Constituents, start)
		factorResults := harmonicFactorsAtTime(p.Harmonics.Constituents, start.Add(time.Duration(partition*0.5)*time.Hour))

		// derivative functions d and d2 do not include time dependence of u or f
		// but they change slowly enough for that to be okay within our partition
		d := func(t float64) float64 {
			var sum float64
			for _, c := range p.Harmonics.Constituents {
				speed, amplitude, f, angle := calcConstituentParts(c, t, harmonicResults[c.Name], factorResults[c.Name])
				sum += speed * amplitude * f * -math.Sin(angle)
			}
			return sum
		}
		d2 := func(t float64) float64 {
			var sum float64
			for _, c := range p.Harmonics.Constituents {
				speed, amplitude, f, angle := calcConstituentParts(c, t, harmonicResults[c.Name], factorResults[c.Name])
				sum += math.Pow(speed, 2.0) * amplitude * f * -math.Cos(angle)
			}
			return sum
		}

		for i := 0; i < intervalCount; i++ {
			a := float64(i)*delta - offset
			b := float64(i+1)*delta - offset
			aTime := start.Add(time.Duration(a) * time.Hour)

			if aTime.After(p.End) {
				break
			}

			if d(a)*d(b) < 0 {
				extremaHourOffset, err := newtonRaphson(d, d2, (a+b)/2, 1e-6, 100)
				if err != nil {
					fmt.Println("Error finding root:", err)
				}
				extremaTime := start.Add(time.Duration(extremaHourOffset * float64(time.Hour)))
				level := p.getLevel(extremaHourOffset, harmonicResults, factorResults)
				hilo := "L"
				if d2(extremaHourOffset) < 0 {
					hilo = "H"
				}
				if extremaTime.After(start) {
					extrema = append(extrema, &PredictionValue{Time: extremaTime, Level: level, Type: hilo})
				}
			}
		}

		start = end
	}

	return extrema
}

func (p *Prediction) timeline() timeline {
	var timeline timeline

	start := p.Start
	end := p.End
	interval := p.Interval

	// handle single-point timeline
	if start.Equal(end) {
		timeline.Times = append(timeline.Times, start)
		return timeline
	}

	if interval == 0 {
		interval = 10 * time.Minute // default to 10 minutes if not provided
	}

	startTime := start.Unix()
	endTime := end.Unix()

	for lastTime := startTime; lastTime < endTime; lastTime += int64(interval.Seconds()) {
		timeline.Times = append(timeline.Times, time.Unix(lastTime, 0))
	}

	return timeline
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

// Newton-Raphson method for finding root
func newtonRaphson(f func(float64) float64, df func(float64) float64, initialGuess float64, tolerance float64, maxIterations int) (float64, error) {
	x := initialGuess
	for i := 0; i < maxIterations; i++ {
		fx := f(x)
		if math.Abs(fx) < tolerance {
			return x, nil
		}
		dfx := df(x)
		if dfx == 0 {
			return 0, fmt.Errorf("derivative is zero")
		}
		x = x - fx/dfx
	}
	return 0, fmt.Errorf("max iterations reached")
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
