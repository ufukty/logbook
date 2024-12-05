package challenge

import (
	"fmt"
	"testing"
)

func TestSolveChallenge(t *testing.T) {
	for difficulty := 2; difficulty < len(alphabet)-1; difficulty++ {
		t.Run(fmt.Sprintf("difficulty %d", difficulty), func(t *testing.T) {
			c, err := CreateChallenge(difficulty)
			if err != nil {
				t.Fatalf("NewChallenge: %v", err)
			}
			t.Log("Challange:", c)

			combination, err := SolveChallenge(difficulty, c.Masked, c.Hashed)
			if err != nil {
				t.Fatalf("SolveChallenge: %v", err)
			}
			original := combination + c.Masked
			if original != c.Original {
				t.Fatalf("assert original, got %s, want %s", original, c.Original)
			}
		})
	}
}

// func BenchmarkDeviation(b *testing.B) {
// 	ts := []time.Duration{}
// 	for i := 0; i < b.N; i++ {
// 		t := time.Now()
// 		ch, err := NewChallenge(500, 3)
// 		if err != nil {
// 			b.Fatal(fmt.Errorf("prep: %w", err))
// 		}
// 		solved, err := SolveChallenge(ch.N, ch.Masked, ch.Hashed)
// 		if err != nil {
// 			b.Fatal(fmt.Errorf("act, solve: %w", err))
// 		}
// 		if ch.Combination != solved {
// 			b.Fatalf("assert got %s want %s", solved, ch.Combination)
// 		}
// 		ts = append(ts, time.Since(t))
// 	}
// 	b.Log(statistics(ts))
// }

// func statistics(ts []time.Duration) string {
// 	// Calculate statistics
// 	total := time.Duration(0)
// 	min := ts[0]
// 	max := ts[0]

// 	for _, t := range ts {
// 		total += t
// 		if t < min {
// 			min = t
// 		}
// 		if t > max {
// 			max = t
// 		}
// 	}
// 	avg := total / time.Duration(len(ts))

// 	// Calculate standard deviation
// 	var varianceSum float64
// 	for _, t := range ts {
// 		diff := float64(t-avg) / float64(time.Millisecond)
// 		varianceSum += diff * diff
// 	}
// 	stdDev := time.Duration(float64(time.Millisecond) * (varianceSum / float64(len(ts))))

// 	return fmt.Sprintf("Benchmark results: runs=%d avg=%v min=%v max=%v stdDev=%v", len(ts), avg, min, max, stdDev)
// }
