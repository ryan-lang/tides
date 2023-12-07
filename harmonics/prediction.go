package harmonics

import (
	"fmt"
	"math"
	"time"

	"github.com/ryan-lang/tides/astronomy"
)

type (
	Prediction struct {
		Timeline     Timeline
		Start        time.Time
		Constituents []*HarmonicConstituent
	}
	PredictionValue struct {
		Time  time.Time
		Hour  float64
		Level float64
	}
	PredictionExtremaValue struct {
		Time  time.Time
		Level float64
		Type  string
	}
	Timeline struct {
		Times        []time.Time
		ElapsedHours []float64
	}
	TideExtrema struct {
		Time  time.Time
		Level float64
		Type  string
	}
)

func (p *Prediction) GetTimelinePrediction() []PredictionValue {
	results := make([]PredictionValue, 0)
	harmonicResults, harmonicFactors := p.prepare()

	for i, time := range p.Timeline.Times {
		hour := p.Timeline.ElapsedHours[i]
		level := p.getLevel(hour, harmonicResults, harmonicFactors[i])

		results = append(results, PredictionValue{
			Time:  time,
			Hour:  hour,
			Level: level,
		})
	}

	return results
}

func (p *Prediction) GetHighLowPrediction() []PredictionExtremaValue {
	return p.getExtrema(0)
}

func (p *Prediction) GetLowsPrediction() []PredictionExtremaValue {
	return nil
}

func (p *Prediction) GetHighsPrediction() []PredictionExtremaValue {
	// results := make([]PredictionHighLowValue, 0)
	// for _, ex := range p.getExtrema(0) {
	// 	if ex.Type == "H" {
	// 		results = append(results, PredictionHighLowValue{
	// 			Time:  ex.Time,
	// 			Level: ex.Level,
	// 			Type:  "H",
	// 		})
	// 	}
	// }
	// return results

	return nil
}

func (p *Prediction) getLevel(hour float64, harmonicResults *harmonicResults, harmonicFactors *harmonicFactors) float64 {
	amplitudes := make([]float64, 0)
	result := 0.0

	for _, constituent := range p.Constituents {
		amplitude := constituent.Amplitude
		phase := constituent.PhaseUTC * astronomy.DEG_TO_RAD
		f := harmonicFactors.forms[constituent.Name]
		speed := harmonicResults.speeds[constituent.Name]
		u := harmonicFactors.nodes[constituent.Name]
		V0 := harmonicResults.values[constituent.Name]
		amplitudes = append(amplitudes, amplitude*f*math.Cos(speed*hour+(V0+u)-phase))
	}

	for _, item := range amplitudes {
		result += item
	}

	return result
}

// "partition" is partition hours; the number of hours for which we consider the node factors to be constant
func (p *Prediction) getExtrema(partition float64) []PredictionExtremaValue {
	extrema := make([]PredictionExtremaValue, 0)

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
	for start.Before(p.End()) {
		end := start.Add(time.Duration(partition) * time.Hour)

		// get the harmonic results at the start of the partition, and the
		// factor results for the first half of the partition
		// TODO: why first half??
		harmonicResults := harmonicResultsAtTime(p.Constituents, start)
		factorResults := harmonicFactorsForRange(p.Constituents, Timeline{Times: []time.Time{start.Add(time.Duration(partition*0.5) * time.Hour)}})

		// derivative functions d and d2 do not include time dependence of u or f
		// but they change slowly enough for that to be okay within our partition
		d := func(t float64) float64 {
			var sum float64
			for _, c := range p.Constituents {
				speed := harmonicResults.speeds[c.Name]
				amplitude := c.Amplitude
				phase := c.PhaseUTC * astronomy.DEG_TO_RAD // TODO: degtorad here??
				f := factorResults[0].forms[c.Name]
				V0 := harmonicResults.values[c.Name]
				u := factorResults[0].nodes[c.Name]
				sum -= speed * amplitude * f * math.Sin(speed*t+(V0+u)-phase)
			}
			return sum
		}
		d2 := func(t float64) float64 {
			var sum float64
			for _, c := range p.Constituents {
				speed := harmonicResults.speeds[c.Name]
				amplitude := c.Amplitude
				phase := c.PhaseUTC * astronomy.DEG_TO_RAD // TODO: degtorad here??
				f := factorResults[0].forms[c.Name]
				V0 := harmonicResults.values[c.Name]
				u := factorResults[0].nodes[c.Name]
				sum -= math.Pow(speed, 2.0) * amplitude * f * math.Cos(speed*t+(V0+u)-phase)
			}
			return sum
		}

		for i := 0; i < intervalCount; i++ {
			a := float64(i)*delta - offset
			b := float64(i+1)*delta - offset
			aTime := start.Add(time.Duration(a) * time.Hour)

			if aTime.After(p.End()) {
				break
			}

			if d(a)*d(b) < 0 {
				extremaHourOffset, err := newtonRaphson(d, d2, (a+b)/2, 1e-6, 100)
				if err != nil {
					fmt.Println("Error finding root:", err)
				}
				extremaTime := start.Add(time.Duration(extremaHourOffset * float64(time.Hour)))
				level := p.getLevel(extremaHourOffset, harmonicResults, factorResults[0])
				hilo := "L"
				if d2(extremaHourOffset) < 0 {
					hilo = "H"
				}
				if extremaTime.After(start) {
					extrema = append(extrema, PredictionExtremaValue{Time: extremaTime, Level: level, Type: hilo})
				}
			}
		}

		start = end
	}

	return extrema
}

func (p *Prediction) prepare() (*harmonicResults, []*harmonicFactors) {
	harmonicResults := harmonicResultsAtTime(p.Constituents, p.Start)
	harmonicFactors := harmonicFactorsForRange(p.Constituents, p.Timeline)
	return harmonicResults, harmonicFactors
}

func (p *Prediction) End() time.Time {
	return p.Timeline.Times[len(p.Timeline.Times)-1]
}

func MakeTimeline(start, end time.Time, interval time.Duration) Timeline {
	if interval == 0 {
		interval = 10 * time.Minute // default to 10 minutes if not provided
	}

	var timeline Timeline
	startTime := start.Unix()
	endTime := end.Unix()

	for lastTime := startTime; lastTime < endTime; lastTime += int64(interval.Seconds()) {
		timeline.Times = append(timeline.Times, time.Unix(lastTime, 0))
		elapsedHours := float64(lastTime-startTime) / 3600.0 // converting seconds to hours
		timeline.ElapsedHours = append(timeline.ElapsedHours, elapsedHours)
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

func harmonicFactorsForRange(constituents []*HarmonicConstituent, tl Timeline) []*harmonicFactors {
	factors := make([]*harmonicFactors, 0)
	for _, t := range tl.Times {
		factors = append(factors, harmonicFactorsAtTime(constituents, t))
	}
	return factors
}

func (p *Prediction) calculateMinDelta(t time.Time) float64 {
	minDelta := math.MaxFloat64
	astro := &astronomy.Astro{Time: t}

	for _, c := range p.Constituents {
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
