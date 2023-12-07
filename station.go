package tides

import "github.com/ryan-lang/tides/harmonics"

type (
	Station struct {
		ID                   string                          `json:"id"`
		Name                 string                          `json:"name"`
		Continent            string                          `json:"continent"`
		Country              string                          `json:"country"`
		Region               string                          `json:"region"`
		Timezone             string                          `json:"timezone"`
		Latitude             float64                         `json:"latitude"`
		Longitude            float64                         `json:"longitude"`
		HarmonicConstituents []harmonics.HarmonicConstituent `json:"harmonic_constituents"`
	}
)
