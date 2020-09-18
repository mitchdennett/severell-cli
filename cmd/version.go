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
	Short: "Print the version for Severell CLI",
	Long:  `Print the version for Severell CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Severell CLI version 0.0.1-SNAPSHOT")
	},
}