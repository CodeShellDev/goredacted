package config

import (
	"encoding/json"

	"github.com/codeshelldev/goredacted/utils/fsutils"
)

type Config struct {
	INCLUDE []XXCLUDE `json:"include"`
	EXCLUDE []XXCLUDE `json:"exclude"`
}

type XXCLUDE struct {
	REGEX	string	`json:"regex"`
	PARSER	string	`json:"parser"`
}

var CONFIG = &Config{}

func Load(path string) error {
	content, err := fsutils.Read(path)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(content), &CONFIG)

	if err != nil {
		return err
	}
	
	return nil
}

func Create(path string) error {
	bytes, err := json.Marshal(CONFIG)

	if err != nil {
		return err
	}

	return fsutils.Write(path, string(bytes))
}