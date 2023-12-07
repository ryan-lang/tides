package download

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ryan-lang/noaa-tidesandcurrents/client/dataApi"
	"github.com/ryan-lang/noaa-tidesandcurrents/client/metadataApi"
	"github.com/ryan-lang/tides/harmonics"
	"github.com/spf13/cobra"
)

var stationId string

var noaaStationCmd = &cobra.Command{
	Use:   "noaaStation",
	Short: "Download NOAA station data & save to local file",
	Long: `Download NOAA station harmonic constituents & datums from the NOAA CO-OPS API.

Example:
tides download noaaStation --station-id 9447130
	`,
	Run: func(cmd *cobra.Command, args []string) {

		harmonicsRes, err := downloadNOAAHarmonics(stationId)
		if err != nil {
			log.Printf("Error downloading NOAA harmonic constituent data: %s\n", err)
		}

		datumRes, err := downloadNOAADatums(stationId)
		if err != nil {
			log.Printf("Error downloading NOAA datum data: %s\n", err)
		}

		document := &harmonics.StationDocument{
			HarmonicConstituents: harmonicsRes,
			Datums:               datumRes,
		}

		json, err := json.Marshal(document)
		if err != nil {
			log.Printf("error marshalling station document: %s", err)
		}

		// write json to file
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", outputPath, stationId), json, 0644)
		if err != nil {
			log.Printf("error writing station document to file: %s", err)
		}
	},
}

func init() {
	DownloadCmd.AddCommand(noaaStationCmd)

	noaaStationCmd.PersistentFlags().StringVarP(&stationId, "station", "s", "", "NOAA station ID")
	noaaStationCmd.MarkPersistentFlagRequired("station")
}

func downloadNOAAHarmonics(stationId string) ([]*harmonics.HarmonicConstituent, error) {

	// do the remote request
	noaaClient := metadataApi.NewClient(true, "github.com/ryan-lang/tides")
	req := metadataApi.NewStationRequest(noaaClient, stationId)
	res, err := req.HarmonicConstituents(context.Background(), &metadataApi.HarmonicConstituentsRequest{Units: "metric"})
	if err != nil {
		return nil, fmt.Errorf("error getting harmonic constituents: %s", err)
	}

	// transmute into our struct format
	s := make([]*harmonics.HarmonicConstituent, len(res.HarmonicConstituents))
	for i, c := range res.HarmonicConstituents {
		s[i] = &harmonics.HarmonicConstituent{
			Name: c.Name,
			// Number:      c.Number,
			// Description: c.Description,
			Amplitude:  c.Amplitude,
			PhaseUTC:   c.PhaseGMT,
			PhaseLocal: c.PhaseLocal,
			Speed:      c.Speed,
		}
	}

	return s, nil
}

func downloadNOAADatums(stationId string) ([]*harmonics.Datum, error) {

	// do the remote request
	noaaClient := dataApi.NewClient(true, "github.com/ryan-lang/tides")
	res, err := noaaClient.Datums(context.Background(), &dataApi.DatumsRequest{StationID: stationId, Units: "metric"})
	if err != nil {
		return nil, fmt.Errorf("error getting station datums: %s", err)
	}

	// pack response into json
	s := make([]*harmonics.Datum, len(res.Datums))
	for i, c := range res.Datums {
		s[i] = &harmonics.Datum{
			Name:  c.Name,
			Value: c.Value,
		}
	}

	return s, nil
}
