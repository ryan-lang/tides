package predict

import (
	"fmt"
	"log"
	"time"

	"github.com/araddon/dateparse"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/ryan-lang/tides/harmonics"
	"github.com/spf13/cobra"
)

var dateUntil, dateSince, dateFrom, dateTo string
var dataDir, stationId, units, datum, intervalStr string
var printUnits, printTimes, extrema bool

var PredictCmd = &cobra.Command{
	Use:   "predict",
	Short: "download tide data from remote sources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// load harmonics data
		har, err := harmonics.LoadFromFile(fmt.Sprintf("%s/%s.json", dataDir, stationId))
		if err != nil {
			log.Fatalf("error loading station data: %s", err)
		}

		// setup dates
		startDate := time.Now()
		endDate := startDate
		interval := time.Minute

		// setup the parser
		w := when.New(nil)
		w.Add(en.All...)
		w.Add(common.All...)

		// parse the dates
		if dateSince != "" {
			startDate = whenParseFatal(w, dateSince)
		} else if dateFrom != "" {
			startDate = dateParseFatal(dateFrom)
		}
		if dateUntil != "" {
			endDate = whenParseFatal(w, dateUntil)
		} else if dateTo != "" {
			endDate = dateParseFatal(dateTo)
		}

		// parse the interval
		if intervalStr != "" {
			interval, err = time.ParseDuration(intervalStr)
			if err != nil {
				log.Fatalf("Failed to parse interval: %v", err)
			}
		}

		// extrema requires a range
		if extrema && endDate.Equal(startDate) {
			endDate = startDate.Add(time.Hour * 24)
		}

		// create a prediction
		prediction := har.NewRangePrediction(startDate, endDate,
			harmonics.WithDatum(datum),
			harmonics.WithUnits(units),
			harmonics.WithInterval(interval),
		)

		if extrema {

			// get prediction
			results := prediction.PredictExtrema()

			// print results
			for _, result := range results {
				if printTimes {
					fmt.Printf("%s\t", result.Time.Format(time.RFC3339))
				}
				fmt.Printf("%s\t", result.Type)
				fmt.Printf("%f", result.Level)
				if printUnits {
					fmt.Printf("%s", prediction.Units)
				}
				fmt.Println()
			}
		} else {

			// get prediction
			results := prediction.Predict()

			// print results
			for _, result := range results {
				if printTimes {
					fmt.Printf("%s\t", result.Time.Format(time.RFC3339))
				}
				fmt.Printf("%f", result.Level)
				if printUnits {
					fmt.Printf("%s", prediction.Units)
				}
				fmt.Println()
			}

		}

	},
}

func init() {
	PredictCmd.PersistentFlags().StringVarP(&stationId, "station", "s", "", "station identifier (e.g. NOAA station ID); must match json file in data directory")
	PredictCmd.PersistentFlags().StringVarP(&dataDir, "data-dir", "d", "./data", "data directory containing station data")
	PredictCmd.PersistentFlags().StringVarP(&datum, "datum", "m", "mllw", "datum to use for prediction (mllw, mhhw, mhw, msl, mslw, msw, naw, stnd)")
	PredictCmd.PersistentFlags().StringVarP(&units, "units", "u", "m", "m (metric) or ft (imperial)")
	PredictCmd.PersistentFlags().StringVarP(&intervalStr, "interval", "i", "1m", "interval between predictions (e.g. 1h, 30m, 15m)")
	PredictCmd.PersistentFlags().BoolVarP(&printUnits, "print-units", "", false, "print units in output")
	PredictCmd.PersistentFlags().BoolVarP(&printTimes, "print-times", "", false, "print times in output")
	PredictCmd.PersistentFlags().BoolVarP(&extrema, "extrema", "e", false, "returns tide extrema (highs and lows) only")
	PredictCmd.PersistentFlags().StringVarP(&dateSince, "since", "", "", "relative start date for prediction (eg. yesterday, last friday)")
	PredictCmd.PersistentFlags().StringVarP(&dateUntil, "until", "", "", "relative end date for prediction (eg. tomorrow, next friday)")
	PredictCmd.PersistentFlags().StringVarP(&dateFrom, "from", "", "", "absolute start date for prediction (eg. 2019-01-01T00:00:00Z)")
	PredictCmd.PersistentFlags().StringVarP(&dateTo, "to", "", "", "absolute end date for prediction (eg. 2019-01-01T00:00:00Z)")
	PredictCmd.MarkPersistentFlagRequired("station")
}

func whenParseFatal(w *when.Parser, s string) time.Time {
	parsed, err := w.Parse(s, time.Now())
	if err != nil {
		log.Fatalf("Failed to parse date: %v", err)
	} else if parsed == nil {
		log.Fatalf("Failed to parse date: %s", dateUntil)
	}
	return parsed.Time
}

func dateParseFatal(s string) time.Time {
	d, err := dateparse.ParseLocal(s)
	if err != nil {
		log.Fatalf("Failed to parse date: %v", err)
	}
	return d
}
