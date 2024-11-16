package challenge

import (
	"fmt"
	"testing"
	"time"
)

func Example_encode() {
	fmt.Println(encode([]byte("Hello world"))) // Output: JBSWY3DPEB3W64TMMQ
}

func TestSolve(t *testing.T) {
	type tc struct {
		input_l, input_n int
	}

	tcs := []tc{
		{20, 1},
		{20, 2},
		{20, 3},
		{20, 4},

		{50, 1},
		{50, 2},
		{50, 3},
		{50, 4},

		{100, 1},
		{100, 2},
		{100, 3},
		{100, 4},

		{200, 1},
		{200, 2},
		{200, 3},
		{200, 4},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%d-%d", tc.input_l, tc.input_n), func(t *testing.T) {
			defer func(s time.Time) { t.Log("took", time.Since(s)) }(time.Now())
			ch, err := NewChallenge(tc.input_l, tc.input_n)
			if err != nil {
				t.Fatalf("failed to create challenge: %v", err)
			}
			t.Log("Challange:", ch)
			solved, err := Solve(ch.N, ch.Que, ch.Hash)
			if err != nil {
				t.Fatalf("failed to solve challenge: %v", err)
			}
			if ch.Original != solved {
				t.Fatalf("assert got %s want %s", solved, ch.Original)
			}
		})
	}
}

func BenchmarkDeviation(b *testing.B) {
	ts := []time.Duration{}
	for i := 0; i < b.N; i++ {
		t := time.Now()
		ch, err := NewChallenge(500, 4)
		if err != nil {
			b.Fatal(fmt.Errorf("prep: %w", err))
		}
		solved, err := Solve(ch.N, ch.Que, ch.Hash)
		if err != nil {
			b.Fatal(fmt.Errorf("act, solve: %w", err))
		}
		if ch.Original != solved {
			b.Fatalf("assert got %s want %s", solved, ch.Original)
		}
		ts = append(ts, time.Since(t))
	}
	b.Log(statistics(ts))
}

func statistics(ts []time.Duration) string {
	// Calculate statistics
	total := time.Duration(0)
	min := ts[0]
	max := ts[0]

	for _, t := range ts {
		total += t
		if t < min {
			min = t
		}
		if t > max {
			max = t
		}
	}
	avg := total / time.Duration(len(ts))

	// Calculate standard deviation
	var varianceSum float64
	for _, t := range ts {
		diff := float64(t-avg) / float64(time.Millisecond)
		varianceSum += diff * diff
	}
	stdDev := time.Duration(float64(time.Millisecond) * (varianceSum / float64(len(ts))))

	return fmt.Sprintf("Benchmark results: runs=%d avg=%v min=%v max=%v stdDev=%v", len(ts), avg, min, max, stdDev)
}
