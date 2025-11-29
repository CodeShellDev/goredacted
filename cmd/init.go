package cmd

import (
	"github.com/codeshelldev/goredacted/internals/config"
	"github.com/codeshelldev/goredacted/internals/defaults"
	"github.com/codeshelldev/goredacted/internals/env"
	"github.com/codeshelldev/goredacted/internals/seed"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
    Short: "Init goredacted",
	Args: init_validate,
    Long:  `Initiliaze goredacted for first run`,
	Run: init_run,
}

func init() {
	Register(initCmd)

	initCmd.Flags().StringVarP(&env.ENV.SEED_PATH, "seed", "s", defaults.DEFAULTS.SEED_PATH, "path of seed file")
	initCmd.Flags().StringVarP(&env.ENV.CONFIG_PATH, "config", "c", defaults.DEFAULTS.CONFIG_PATH, "path of config file")
}

func init_validate(cmd *cobra.Command, args []string) error {
	return nil
}

func init_run(cmd *cobra.Command, args []string) {
	err := seed.Create(env.ENV.SEED_PATH)
	if err != nil {
		cmd.PrintErrln("error writing to", env.ENV.SEED_PATH, ":", err.Error())
		return
	}

	err = config.Create(env.ENV.CONFIG_PATH)
	if err != nil {
		cmd.PrintErrln("error writing to", env.ENV.CONFIG_PATH, ":", err.Error())
		return
	}
}