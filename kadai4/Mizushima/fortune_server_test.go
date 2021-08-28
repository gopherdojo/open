package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOmikujiHandler(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	omikujiHandler(w, r)
	rw := w.Result()
	defer func() {
		if err := rw.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if rw.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", rw.StatusCode)
	}

	var res Res
	dec := json.NewDecoder(rw.Body)
	if err := dec.Decode(&res); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	s := []string{"大吉", "中吉", "小吉", "凶"}
	if res.Result == "" || !contains(s, res.Result) {
		t.Fatal("Error: json returned is not valid")
	}
}

func TestResult(t *testing.T) {
	cases := map[string]struct {
		date     time.Time
		index    int
		expected string
	}{
		"the 1st of the first three days of the new year.": {
			date:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			index:    1,
			expected: "大吉",
		},
		"the 2nd of the first three days of the new year.": {
			date:     time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			index:    1,
			expected: "大吉",
		},
		"the 3rd of the first three days of the new year.": {
			date:     time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
			index:    1,
			expected: "大吉",
		},
		"the 4th of the first three days of the new year, and the index is 5.": {
			date:     time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
			index:    5,
			expected: "凶",
		},
		"August 25 and the index is 0": {
			date:     time.Date(2021, 8, 25, 0, 0, 0, 0, time.Local),
			index:    0,
			expected: "大吉",
		},
		"August 25 and the index is 1": {
			date:     time.Date(2021, 8, 25, 0, 0, 0, 0, time.Local),
			index:    1,
			expected: "中吉",
		},
		"August 25 and the index is 2": {
			date:     time.Date(2021, 8, 25, 0, 0, 0, 0, time.Local),
			index:    2,
			expected: "中吉",
		},
		"August 25 and the index is 3": {
			date:     time.Date(2021, 8, 25, 0, 0, 0, 0, time.Local),
			index:    3,
			expected: "小吉",
		},
		"August 25 and the index is 4": {
			date:     time.Date(2021, 8, 25, 0, 0, 0, 0, time.Local),
			index:    4,
			expected: "小吉",
		},
		"August 25 and the index is 5": {
			date:     time.Date(2021, 8, 25, 0, 0, 0, 0, time.Local),
			index:    5,
			expected: "凶",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := result(c.index, c.date)
			if actual != c.expected {
				t.Fatalf("expected %v, but got %v", c.expected, actual)
			}
		})
	}
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
