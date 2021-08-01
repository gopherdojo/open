/*
Package mypkg is my package.
*/
package converter

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
)

// Converter is convert images source to destination
type Converter struct {
	Src         string                                          // Source image extension
	Dest        string                                          // Destination image extension
	DecodeFuncs map[string]func(io.Reader) (image.Image, error) // DecodeFuncs use decode
	EncodeFuncs map[string]interface{}                          // EncodeFuncs use encode
}

// NewConverter is return new converter
func NewConverter(src, dest string) (*Converter, error) {
	allowExtensions := []string{"jpg", "png", "gif"}
	contains := func(s []string, e string) bool {
		for _, a := range s {
			if a == e {
				return true
			}
		}
		return false
	}
	if contains(allowExtensions, src) && contains(allowExtensions, dest) {
		return &Converter{
			Src:  src,
			Dest: dest,
			DecodeFuncs: map[string]func(io.Reader) (image.Image, error){
				"jpg": jpeg.Decode,
				"gif": gif.Decode,
				"png": png.Decode,
			},
			EncodeFuncs: map[string]interface{}{
				"jpg": jpeg.Encode,
				"gif": gif.Encode,
				"png": png.Encode,
			}}, nil
	}
	return nil, fmt.Errorf("jpg, png, gif extenstions are only supported")
}

// Convert is convert image file. If file is not image return error,
func (c *Converter) Convert(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	srcImage, err := c.DecodeFuncs[c.Src](f)
	if err != nil {
		return err
	}
	out, err := os.Create(filePath[0:len(filePath)-len(path.Ext(filePath))] + "." + c.Dest)
	if err != nil {
		return err
	}
	switch fn := c.EncodeFuncs[c.Dest].(type) {
	case func(io.Writer, image.Image) error:
		err = fn(out, srcImage)
	case func(io.Writer, image.Image, *jpeg.Options) error:
		err = fn(out, srcImage, &jpeg.Options{Quality: 100})
	case func(io.Writer, image.Image, *gif.Options) error:
		err = fn(out, srcImage, &gif.Options{NumColors: 256})
	}
	if err != nil {
		return err
	}
	return nil
}
