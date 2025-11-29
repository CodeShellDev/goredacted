package defaults

import (
	"path/filepath"

	"github.com/codeshelldev/goredacted/internals/env"
	"github.com/codeshelldev/goredacted/utils/fsutils"
)

var DEFAULTS = env.Env{
	SEED_PATH: filepath.Join(fsutils.Home(), ".config", "grdct", "tcdrg.grdct"),
	CONFIG_PATH: filepath.Join(fsutils.Home(), ".config", "grdct", "config.json"),
	BACKGROUND: false,
}