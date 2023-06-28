package file

import (
	"embed"
	"fmt"
	"os"
)

func ReadFile(filename string) (string, error) {
	fmt.Println(os.Getwd())
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
