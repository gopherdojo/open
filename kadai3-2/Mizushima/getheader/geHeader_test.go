package getheader_test

import (
	"errors"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/getheader"
)

var heads1 = []string{
	"Content-Length",
	"Accept-Ranges",
	"Content-Type",
	"Access-Control-Allow-Methods",
}

var heads2 = []string{
	"Accept-Ranges",
	"Content-Type",
	"Access-Control-Allow-Methods",
}

var vals1 = [][]string{{"146515"}, {"bytes"}, {"image/jpeg"}, {"GET", "OPTIONS"}}

var vals2 = [][]string{{"bytes"}, {"image/jpeg"}, {"GET", "OPTIONS"}}

func Test_ResHeader(t *testing.T) {
	t.Helper()

	cases := []struct {
		name     string
		input    string
		heads    []string
		vals     [][]string
		expected []string
	}{
		{
			name:     "case 1",
			input:    "Content-Length",
			heads:    heads1,
			vals:     vals1,
			expected: []string{"146515"},
		},
		{
			name:     "case 2",
			input:    "Accept-Ranges",
			heads:    heads2,
			vals:     vals2,
			expected: []string{"bytes"},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			resp, err := makeResponse(t, c.heads, c.vals)
			if err != nil {
				t.Error(err)
			}
			actual, err := getheader.ResHeader(os.Stdout, resp, c.input)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf("expected %s, but got %s", c.expected, actual)
			}
		})
	}
}

func Test_GetSize(t *testing.T) {

	cases := []struct {
		name     string
		heads    []string
		vals     [][]string
		expected uint
	}{
		{
			name:     "case 1",
			heads:    heads1,
			vals:     vals1,
			expected: uint(146515),
		},
		{
			name:     "case 2",
			heads:    heads2,
			vals:     vals2,
			expected: 0,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			resp, err := makeResponse(t, c.heads, c.vals)

			if err != nil {
				t.Error(err)
			}
			actual, err := getheader.GetSize(resp)
			if err != nil {
				if actual != 0 || c.expected != 0 {
					t.Error(err)
				} else {
					if err.Error() != "cannot find Content-Length header" {
						t.Errorf("expected error: cannot find Content-Length header, but %w", err)
					}
				}
			}
			if actual != c.expected {
				t.Errorf("expected %d, but got %d", c.expected, actual)
			}
		})
	}
}

func makeResponse(t *testing.T, heads []string, vals [][]string) (*http.Response, error) {
	t.Helper()

	var resp = make(map[string][]string)

	if len(heads) != len(vals) {
		return nil, errors.New("expected the length of heads and vals sre same")
	}

	for i := 0; i < len(heads); i++ {
		if _, ok := resp[heads[i]]; ok {
			return nil, errors.New("Duplicate elements in heads")
		}

		resp[heads[i]] = vals[i]
	}

	return &http.Response{Header: http.Header(resp)}, nil
}
