package lec

import (
	"bufio"
	"fmt"
	"logbook/internal/average"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

const testFileTimeFormat = "2006-01-02T15:04:05-07:00"

func read() ([]time.Time, error) {
	f, err := os.Open("testdata/ts.txt")
	if err != nil {
		return nil, fmt.Errorf("prep, load doc: %w", err)
	}
	defer f.Close()

	r := []time.Time{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t, err := time.Parse(testFileTimeFormat, scanner.Text())
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

type dumpable interface {
	dump() string
}

func dump(inf dumpable) error {
	err := os.MkdirAll("testresults", 0755)
	if err != nil {
		return fmt.Errorf("mkdir testresults: %w", err)
	}
	f, err := os.Create(filepath.Join("testresults", fmt.Sprintf("%s.txt", time.Now().Format("2006-01-02-15-04-05"))))
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	fmt.Fprint(f, inf.dump())
	return nil
}

func parse(t *testing.T, s string) time.Time {
	tt, err := time.Parse(testFileTimeFormat, s)
	if err != nil {
		t.Fatal(fmt.Errorf("parse: %w", err))
	}
	return tt.Truncate(average.Day)
}

func TestMinpow(t *testing.T) {
	tcs := map[int]int{
		0:  0,
		1:  1,
		2:  1,
		3:  2,
		4:  2,
		5:  3,
		8:  3,
		9:  4,
		16: 4,
		17: 5,
		32: 5,
		33: 6,
		64: 6,
		65: 7,
	}

	for i, e := range tcs {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g := minpow(i)
			if g != e {
				t.Fatalf("expected %d != got %d", e, g)
			}
		})
	}
}
