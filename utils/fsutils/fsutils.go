package fsutils

import (
	"os"
	"os/user"
	"path/filepath"
)

func Write(path string, data string) error {
    err := os.WriteFile(path, []byte(data), 0600)

    if err != nil {
        return err
    }

	return nil
}

func Read(path string) (string, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func GetDir(path string) string {
	return filepath.Dir(path)
}

func MakeDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func Exists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func Home() string {
	usr, err := user.Current()

	if err == nil {
		return usr.HomeDir
	}

	return ""
}