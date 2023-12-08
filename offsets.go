package tides

type (
	TidePredOffsets struct {
		RefStationID         string  `json:"ref_station_id"`
		HeightOffsetHighTide float64 `json:"height_offset_high_tide"`
		HeightOffsetLowTide  float64 `json:"height_offset_low_tide"`
		TimeOffsetHighTide   float64 `json:"time_offset_high_tide"` // in minutes
		TimeOffsetLowTide    float64 `json:"time_offset_low_tide"`  // in minutes
	}
)
