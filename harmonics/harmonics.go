package harmonics

import (
	"time"

	"github.com/ryan-lang/tides/astronomy"
)

const (
	DEFAULT_TIMELINE_INTERVAL = time.Minute
)

type (
	Harmonics struct {
		Constituents []*HarmonicConstituent
	}
	HarmonicConstituent struct {
		Name       string                   `json:"name"`
		Model      HarmonicConstituentModel `json:"-"`
		PhaseUTC   float64                  `json:"phase_UTC"`
		PhaseLocal float64                  `json:"phase_local"` // TODO how/hwere is this used
		Amplitude  float64                  `json:"amplitude"`
		Speed      float64                  `json:"speed"` // TODO how/hwere is this used
	}
	HarmonicConstituentModel interface {
		GetName() string
		Speed(*astronomy.Astro) float64
		Value(*astronomy.Astro) float64
		NodeFactor(*astronomy.Astro) float64
		FormFactor(*astronomy.Astro) float64
	}

	harmonicResults struct {
		values map[string]float64
		speeds map[string]float64
	}
	harmonicFactors struct {
		nodes map[string]float64
		forms map[string]float64
	}
)

func (h *Harmonics) NewPrediction(start, end time.Time) (*Prediction, error) {
	return &Prediction{
		Timeline:     MakeTimeline(start, end, DEFAULT_TIMELINE_INTERVAL),
		Constituents: h.Constituents,
		Start:        start,
	}, nil
}

func harmonicResultsAtTime(constituents []*HarmonicConstituent, t time.Time) *harmonicResults {

	// Create maps to store base values and speeds for each constituent.
	result := &harmonicResults{
		values: make(map[string]float64),
		speeds: make(map[string]float64),
	}

	// Initialize the starting astronomical conditions based on the prediction start time.
	astro := &astronomy.Astro{Time: t}

	// Iterate over each constituent to calculate and store their base value and speed at the start time.
	for _, constituent := range constituents {
		value := constituent.Model.Value(astro)
		speed := constituent.Model.Speed(astro)
		result.values[constituent.Name] = astronomy.DEG_TO_RAD * value
		result.speeds[constituent.Name] = astronomy.DEG_TO_RAD * speed
	}

	return result
}

func harmonicFactorsAtTime(constituents []*HarmonicConstituent, t time.Time) *harmonicFactors {

	factors := &harmonicFactors{
		nodes: make(map[string]float64),
		forms: make(map[string]float64),
	}

	stepAstro := &astronomy.Astro{Time: t}

	// Calculate node and form factors for each constituent at this time step.
	// Values are adjusted to ensure they fall within the [0, 360) range and converted to radians as needed.
	for _, constituent := range constituents {
		nodeFactor := modulus(constituent.Model.NodeFactor(stepAstro), 360)
		formFactor := modulus(constituent.Model.FormFactor(stepAstro), 360)

		factors.nodes[constituent.Name] = astronomy.DEG_TO_RAD * nodeFactor
		factors.forms[constituent.Name] = formFactor
	}

	return factors
}
