package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCmd root of cmd command
var RootCmd = &cobra.Command{
	Use:   "cacli",
	Short: "Cacli help you to build clean code in golang project",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute to execute cobra cmd command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
