package download_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/download"
)

var testdataPathMap = map[int][]string{
	0: {"../testdata/003", mustGetSize("../testdata/003")},
	1: {"../testdata/z4d4kWk.jpg", mustGetSize("../testdata/z4d4kWk.jpg")},
	// 2 : "../documents/http.request.txt",
}

func TestDownloader_SingleProcess(t *testing.T) {
	t.Helper()

	ts, clean := newTestServer(t, nonRangeAccessHandler, 0)
	defer clean()

	// get a url.URL object
	urlObj := getURLObject(t, ts.URL)

	// get a file for output
	output, clean := makeTempFile(t)
	defer clean()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// get the file size to be downloaded.
	size := getSizeForTest(t, resp)

	// this test is non-parallel download.
	part := size
	procs := uint(1)
	isPara := false
	tmpDirName := ""
	ctx := context.Background()

	t.Run("case 1", func(t *testing.T) {
		err := download.Downloader(urlObj, output, size, part, procs, isPara, tmpDirName, ctx)
		if err != nil {
			t.Error(err)
		}

		actual := new(bytes.Buffer).Bytes()
		_, err = output.Read(actual)
		if err != nil {
			t.Error(err)
		}

		expected, err := os.ReadFile(testdataPathMap[0][0])
		if err != nil {
			t.Error(err)
		}

		if reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %s, but got %s", expected, actual)
		}
	})
}

func TestDownloader_SingleProcessTimeout(t *testing.T) {
	t.Helper()

	ts, clean := newTestServer(t, nonRangeAccessTooLateHandler, 0)
	defer clean()

	// get a url.URL object
	urlObj := getURLObject(t, ts.URL)

	// get a file for output
	output, clean := makeTempFile(t)
	defer clean()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// get the file size to be downloaded.
	size := getSizeForTest(t, resp)

	// this test is non-parallel download.
	part := size
	procs := uint(1)
	isPara := false
	tmpDirName := ""
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	t.Run("case 1", func(t *testing.T) {
		actual := download.Downloader(urlObj, output, size, part, procs, isPara, tmpDirName, ctx)
		if err != nil {
			t.Error(err)
		}
		expected := fmt.Errorf("request.Request err: Get \"%s\": %w", urlObj, context.DeadlineExceeded)
		if actual.Error() != expected.Error() {
			t.Errorf("expected %s, \nbut got %s", expected, actual)
		}
	})
}

func TestDownloader_ParallelProcess(t *testing.T) {
	t.Helper()

	// make a server for test.
	ts, clean := newTestServer(t, rangeAccessHandler, 1)
	defer clean()

	// get a url.URL object
	urlObj := getURLObject(t, ts.URL)

	// get a file for output
	output, clean := makeTempFile(t)
	defer clean()

	// get a response from test server.
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// get the file size to be downloaded.
	size := getSizeForTest(t, resp)

	// this test is parallel download.
	procs := uint(runtime.NumCPU())
	part := size / procs
	isPara := true
	tmpDirName := "test"
	ctx := context.Background()

	t.Run("case 1", func(t *testing.T) {
		if err := os.Mkdir(tmpDirName, 0775); err != nil {
			t.Error(err)
		}
		err := download.Downloader(urlObj, output, size, part, procs, isPara, tmpDirName, ctx)
		if err != nil {
			t.Errorf("err: %w", err)
		}
		defer func() {
			err := os.RemoveAll(tmpDirName)
			if err != nil {
				t.Error(err)
			}
		}()
	})
}

func newTestServer(t *testing.T,
	handler func(t *testing.T, w http.ResponseWriter, r *http.Request, testDataKey int),
	testDataKey int) (*httptest.Server, func()) {

	t.Helper()

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			handler(t, w, r, testDataKey)
		},
	))

	return ts, func() { ts.Close() }
}

func nonRangeAccessHandler(t *testing.T, w http.ResponseWriter, r *http.Request, testDataKey int) {
	t.Helper()

	body, err := os.ReadFile(testdataPathMap[testDataKey][0])
	if err != nil {
		t.Fatal(err)
	}
	w.Header().Set("Content-Length", testdataPathMap[testDataKey][1])
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, body)
}

func nonRangeAccessTooLateHandler(t *testing.T, w http.ResponseWriter, r *http.Request, testDataKey int) {
	t.Helper()

	body, err := os.ReadFile(testdataPathMap[testDataKey][0])
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)
	w.Header().Set("Content-Length", testdataPathMap[testDataKey][1])
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, body)
}

func rangeAccessHandler(t *testing.T, w http.ResponseWriter, r *http.Request, testDataKey int) {
	t.Helper()

	w.Header().Set("Content-Length", testdataPathMap[testDataKey][1])
	w.Header().Set("Access-Range", "bytes")

	rangeHeader := r.Header.Get("Range")

	body := retBody(t, rangeHeader, testdataPathMap[testDataKey][0])
	w.WriteHeader(http.StatusPartialContent)
	fmt.Fprint(w, body)
}

func retBody(t *testing.T, rangeHeader string, testDataPath string) []byte {
	t.Helper()

	b, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatal(err)
	}

	if rangeHeader == "" {
		return b
	}

	rangeVals := strings.Split(rangeHeader, "=")
	if rangeVals[0] != "bytes" {
		t.Fatal(errors.New("err : Range header expected \"bytes\""))
	}

	rangeBytes := strings.Split(rangeVals[1], "-")
	start, err := strconv.Atoi(rangeBytes[0])
	if err != nil {
		t.Fatal(err)
	}

	end, err := strconv.Atoi(rangeBytes[1])
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Printf("length of b: %d\n", len(b))
	// fmt.Printf("length of b[start:end+1]: %d\n", len(b[start:end+1]))
	return b[start : end+1]
}

func getURLObject(t *testing.T, urlStr string) *url.URL {
	t.Helper()

	urlObj, err := url.ParseRequestURI(urlStr)
	if err != nil {
		t.Error(err)
	}

	return urlObj
}

func makeTempFile(t *testing.T) (*os.File, func()) {
	t.Helper()

	dir, err := ioutil.TempDir("", "test_download")
	if err != nil {
		t.Fatal(err)
	}

	out, err := os.Create(dir + "/test")
	if err != nil {
		t.Fatal(err)
	}

	return out,
		func() {
			err = out.Close()
			if err != nil {
				t.Fatal(err)
			}
			err = os.RemoveAll(dir)
			if err != nil {
				t.Fatal(err)
			}
		}
}

// GetSize returns size from response header.
func getSizeForTest(t *testing.T, r *http.Response) uint {
	t.Helper()

	contLen, is := r.Header["Content-Length"]
	// fmt.Println(h)
	if !is {
		t.Errorf("cannot find Content-Length header")
	}

	ret, err := strconv.ParseUint(contLen[0], 10, 32)
	if err != nil {
		t.Error(err)
	}
	return uint(ret)
}

func mustGetSize(path string) string {

	fileinfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	return strconv.Itoa(int(fileinfo.Size()))
}
