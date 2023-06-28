package projectassistant

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type projectFile struct {
	name         string
	filepath     string
	isDirectory  bool
	size         int64
	lastModified time.Time
	metadata     map[string]any
}

type files []projectFile

func GetFileStructure(root string) (files, error) {
	return walk(root)
}

func (ff files) String() string {
	var result string
	for _, f := range ff {
		b, err := json.Marshal(f.metadata)
		if err != nil {
			panic(err)
		}
		metadata := string(b)
		result += f.name + ":" + f.filepath + ":" + strconv.FormatBool(f.isDirectory) + ":" + strconv.FormatInt(f.size, 10) + ":" + f.lastModified.String() + ":" + metadata + ";"
	}
	return result
}

func (ff files) MarkdownTable() string {
	table := "| Name | Filepath | Is Directory | Size | Last Modified | Metadata |\n"
	table += "| ---- | -------- | ------------ | ---- | ------------- | ------- |\n"

	for _, f := range ff {
		b, err := json.Marshal(f.metadata)
		if err != nil {
			panic(err)
		}
		metadata := string(b)
		table += "| " + f.name + " | " + f.filepath + " | " + strconv.FormatBool(f.isDirectory) + " | " + strconv.FormatInt(f.size, 10) + " | " + f.lastModified.String() + " | " + metadata + " |\n"
	}

	return table
}

func walk(root string) ([]projectFile, error) {
	files := []projectFile{}
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			f := projectFile{
				name:         info.Name(),
				filepath:     path,
				isDirectory:  info.IsDir(),
				size:         info.Size(),
				lastModified: info.ModTime(),
			}
			if !info.IsDir() {
				meta, err := getMetadata(path)
				if err != nil {
					return err
				}
				f.metadata = meta
			}
			files = append(files, f)
			return nil
		})
	return files, err
}

func getMetadata(path string) (map[string]any, error) {

	switch filepath.Ext(path) {
	case ".go":
		return goFileMetadata(path)
	default:
		return map[string]any{}, nil
	}
}

func goFileMetadata(filepath string) (map[string]any, error) {
	meta := map[string]any{}

	f, err := os.Open(filepath)
	if err != nil {
		return meta, fmt.Errorf("failed to open %q: %w", filepath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "package") {
			meta["package_name"] = scanner.Text()
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return meta, fmt.Errorf("scanner error for %q: %w", filepath, err)
	}

	return meta, nil
}
