package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cacli",
	Short: "Cacli help you to build clean code in golang project",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute to execute cobra cmd command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
