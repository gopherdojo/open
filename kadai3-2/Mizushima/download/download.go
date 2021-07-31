// download package implements parallel download and non-parallel
// download.
package download

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"

	"github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/request"
	"golang.org/x/sync/errgroup"
)

// PDownloader is user-defined struct for the download process.
// It's not limited to parallel downloads.
type PDownloader struct {
	url      *url.URL // URL for the download
	output   *os.File // Where to save the downloaded file
	fileSize uint     // size of the downloaded file
	part     uint     // Number of divided bytes
	procs    uint     // Number of parallel download process
}

// newPDownloader is constructor for PDownloader.
func newPDownloader(url *url.URL, output *os.File, fileSize uint, part uint, procs uint) *PDownloader {
	return &PDownloader{
		url:      url,
		output:   output,
		fileSize: fileSize,
		part:     part,
		procs:    procs,
	}
}

// Downloader gets elements of PDownloader, the download is parallel or not, temprary
// directory name and context.Context, and drives DownloadFile method if isPara is false
// or PDownload if isPara is true.
//
func Downloader(url *url.URL,
	output *os.File, fileSize uint, part uint, procs uint, isPara bool,
	tmpDirName string, ctx context.Context) error {
	pd := newPDownloader(url, output, fileSize, part, procs)
	if !isPara {
		fmt.Printf("%s do not accept range access: downloading by single process\n", url)
		err := pd.DownloadFile(ctx)
		if err != nil {
			return err
		}
	} else {
		grp, ctx := errgroup.WithContext(ctx)
		if err := pd.PDownload(grp, tmpDirName, procs, ctx); err != nil {
			return err
		}

		if err := grp.Wait(); err != nil {
			return err
		}
	}
	return nil
}

// DownloadFile drives a non-parallel download
func (pd *PDownloader) DownloadFile(ctx context.Context) (err error) {

	resp, err := request.Request(ctx, "GET", pd.url.String(), "", "")
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()

	_, err = io.Copy(pd.output, resp.Body)
	if err != nil {
		return
	}

	return nil
}

// PDownload drives parallel download. downloaded file is in temporary
// directory named tmpDirName.
func (pd *PDownloader) PDownload(grp *errgroup.Group,
	tmpDirName string, procs uint, ctx context.Context) error {
	var start, end, idx uint

	for idx = uint(0); idx < procs; idx++ {
		if idx == 0 {
			start = 0
		} else {
			start = idx*pd.part + 1
		}

		// if idx is the end
		if idx == pd.procs-1 {
			end = pd.fileSize
		} else {
			end = (idx + 1) * pd.part
		}

		// idxを代入し直す
		// https://qiita.com/harhogefoo/items/7ccb4e353a4a01cfa773
		idx := idx
		// fmt.Printf("start: %d, end: %d, pd.part: %d\n", start, end, pd.part)
		bytes := fmt.Sprintf("bytes=%d-%d", start, end)

		grp.Go(func() error {
			fmt.Printf("grp.Go: tmpDirName: %s, bytes %s, idx: %d\n", tmpDirName, bytes, idx)
			return pd.ReqToMakeCopy(tmpDirName, bytes, idx, ctx)
		})
	}
	return nil
}

// ReqToMakeCopy sends a "GET" request with "Range" field with "bytes" range.
// And gets response and make a copy to a temprary file in temprary directory from response body.
//
func (pd *PDownloader) ReqToMakeCopy(tmpDirName, bytes string, idx uint, ctx context.Context) (err error) {
	// fmt.Printf("ReqToMakeCopy: tmpDirName: %s, bytes %s, idx: %d\n", tmpDirName, bytes, idx)
	resp, err := request.Request(ctx, "GET", pd.url.String(), "Range", bytes)
	if err != nil {
		return err
	}

	tmpOut, err := os.Create(tmpDirName + "/" + strconv.Itoa(int(idx)))
	if err != nil {
		return err
	}
	// fmt.Printf("tmpOut.Name(): %s\n", tmpOut.Name())
	defer func() {
		err = tmpOut.Close()
	}()

	// b := make([]byte, 1000)
	// resp.Body.Read(b)
	// fmt.Printf("resp.body: %s\n", string(b))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// fmt.Printf("err: %s\n", err)
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			return err
		}
	}

	// fmt.Printf("response body: length: %d\n", len(body))

	length, err := tmpOut.Write(body)
	if err != nil {
		return err
	}
	fmt.Printf("%d/%d was downloaded len=%d\n", idx, pd.procs, length)
	return nil
}
