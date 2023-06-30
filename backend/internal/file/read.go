package file

import (
	"embed"
	"os"
)

func ReadFile(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ReadEmbededFile(fs embed.FS, filename string) (string, error) {
	b, err := fs.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
