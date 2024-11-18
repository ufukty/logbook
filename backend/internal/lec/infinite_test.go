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

	inf := New(ts[0], average.Year, average.Day)
	for _, t := range ts {
		fmt.Println(t)
		inf.Save(t, 1)
	}

	q, err := inf.Query(ts[0], ts[len(ts)-1])
	if err != nil {
		t.Fatal(fmt.Errorf("act, query: %w", err))
	}
	fmt.Println(q)
}
