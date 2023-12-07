package download

import (
	"fmt"

	"github.com/spf13/cobra"
)

var outputPath string

var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download tide data from remote sources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("download called")
	},
}

func init() {
	DownloadCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", "./data", "output directory")
}
