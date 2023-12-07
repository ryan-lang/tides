package download

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ryan-lang/noaa-tidesandcurrents/client/metadataApi"
	"github.com/ryan-lang/tides/harmonics"
	"github.com/spf13/cobra"
)

type (
	harmonicConstituentDocument struct {
		HarmonicConstituents []harmonics.HarmonicConstituent `json:"harmonic_constituents"`
	}
)

var stationId string

// noaaHarmonicsCmd represents the noaaHarmonics command
var noaaHarmonicsCmd = &cobra.Command{
	Use:   "noaaHarmonics",
	Short: "NOAA harmonic constituent data",
	Long: `Download NOAA harmonic constituent data from the NOAA CO-OPS API.

Example:
tides download noaaHarmonics --station-id 9447130
	`,
	Run: func(cmd *cobra.Command, args []string) {
		err := downloadNOAAHarmonics(stationId, fmt.Sprintf("%s/%s.json", outputPath, stationId))
		if err != nil {
			fmt.Printf("Error downloading NOAA harmonic constituent data: %s\n", err)
		}
	},
}

func init() {
	DownloadCmd.AddCommand(noaaHarmonicsCmd)

	noaaHarmonicsCmd.PersistentFlags().StringVarP(&stationId, "station", "s", "", "NOAA station ID")
	noaaHarmonicsCmd.MarkPersistentFlagRequired("station")
}

func downloadNOAAHarmonics(stationId, outputPath string) error {

	// do the remote request
	noaaClient := metadataApi.NewClient(true, "github.com/ryan-lang/tides")
	req := metadataApi.NewStationRequest(noaaClient, stationId)
	res, err := req.HarmonicConstituents(context.Background(), &metadataApi.HarmonicConstituentsRequest{Units: "metric"})
	if err != nil {
		return fmt.Errorf("error getting harmonic constituents: %s", err)
	}

	// pack response into json
	s := make([]harmonics.HarmonicConstituent, len(res.HarmonicConstituents))
	for i, c := range res.HarmonicConstituents {
		s[i] = harmonics.HarmonicConstituent{
			Name: c.Name,
			// Number:      c.Number,
			// Description: c.Description,
			Amplitude:  c.Amplitude,
			PhaseUTC:   c.PhaseGMT,
			PhaseLocal: c.PhaseLocal,
			Speed:      c.Speed,
		}
	}
	json, err := json.Marshal(harmonicConstituentDocument{HarmonicConstituents: s})
	if err != nil {
		return fmt.Errorf("error marshalling harmonic constituents: %s", err)
	}

	// write json to file
	err = ioutil.WriteFile(outputPath, json, 0644)
	if err != nil {
		return fmt.Errorf("error writing harmonic constituents to file: %s", err)
	}

	return nil
}
