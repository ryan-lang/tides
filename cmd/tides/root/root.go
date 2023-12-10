package root

import (
	"os"

	"github.com/ryan-lang/tides/cmd/tides/root/download"
	"github.com/ryan-lang/tides/cmd/tides/root/predict"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tides",
	Short: "A tide prediction calculator",
	Long:  `A tide prediction calculator that uses harmonic constituent data and astronomy math to provide tide predictions in go.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(download.DownloadCmd)
	rootCmd.AddCommand(predict.PredictCmd)
}
