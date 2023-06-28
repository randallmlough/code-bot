package projectassistant

import (
	"os"

	"github.com/randallmlough/code-bot/internal/file"
)

func createFile(body string, filepath string) error {
	return file.CreateFile(body, filepath)
}

func updateFile(body string, filepath string) error {
	input, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	input = append(input, []byte("\n")...)
	input = append(input, []byte(body)...)

	return createFile(string(input), filepath)
}
