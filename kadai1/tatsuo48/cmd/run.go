package cmd

import (
	"os"
	"path/filepath"
	"tatsuo48/pkg/converter"
)

func Run(dir, src, dest string) error {
	if err := traverseDir(dir, src, dest); err != nil {
		return err
	}
	return nil
}

func traverseDir(dir, src, dest string) error {
	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	converter, err := converter.NewConverter(src, dest)
	if err != nil {
		return err
	}

	for _, entry := range dirEntry {
		relativePath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			traverseDir(relativePath, src, dest)
		}
		if err := converter.Convert(relativePath); err != nil {
			continue
		}
	}
	return nil
}
