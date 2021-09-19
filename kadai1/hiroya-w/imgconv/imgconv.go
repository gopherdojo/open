package imgconv

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

func Converter(directory, inputType, outputType string) error {
	imgPaths, err := getFiles(directory, inputType)
	if err != nil {
		return err
	}

	for _, path := range imgPaths {
		if err := convert(path, outputType); err != nil {
			return err
		}
	}
	return nil
}

func getFiles(directory, inputType string) ([]string, error) {
	var imgPaths []string

	if f, err := os.Stat(directory); err != nil {
		return nil, err
	} else if !f.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", directory)
	}

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == "."+inputType {
			imgPaths = append(imgPaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return imgPaths, nil
}

func convert(filePath, outputType string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	fileName := filepath.Base(filePath)
	output, err := os.Create(fileName + "." + outputType)
	if err != nil {
		return err
	}
	defer output.Close()

	switch outputType {
	case "jpg", "jpeg":
		return convertJPG(img, output)
	case "png":
		return convertPNG(img, output)
	case "gif":
		return convertGIF(img, output)
	default:
		return fmt.Errorf("%s is not a supported output type", outputType)
	}
}

func convertJPG(img image.Image, output *os.File) error {
	if err := jpeg.Encode(output, img, nil); err != nil {
		return err
	}
	return nil
}

func convertPNG(img image.Image, output *os.File) error {
	if err := png.Encode(output, img); err != nil {
		return err
	}
	return nil
}

func convertGIF(img image.Image, output *os.File) error {
	if err := gif.Encode(output, img, nil); err != nil {
		return err
	}
	return nil
}
