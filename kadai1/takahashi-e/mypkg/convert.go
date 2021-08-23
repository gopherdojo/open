// Package cmd will command and control you.
package mypkg

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
)

//Assert error Utility
func Assert(msg string, err error) error {
	if err != nil {
		fmt.Println(msg, "error")
		return err
	}
	return err
}

// Return filepath which is selected extension below selected directory.\
func (c *Arguments) ReturnFilePath() error {

	err := os.Chdir(c.SelectedDirectory)
	p, err := os.Getwd()
	err = filepath.Walk(p,
		func(paths string, info fs.FileInfo, err error) error {

			if filepath.Ext(paths) == c.SelectedFileType {
				c.StringPath = append(c.StringPath, paths)
			}
			return nil
		})
	return err
}

// Replace from selected file extension to othe file extension.
func (c *Arguments) ReplaceExt(filePath, from, to string) string {
	ext := filepath.Ext(filePath)
	if len(from) > 0 && ext != from {
		return filePath
	}
	return filePath[:len(filePath)-len(ext)] + to
}

// Arrange case for conversion.
func (c *Arguments) CaseConvert(from string, to string, fn string, img image.Image) error {
	fno := c.ReplaceExt(fn, from, "-out"+to)
	fo, err := os.Create(fno)
	defer fo.Close()

	err = Assert("OS.Create", err)

	switch to {
	case ".gif":
		return gif.Encode(fo, img, nil)

	case ".png":

		return png.Encode(fo, img)
	case ".jpg":
		return jpeg.Encode(fo, img, nil)
	}

	return err
}

// Convert image in  each case
func (c *Arguments) ConvertImage(fn string) error {
	f, err := os.Open(fn)
	err = Assert("OS.Open", err)
	defer f.Close()

	img, _, err := image.Decode(f)
	err = Assert("Decode", err)

	imageType := filepath.Ext(fn)

	switch imageType {

	case ".jpg":
		switch c.ConvertedFileType {
		case ".gif":
			err = c.CaseConvert(".jpg", ".gif", fn, img)
		default:
			err = c.CaseConvert(".jpg", ".png", fn, img)
		}

	case ".jpeg":
		switch c.ConvertedFileType {
		case ".gif":
			err = c.CaseConvert(".jpg", ".gif", fn, img)
		default:
			err = c.CaseConvert(".jpg", ".png", fn, img)
		}

	case ".png":
		switch c.ConvertedFileType {
		case ".jpeg":
			err = c.CaseConvert(".png", ".jpeg", fn, img)
		case ".jpg":
			err = c.CaseConvert(".png", ".jpg", fn, img)
		case ".gif":
			err = c.CaseConvert(".png", ".gif", fn, img)
		default:
			err = c.CaseConvert(".png", ".jpg", fn, img)
		}
	}
	return err
}

// Run in main.go
func (c *Arguments) Run() error {
	err := c.ReturnFilePath()
	err = Assert("ReturnFilePath", err)
	for _, v := range c.StringPath {
		err := c.ConvertImage(v)
		err = Assert("ConvertImage", err)
	}
	return nil
}
