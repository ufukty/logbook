package challenge

import (
	"fmt"
)

const (
	ML = 3
	IL = 20
)

type Challenge struct {
	Difficulty int    `json:"D"`
	Masked     string `json:"M"`
	Hashed     string `json:"H"`
	Original   string `json:"O"`
}

func (c Challenge) String() string {
	return fmt.Sprintf("(Original: %s) (Masked: %s) (Hash: %s) \n", c.Original, c.Masked, c.Hashed)
}

func mask(o string) string {
	return string(o[ML:])
}

func createOriginal(difficulty int) (string, error) {
	o1, err := randstring(alphabet[:difficulty], ML)
	if err != nil {
		return "", fmt.Errorf("randstring for masked part: %w", err)
	}
	o2, err := randstring(alphabet, IL-ML)
	if err != nil {
		return "", fmt.Errorf("randstring for clear part: %w", err)
	}
	return o1 + o2, nil
}

func CreateChallenge(difficulty int) (Challenge, error) {
	o, err := createOriginal(difficulty)
	if err != nil {
		return Challenge{}, fmt.Errorf("randstring: %w", err)
	}
	c := Challenge{
		Masked:     mask(o),
		Hashed:     hash(o),
		Original:   o,
		Difficulty: difficulty,
	}
	return c, nil
}

var (
	ErrMinDifficulty = fmt.Errorf("difficulty needs to be bigger than 1")
	ErrMaxDifficulty = fmt.Errorf("difficulty needs to be smaller than the alphabet size")
)

func CreateBatch(difficulty, CPB int) ([]Challenge, error) {
	if difficulty < 2 {
		return nil, ErrMinDifficulty
	}
	if len(alphabet) <= difficulty {
		return nil, ErrMaxDifficulty
	}
	cs := []Challenge{}
	for range CPB {
		c, err := CreateChallenge(difficulty)
		if err != nil {
			return nil, fmt.Errorf("CreateChallenge: %w", err)
		}
		cs = append(cs, c)
	}
	return cs, nil
}
