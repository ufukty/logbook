package lec

import (
	"bufio"
	"fmt"
	"logbook/internal/average"
	"os"
	"path/filepath"
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

func dump(inf *Infinite) error {
	err := os.MkdirAll("testresults", 0755)
	if err != nil {
		return fmt.Errorf("mkdir testresults: %w", err)
	}
	f, err := os.Create(filepath.Join("testresults", fmt.Sprintf("%s.txt", time.Now().Format("2006-01-02-15-04-05"))))
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	fmt.Fprint(f, inf.Dump())
	return nil
}

func TestInfinite(t *testing.T) {
	// Use either read() or generateTestData to prepare timestamps
	ts, err := read()
	if err != nil {
		t.Fatal(fmt.Errorf("prep: %w", err))
	}

	var (
		from = ts[0].Truncate(average.Day)
		to   = ts[len(ts)-1].Truncate(average.Day).Add(average.Day)
	)

	fmt.Printf("testdata: (%s -> %s)\n", ts[0], ts[len(ts)-1])
	fmt.Printf("testdata: (%s -> %s) truncated\n", from, to)

	realdiff := ts[len(ts)-1].Sub(ts[0])
	truncateddiff := to.Sub(from)
	fmt.Printf("(real period: %s [%f years]) (Δtruncated: %s [%f years]) (diff: %s)\n",
		realdiff, float64(realdiff)/float64(average.Year),
		truncateddiff, float64(truncateddiff)/float64(average.Year),
		time.Duration(truncateddiff-realdiff).Abs(),
	)

	inf := New(from, average.Year+average.Day, average.Day)
	for _, t := range ts {
		inf.Save(t, 1)
	}

	err = dump(inf)
	if err != nil {
		t.Fatal(fmt.Errorf("dump: %w", err))
	}

	julyTwo, err := time.Parse(testFileTimeFormat, "2024-07-02T14:45:22+03:00")
	if err != nil {
		t.Fatal(fmt.Errorf("prep, julyTwo: %w", err))
	}
	julyTwo = julyTwo.Truncate(average.Day)

	type (
		input struct {
			from, to time.Time
		}
		output struct {
			q   int
			err bool
		}
		tc struct {
			input
			output
		}
	)
	tcs := map[string]tc{
		"Empty range": {input{from, from}, output{err: true}},

		"Overflow range":               {input{from, to.Add(average.Day * 10)}, output{err: true}},
		"Out of bound (1 day after)":   {input{to.Add(average.Day), to.Add(average.Day * 2)}, output{err: true}},
		"Out of bound (2 days after)":  {input{to.Add(average.Day * 2), to.Add(average.Day * 3)}, output{err: true}},
		"Out of bound (1 week after)":  {input{to.Add(average.Week), to.Add(average.Week * 2)}, output{err: true}},
		"Out of bound (1 week before)": {input{from.Add(-3 * average.Week), from.Add(-2 * average.Week)}, output{err: true}},

		"Full range (adjusted)": {input{from, to.Add(average.Day)}, output{q: 1000}},
		"Full range (year)":     {input{from, from.Add(average.Year + average.Day)}, output{q: 1000}},

		"Single day": {input{from, from.Add(average.Day)}, output{q: 1}},
		"July 2nd":   {input{julyTwo, julyTwo.Add(average.Day)}, output{q: 6}},
	}

	for tn, tc := range tcs {
		t.Run(tn, func(t *testing.T) {
			q, err := inf.Query(tc.input.from, tc.input.to)
			if tc.output.err {
				if err == nil {
					t.Fatalf("act, expected error. got '%d'", q)
				}
				return
			}
			if err != nil {
				t.Fatalf("act: %v", err)
			}
			if q != tc.output.q {
				t.Errorf("assert, expected %d, got %d, input range [%s, %s]", tc.output.q, q, tc.input.from, tc.input.to)
			}
		})
	}
}
