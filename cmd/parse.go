package cmd

import (
	"strings"

	"github.com/codeshelldev/goredacted/internals/config"
	"github.com/codeshelldev/goredacted/internals/defaults"
	"github.com/codeshelldev/goredacted/internals/env"
	"github.com/codeshelldev/goredacted/internals/seed"
	service "github.com/codeshelldev/goredacted/internals/service"
	"github.com/codeshelldev/goredacted/utils/fsutils"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
    Short: "Parse and redact stdin",
	Args: parse_validate,
    Long:  `Parse and redact using standard input`,
	Run: parse_run,
}

func init() {
	Register(parseCmd)

	parseCmd.Flags().StringVarP(&env.ENV.SEED_PATH, "seedfile", "f", defaults.DEFAULTS.SEED_PATH, "path of seed file")
	parseCmd.Flags().StringVarP(&env.ENV.SEED, "seed", "s", "", "seed")
	parseCmd.Flags().StringVarP(&env.ENV.CONFIG_PATH, "config", "c", defaults.DEFAULTS.CONFIG_PATH, "path of config file")

	parseCmd.MarkFlagsMutuallyExclusive("seed", "seedfile")
}

func parse_validate(cmd *cobra.Command, args []string) error {
	return nil
}

func parse_run(cmd *cobra.Command, args []string) {
	if strings.TrimSpace(env.ENV.SEED) != "" {
		seed.Load(env.ENV.SEED)
	} else {
		if !fsutils.Exists(env.ENV.SEED_PATH) {
			cmd.PrintErrln(env.ENV.SEED_PATH, "not found")
			return
		}

		err := seed.LoadFile(env.ENV.SEED_PATH)
		if err != nil {
			cmd.PrintErrln("seed file error:", err.Error())
			return
		}
	}

	err := config.Load(env.ENV.CONFIG_PATH)
	if err != nil {
		cmd.PrintErrln("config error:", err.Error())
		return
	}

	err = service.Redact()

	if err != nil {
		cmd.PrintErrln("error parsing stdin:", err.Error())
		return
	}
}