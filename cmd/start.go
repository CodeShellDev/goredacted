package cmd

import (
	"os"

	"github.com/codeshelldev/goredacted/internals/config"
	"github.com/codeshelldev/goredacted/internals/defaults"
	"github.com/codeshelldev/goredacted/internals/env"
	"github.com/codeshelldev/goredacted/internals/seed"
	service "github.com/codeshelldev/goredacted/internals/service"
	"github.com/codeshelldev/goredacted/utils/fsutils"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
    Short: "Start goredacted",
	Args: start_validate,
    Long:  `Start main loop for goredacted`,
	Run: start_run,
}

func init() {
	Register(startCmd)
	
	startCmd.Flags().StringVarP(&env.ENV.SEED_PATH, "seed", "s", defaults.DEFAULTS.SEED_PATH, "path of seed file")
	startCmd.Flags().StringVarP(&env.ENV.CONFIG_PATH, "config", "c", defaults.DEFAULTS.CONFIG_PATH, "path of config file")
	startCmd.Flags().StringVarP(&env.ENV.LISTEN_ADDR, "listen", "l", "", "http endpoint to communicate on")

	startCmd.Flags().BoolVarP(&env.ENV.BACKGROUND, "detach", "d", defaults.DEFAULTS.BACKGROUND, "detach and spawn background service")
}

func start_validate(cmd *cobra.Command, args []string) error {
	return nil
}

func start_run(cmd *cobra.Command, args []string) {
	if !fsutils.Exists(env.ENV.SEED_PATH) {
		cmd.PrintErrln(env.ENV.SEED_PATH, "not found")
		return
	}

	seed.LoadFile(env.ENV.SEED_PATH)
	config.Load(env.ENV.CONFIG_PATH)

	if os.Getenv("DAEMON_MODE") == "1" || !env.ENV.BACKGROUND {
		if env.ENV.LISTEN_ADDR != "" {
			service.Listen(env.ENV.LISTEN_ADDR)
		} else {
			service.Hook()
		}
		return
	}

	if env.ENV.BACKGROUND {
		err := service.Spawn(os.Args[1:]...)

		if err != nil {
			cmd.PrintErrln("error starting background process:", err.Error())
			return
		}
	}
}