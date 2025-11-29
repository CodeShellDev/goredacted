package clipboardutils

import (
	"github.com/atotto/clipboard"
)

func Read() string {
	content, err := clipboard.ReadAll()

	if err != nil {
		return ""
	}

	return content
}

func Write(content string) error {
	return clipboard.WriteAll(content)
}