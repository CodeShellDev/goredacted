package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goredacted",
    Short: "Go Redact CLI",
    Long:  `Go CLI Programm to Redact and Mask sensitive information.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Register(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func init() {
	
}