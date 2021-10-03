package imgconv

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type CLI struct {
	OutStream, ErrStream io.Writer
}

// validateType validates the type of the image
func validateType(t string) error {
	switch t {
	case "jpg", "jpeg", "png", "gif":
		return nil
	default:
		return fmt.Errorf("invalid type: %s", t)
	}
}

func (cli *CLI) Run() int {
	config := &Config{}
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.StringVar(&config.InputType, "input-type", "jpg", "input type[jpg|jpeg|png|gif]")
	fs.StringVar(&config.OutputType, "output-type", "png", "output type[jpg|jpeg|png|gif]")
	fs.SetOutput(cli.ErrStream)
	fs.Usage = func() {
		fmt.Fprintf(cli.ErrStream, "Usage: %s [options] DIRECTORY\n", "imgconv")
		fs.PrintDefaults()
	}

	fs.Parse(os.Args[1:])

	if err := validateType(config.InputType); err != nil {
		fmt.Fprintf(cli.ErrStream, "invalid input type: %s\n", err)
		return 1
	}

	if err := validateType(config.OutputType); err != nil {
		fmt.Fprintf(cli.ErrStream, "invalid output type: %s\n", err)
		return 1
	}

	if config.InputType == config.OutputType {
		fmt.Fprintf(cli.ErrStream, "input type and output type must be different\n")
		return 1
	}

	if fs.Arg(0) == "" {
		fmt.Fprintf(cli.ErrStream, "directory is required\n")
		return 1
	}

	config.Directory = fs.Arg(0)

	imgConv := &ImgConv{
		OutStream: cli.OutStream,
	}
	dec := NewDecoder()
	enc, err := NewEncoder(config.InputType)
	if err != nil {
		fmt.Fprintf(cli.ErrStream, "failed to create encoder: %s\n", err)
		return 1
	}
	imgConv.Run(dec, enc, config.Directory)

	return 0
}
