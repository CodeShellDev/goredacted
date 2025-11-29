package cmd

import (
	service "github.com/codeshelldev/goredacted/internals/service"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
    Short: "Stop background task",
	Args: stop_validate,
    Long:  `Stop background service for goredacted`,
	Run: stop_run,
}

func init() {
	Register(stopCmd)
}

func stop_validate(cmd *cobra.Command, args []string) error {
	return nil
}

func stop_run(cmd *cobra.Command, args []string) {
	err := service.Despawn()

	if err != nil {
		cmd.PrintErrln("error stopping background service:", err.Error())
		return
	}
}