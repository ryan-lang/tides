package harmonics

import (
	"fmt"
	"strings"
)

type (
	Datum struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}
)

func (h *Harmonics) GetDatum(name string) *Datum {
	for _, d := range h.Datums {
		if strings.EqualFold(d.Name, name) {
			return d
		}
	}

	return nil
}

func (h *Harmonics) DatumConvert(from, to string, val float64) (float64, error) {
	fromDatum := h.GetDatum(from)
	if fromDatum == nil {
		return 0, fmt.Errorf("datum not found: %s", from)
	}

	toDatum := h.GetDatum(to)
	if toDatum == nil {
		return 0, fmt.Errorf("datum not found: %s", to)
	}

	return val + fromDatum.Value - toDatum.Value, nil
}
