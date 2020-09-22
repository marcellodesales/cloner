package cmd

import (
	"github.com/marcellodesales/cloner/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows version details",
	Long:  `Shows the version information of this CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.ShowVersionDetails()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
