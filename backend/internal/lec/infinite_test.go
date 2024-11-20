package lec

import (
	"bufio"
	"fmt"
	"logbook/internal/average"
	"os"
	"testing"
	"time"
)

func read() ([]time.Time, error) {
	f, err := os.Open("testdata/ts.txt")
	if err != nil {
		return nil, fmt.Errorf("prep, load doc: %w", err)
	}
	defer f.Close()

	r := []time.Time{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t, err := time.Parse("2006-01-02T15:04:05-07:00", scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("parse: %w", err)
		}
		r = append(r, t)
	}

	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("scanner: %w", err)
	}

	return r, nil
}

func TestInfinite(t *testing.T) {
	ts, err := read()
	if err != nil {
		t.Fatal(fmt.Errorf("prep: %w", err))
	}

	var (
		from = ts[0].Truncate(average.Day)
		to   = ts[len(ts)-1].Truncate(average.Day)
	)

	fmt.Println(ts[0])
	fmt.Println(from)
	fmt.Println(ts[len(ts)-1])
	fmt.Println(to)

	inf := New(from, average.Year, average.Day)
	for _, t := range ts {
		inf.Save(t, 1)
	}

	type (
		input struct {
			from, to time.Time
		}
		output struct {
			number int
		}
		tc struct {
			input
			output
		}
	)
	tcs := map[string]tc{
		"full range": {input{from, to.Add(average.Day)}, output{1000}},
		"overflow":   {input{from, to.Add(average.Day * 10)}, output{1000}},
		"half range": {input{from, from.Add(time.Duration(int64(0.5 * float64(int64(to.Sub(from).Nanoseconds())))))}, output{496}},
	}

	for tn, tc := range tcs {
		t.Run(tn, func(t *testing.T) {
			q, err := inf.Query(tc.from, tc.to)
			if err != nil {
				t.Fatal(fmt.Errorf("act, query: %w", err))
			}
			if q != tc.output.number {
				t.Errorf("expected to get '%d', got '%d'", tc.output.number, q)
			}
		})
	}
}
