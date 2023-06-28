package file

import (
	"io"
	"os"
)

func CreateFile(fileContents string, filePathName string) error {
	file, err := os.Create(filePathName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, fileContents)
	if err != nil {
		return err
	}

	return file.Sync()
}
