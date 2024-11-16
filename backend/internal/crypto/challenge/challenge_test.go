package challenge

import (
	"fmt"
	"testing"
)

func TestCombinate(t *testing.T) {
	type tc struct {
		got  string
		want string
	}

	tcs := []tc{
		{combinate(1, 0), "A"},
		{combinate(2, 0), "AA"},
		{combinate(2, 2), "AC"},
		{combinate(4, 4), "AAAE"},
		{combinate(1, 32), "A"},
		{combinate(1, 33), "B"},
	}

	for _, tc := range tcs {
		if tc.got != tc.want {
			t.Errorf("want %s got %s", tc.want, tc.got)
		}
	}
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
		{20, 5},
		{20, 6},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%d-%d", tc.input_l, tc.input_n), func(t *testing.T) {
			ch, err := NewChallenge(tc.input_l, tc.input_n)
			if err != nil {
				t.Fatalf("failed to create challenge: %v", err)
			}
			t.Log(ch)
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

func BenchmarkSolve(b *testing.B) {
	n := 0
	for i := 0; i < b.N; i++ {
		n = b.N
		ch, err := NewChallenge(20, 1)
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
	}
	fmt.Println(">", n)
}
