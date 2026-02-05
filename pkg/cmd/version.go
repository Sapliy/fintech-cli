package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Sapliy CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Sapliy CLI v%s\n", rootCmd.Version)
	},
}
