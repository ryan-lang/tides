package harmonics

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ryan-lang/tides/constituents"
)

func LoadHarmonicConstituentsFromFile(filePath string) ([]*HarmonicConstituent, error) {

	// read the file
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening harmonic constituents file: %s", err)
	}

	// parse the json
	var doc struct {
		HarmonicConstituents []*HarmonicConstituent `json:"harmonic_constituents"`
	}
	err = json.NewDecoder(f).Decode(&doc)
	if err != nil {
		return nil, fmt.Errorf("error parsing harmonic constituents json: %s", err)
	}

	// associate each constituent with its model
	for _, c := range doc.HarmonicConstituents {
		c.Model = GetConstituentForName(c.Name)
	}

	return doc.HarmonicConstituents, nil
}

func GetConstituentForName(name string) HarmonicConstituentModel {
	switch name {
	case "Z0":
		return &constituents.CONSTITUENT_Z0
	case "SA":
		return &constituents.CONSTITUENT_SA
	case "SSA":
		return &constituents.CONSTITUENT_SSA
	case "MM":
		return &constituents.CONSTITUENT_MM
	case "MF":
		return &constituents.CONSTITUENT_MF
	case "Q1":
		return &constituents.CONSTITUENT_Q1
	case "O1":
		return &constituents.CONSTITUENT_O1
	case "K1":
		return &constituents.CONSTITUENT_K1
	case "J1":
		return &constituents.CONSTITUENT_J1
	case "M1":
		return &constituents.CONSTITUENT_M1
	case "P1":
		return &constituents.CONSTITUENT_P1
	case "S1":
		return &constituents.CONSTITUENT_S1
	case "OO1":
		return &constituents.CONSTITUENT_OO1
	case "2N2":
		return &constituents.CONSTITUENT_2N2
	case "N2":
		return &constituents.CONSTITUENT_N2
	case "NU2":
		return &constituents.CONSTITUENT_NU2
	case "M2":
		return &constituents.CONSTITUENT_M2
	case "LAM2":
		return &constituents.CONSTITUENT_LAM2
	case "L2":
		return &constituents.CONSTITUENT_L2
	case "T2":
		return &constituents.CONSTITUENT_T2
	case "S2":
		return &constituents.CONSTITUENT_S2
	case "R2":
		return &constituents.CONSTITUENT_R2
	case "K2":
		return &constituents.CONSTITUENT_K2
	case "M3":
		return &constituents.CONSTITUENT_M3
	case "MSF":
		return &constituents.CONSTITUENT_MSF
	case "2Q1":
		return &constituents.CONSTITUENT_2Q1
	case "RHO":
		return &constituents.CONSTITUENT_RHO
	case "MU2":
		return &constituents.CONSTITUENT_MU2
	case "2SM2":
		return &constituents.CONSTITUENT_2SM2
	case "2MK3":
		return &constituents.CONSTITUENT_2MK3
	case "MK3":
		return &constituents.CONSTITUENT_MK3
	case "MN4":
		return &constituents.CONSTITUENT_MN4
	case "M4":
		return &constituents.CONSTITUENT_M4
	case "MS4":
		return &constituents.CONSTITUENT_MS4
	case "S4":
		return &constituents.CONSTITUENT_S4
	case "M6":
		return &constituents.CONSTITUENT_M6
	case "S6":
		return &constituents.CONSTITUENT_S6
	case "M8":
		return &constituents.CONSTITUENT_M8
	default:
		log.Printf("! No constituent found for name: %s", name)
		return &constituents.CONSTITUENT_Z0
	}
}
