package seed

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/codeshelldev/goredacted/utils/fsutils"
)

var SEED string

func Generate() (string, error) {
    key := make([]byte, 32)

    _, err := rand.Read(key)

    if err != nil {
        return "", err
    }

    return hex.EncodeToString(key), nil
}

func Create(path string) error {
	key, err := Generate()

	if err != nil {
		return err
	}

	fsutils.MakeDir(fsutils.GetDir(path))
	fsutils.Write(path, key)

	return nil
}

func LoadFile(path string) error {
	key, err := fsutils.Read(path)

	if err != nil {
		return err
	}

	Load(key)
	
	return nil
}

func Load(key string) {
	SEED = key
}