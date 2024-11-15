package challenge

import (
	"testing"
)

func TestSolve(t *testing.T) {
	ch, err := newChallenge(1)
	if err != nil {
		t.Fatalf("failed to create challenge: %v", err)
	}

	solvedOriginal, err := Solve(ch.Que, ch.Hash)
	if err != nil {
		t.Fatalf("failed to solve challenge: %v", err)
	}

	if solvedOriginal != ch.Original {
		t.Errorf("solve failed: expected %s, got %s", ch.Original, solvedOriginal)
	} else {
		t.Logf("solve succeeded: %s", solvedOriginal)
	}
}
