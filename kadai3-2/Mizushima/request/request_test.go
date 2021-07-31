package request_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/request"
)

var testdataPathMap = map[int][]string{
	0: {"../testdata/003", mustGetSize("../testdata/003")},
	1: {"../testdata/z4d4kWk.jpg", mustGetSize("../testdata/z4d4kWk.jpg")},
}

func Test_RequestStandard(t *testing.T) {
	t.Helper()

	cases := map[string]struct {
		key      int // key for testdataPathMap
		handler  func(t *testing.T, w http.ResponseWriter, r *http.Request, testDataKey int)
		expected *http.Response
	}{
		"case 1": {
			key:     0,
			handler: nonRangeAccessHandler,
			expected: &http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: map[string][]string{
					"Content-Length": {testdataPathMap[0][1]},
					"Date":           {mustTimeLayout(t, time.Now())},
				},
			},
		},
		"case 2": {
			key:     1,
			handler: rangeAccessHandler,
			expected: &http.Response{
				Status:     "206 Partial Content",
				StatusCode: 206,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: map[string][]string{
					"Access-Range":   {"bytes"},
					"Content-Length": {testdataPathMap[1][1]},
					"Date":           {mustTimeLayout(t, time.Now())},
				},
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			ts, clean := newTestServer(t, c.handler, c.key)
			defer clean()
			actual, err := request.Request(context.Background(), "GET", ts.URL, "", "")
			if err != nil {
				t.Fatal(err)
			}
			// fmt.Println("actual:", actual.Header)
			if !reflect.DeepEqual(actual.Header, c.expected.Header) {
				dumped_expected, err := httputil.DumpResponse(c.expected, false)
				if err != nil {
					t.Fatal(err)
				}
				dumped_actual, err := httputil.DumpResponse(actual, false)
				if err != nil {
					t.Fatal(err)
				}
				t.Errorf("expected,\n%vbut got,\n%v", string(dumped_expected), string(dumped_actual))
			}
		})
	}
}

func Test_RequestTimeout(t *testing.T) {
	t.Helper()

	name := "case timeout"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	t.Run(name, func(t *testing.T) {
		ts, clean := newTestServer(t, nonRangeAccessTooLateHandler, 1)
		defer clean()

		expected := fmt.Errorf("request.Request err: Get \"%s\": %w", ts.URL, context.DeadlineExceeded)
		_, err := request.Request(ctx, "GET", ts.URL, "", "")

		// fmt.Println("actual:\n", actual)
		if err.Error() != expected.Error() {
			t.Errorf("expected %s, but got %s", expected.Error(), err.Error())
		}
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

// mustGetSize returns the size of the file in "path" as a string for "Content-Length" in http header.
func mustGetSize(path string) string {

	fileinfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	return strconv.Itoa(int(fileinfo.Size()))
}

// mustTimeLayout returns the time now in format like this : "Mon, 12 Jul 2021 09:22:22 GMT"
func mustTimeLayout(t *testing.T, tm time.Time) string {
	t.Helper()

	// get the gmt time
	location, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatal(err)
	}
	tm = tm.In(location)

	// the layout is like this : "Mon, 12 Jul 2021 09:22:22 GMT"
	return tm.Format(time.RFC1123)
}
