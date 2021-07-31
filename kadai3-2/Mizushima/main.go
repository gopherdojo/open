package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	flags "github.com/jessevdk/go-flags"

	"github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/download"
	"github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/getheader"
	"github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/listen"
	"github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/request"
)

// struct for options
type Options struct {
	Help   bool   `short:"h" long:"help"`
	Procs  uint   `short:"p" long:"procs"`
	Output string `short:"o" long:"output" default:"./"`
	Tm     int    `short:"t" long:"timeout" default:"120"`
}

// parse options
func (opts *Options) parse(argv []string) ([]string, error) {
	p := flags.NewParser(opts, flags.PrintErrors)
	args, err := p.ParseArgs(argv)
	if err != nil {
		_, err2 := os.Stderr.Write(opts.usage())
		if err2 != nil {
			return nil, fmt.Errorf("%s: invalid command line options: cannot print usage: %s", err, err2)
		}
		return nil, fmt.Errorf("%w: invalid command line options", err)
	}

	return args, nil
}

// usage prints a description of avilable options
func (opts Options) usage() []byte {
	buf := bytes.Buffer{}

	fmt.Fprintln(&buf,
		`Usage: paraDW [options] URL (URL2, URL3, ...)

	Options:
	-h,   --help               print usage and exit
	-p,   --procs <num>        the number of split to download (default: the number of CPU cores)
	-o,   --output <filename>  path of the file downloaded (default: current directory)
	-t,   --timeout <num>      Time limit of return of http response in seconds (default: 120)
	`,
	)

	return buf.Bytes()
}

func main() {

	// parse options
	var opts Options
	argv := os.Args[1:]
	if len(argv) == 0 {
		if _, err := os.Stdout.Write(opts.usage()); err != nil {
			log.Fatalf("err: %s: %s\n", errors.New("no options"), err)
		}
		log.Fatalf("err: %s\n", errors.New("no options"))
	}

	urlsStr, err := opts.parse(argv)
	if err != nil {
		log.Fatalf("err: %s\n", err)
	}

	var urls []*url.URL
	for _, u := range urlsStr {
		url, err := url.ParseRequestURI(u)
		if err != nil {
			log.Fatalf("err: url.ParseRequestURI: %s\n", err)
		}
		urls = append(urls, url)
	}

	fmt.Printf("timeout: %d\n", opts.Tm)

	if opts.Help {
		if _, err := os.Stdout.Write(opts.usage()); err != nil {
			log.Fatalf("err: cannot print usage: %s", err)
		}
		log.Fatal(errors.New("print usage"))
	}

	// if procs was inputted, set the number of runtime.NumCPU() to opts.Procs.
	if opts.Procs == 0 {
		opts.Procs = uint(runtime.NumCPU())
	}

	// if opts.Output inputted and the end of opts.Output is not '/',
	// add '/'.
	if len(opts.Output) > 0 && opts.Output[len(opts.Output)-1] != '/' {
		opts.Output += "/"
	}

	// download from each url in urls
	for i, urlObj := range urls {
		downloadFromUrl(i, opts, urlObj)
	}
}

// downloadFromUrl does the download processing from url object.
func downloadFromUrl(i int, opts Options, urlObj *url.URL) {

	// make a timeout context from a empty context
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), time.Duration(opts.Tm)*time.Second)
	defer cancelTimeout()

	// send "HEAD" request, and gets response.
	resp, err := request.Request(ctxTimeout, "HEAD", urlObj.String(), "", "")
	if err != nil {
		log.Fatalf("err: %s\n", err)
	}

	// get the size from the response header.
	fileSize, err := getheader.GetSize(resp)
	if err != nil {
		log.Fatalf("err: getheader.GetSize: %s\n", err)
	}
	if err = resp.Body.Close(); err != nil {
		log.Fatalf("err: %s", err)
	}

	// get how many bytes to download at a time
	partial := fileSize / opts.Procs

	outputPath := opts.Output + filepath.Base(urlObj.String())
	// if there is the same file in opts.Output, delete that file in advance.
	if isExists(outputPath) {
		err := os.Remove(outputPath)
		if err != nil {
			log.Fatalf("err: isExists: os.Remove: %s\n", err)
		}
	}

	// make a file for download
	out, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("err: os.Create: %s\n", err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Fatalf("err: %s", err)
		}
	}()

	// make a temporary directory for parallel download
	tmpDirName := opts.Output + strconv.Itoa(i)
	err = os.Mkdir(tmpDirName, 0777)
	if err != nil {
		if err3 := out.Close(); err3 != nil {
			log.Fatalf("err: %s", err3)
		}
		if err2 := os.Remove(opts.Output + filepath.Base(urlObj.String())); err2 != nil {
			log.Fatalf("err: os.Mkdir: %s\nerr: os.Remove: %s\n", err, err2)
		}
		log.Fatalf("err: os.Mkdir: %s\n", err)
	}

	clean := func() {
		if err := out.Close(); err != nil {
			log.Fatalf("err: out.Close: %s\n", err)
		}
		// delete the tmporary directory
		if err := os.RemoveAll(tmpDirName); err != nil {
			log.Fatalf("err: RemoveAll: %s\n", err)
		}
		if err := os.Remove(opts.Output + filepath.Base(urlObj.String())); err != nil {
			log.Fatalf("err: os.Remove: %s\n", err)
		}
	}
	ctx, cancel := listen.Listen(ctxTimeout, os.Stdout, clean)
	defer cancel()

	var isPara bool = true
	_, err = getheader.ResHeader(os.Stdout, resp, "Accept-Ranges")
	if err != nil && err.Error() == "cannot find Accept-Ranges header" {
		isPara = false
	} else if err != nil {
		clean()
		log.Fatalf("err: getheader.ResHeader: %s\n", err)
	}

	// drive a download process
	err = download.Downloader(urlObj, out, fileSize, partial, opts.Procs, isPara, tmpDirName, ctx)
	if err != nil {
		log.Fatalf("err: %s\n", err)
	}

	fmt.Printf("download complete: %s\n", urlObj.String())

	// Merge the temporary files into "out", when parallel download executed.
	if isPara {
		err = MergeFiles(tmpDirName, opts.Procs, fileSize, out)
		if err != nil {
			log.Fatalf("err: MergeFiles: %s\n", err)
		}
	}

	// delete the tmporary directory only
	if err := os.RemoveAll(tmpDirName); err != nil {
		log.Fatalf("err: RemoveAll: %s\n", err)
	}
}

// MergeFiles merges temporary files made for parallel download into "output".
func MergeFiles(tmpDirName string, procs, fileSize uint, output *os.File) error {
	for i := uint(0); i < procs; i++ {

		body, err := os.ReadFile(tmpDirName + "/" + strconv.Itoa(int(i)))
		if err != nil {
			return err
		}

		if _, err = fmt.Fprint(output, string(body)); err != nil {
			return err
		}
		fmt.Printf("target file: %s, len=%d written\n", output.Name(), len(string(body)))
	}
	return nil
}

// isExists returns the file exists or not, in 'path'.
func isExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
