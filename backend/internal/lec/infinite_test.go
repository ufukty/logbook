package lec

import (
	"fmt"
	"logbook/internal/average"
	"testing"
	"time"
)

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

	inf := NewInfinite(from, average.Day)
	for _, t := range ts {
		inf.Save(t, 1)
	}

	err = dump(inf)
	if err != nil {
		t.Fatal(fmt.Errorf("dump: %w", err))
	}

	type (
		input struct {
			from, to time.Time
		}
		output struct {
			q   int
			err bool
		}
		tc struct {
			name string
			input
			output
		}
	)
	tcs := []tc{
		{"Empty range", input{from, from}, output{err: true}},

		{"Overflow range", input{from, to.Add(average.Day * 10)}, output{err: true}},
		{"Out of bound (1 day after)", input{to.Add(average.Day), to.Add(average.Day * 2)}, output{err: true}},
		{"Out of bound (2 days after)", input{to.Add(average.Day * 2), to.Add(average.Day * 3)}, output{err: true}},
		{"Out of bound (1 week after)", input{to.Add(average.Week), to.Add(average.Week * 2)}, output{err: true}},
		{"Out of bound (1 week before)", input{from.Add(-3 * average.Week), from.Add(-2 * average.Week)}, output{err: true}},

		{"Full range (adjusted)", input{from, to}, output{q: 1000}},
		{"Full range (year)", input{from, from.Add(average.Year + average.Day)}, output{q: 1000}},

		{"Jun, 1st", input{from, from.Add(average.Day)}, output{q: 3}},
		{"Feb, 1st", input{parse(t, "2024-02-01T00:00:00+00:00"), parse(t, "2024-02-02T00:00:00+00:00")}, output{q: 3}},
		{"Feb, 1st-30th", input{parse(t, "2024-02-01T00:00:00+00:00"), parse(t, "2024-03-01T00:00:00+00:00")}, output{q: 85}},
		{"Mar, 1st-31th", input{parse(t, "2024-03-01T00:00:00+00:00"), parse(t, "2024-04-01T00:00:00+00:00")}, output{q: 72}},
		{"Feb, 1st to Mar, 31th", input{parse(t, "2024-02-01T00:00:00+00:00"), parse(t, "2024-04-01T00:00:00+00:00")}, output{q: 157}},
		{"July, 2nd", input{parse(t, "2024-07-02T00:00:00+00:00"), parse(t, "2024-07-03T00:00:00+00:00")}, output{q: 6}},
		{"July, 4th", input{parse(t, "2024-07-04T00:00:00+00:00"), parse(t, "2024-07-05T00:00:00+00:00")}, output{q: 3}},
		{"July, 2nd-5th", input{parse(t, "2024-07-02T00:00:00+00:00"), parse(t, "2024-07-05T00:00:00+00:00")}, output{q: 10}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
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
				t.Errorf("assert, expected %d, got %d", tc.output.q, q)
			}
		})
	}
}
