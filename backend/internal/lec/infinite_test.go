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

const testformat = "2006-01-02T15:04:05-07:00"

func read() ([]time.Time, error) {
	f, err := os.Open("testdata/ts.txt")
	if err != nil {
		return nil, fmt.Errorf("prep, load doc: %w", err)
	}
	defer f.Close()

	r := []time.Time{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t, err := time.Parse(testformat, scanner.Text())
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
	fmt.Fprintln(f, inf.Dump())
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

	fmt.Println(ts[0])
	fmt.Println(from)
	fmt.Println(ts[len(ts)-1])
	fmt.Println(to)

	inf := New(from, average.Year+average.Day, average.Day)
	for _, t := range ts {
		inf.Save(t, 1)
	}

	err = dump(inf)
	if err != nil {
		t.Fatal(fmt.Errorf("dump: %w", err))
	}

	julyTwo, err := time.Parse(testformat, "2024-07-02T14:45:22+03:00")
	if err != nil {
		t.Fatal(fmt.Errorf("prep, julyTwo: %w", err))
	}
	julyTwo = julyTwo.Truncate(average.Day)

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
		"Empty range": {input{from, from}, output{0}},

		"Full range (adjusted)": {input{from, to.Add(average.Day)}, output{len(ts)}},
		"Full range (year)":     {input{from, from.Add(average.Year + average.Day)}, output{len(ts) - 1}},

		"Overflow range":               {input{from, to.Add(average.Day * 10)}, output{len(ts)}},
		"Out of bound":                 {input{to.Add(average.Day), to.Add(average.Day * 2)}, output{0}},
		"Out of bound (2 days after)":  {input{to.Add(average.Day * 2), to.Add(average.Day * 3)}, output{0}},
		"Out of bound (1 week after)":  {input{to.Add(average.Week), to.Add(average.Week * 2)}, output{0}},
		"Out of bound (1 week before)": {input{from.Add(-3 * average.Week), from.Add(-2 * average.Week)}, output{0}},

		"Single day": {input{from, from.Add(average.Day)}, output{1}},
		"July 2nd":   {input{julyTwo, julyTwo.Add(average.Day)}, output{6}},
	}

	for tn, tc := range tcs {
		t.Run(tn, func(t *testing.T) {
			q, err := inf.Query(tc.input.from, tc.input.to)
			if err != nil {
				t.Fatalf("query failed: %v", err)
			}
			if q != tc.output.number {
				t.Errorf("expected %d, got %d, input range [%s, %s]", tc.output.number, q, tc.input.from, tc.input.to)
			}
		})
	}
}
