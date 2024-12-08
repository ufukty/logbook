package challenge

import (
	"fmt"
)

const (
	BL  = 40   // batch id length
	ML  = 3    // mask length
	OL  = 20   // 'original' length
	PHL = 1000 // pre-hash length
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
	o2, err := randstring(alphabet, OL-ML)
	if err != nil {
		return "", fmt.Errorf("randstring for clear part: %w", err)
	}
	return o1 + o2, nil
}

func CreateChallenge(bid string, difficulty int) (Challenge, error) {
	r, err := createOriginal(difficulty)
	if err != nil {
		return Challenge{}, fmt.Errorf("randstring: %w", err)
	}
	c := Challenge{
		Masked:     mask(r),
		Hashed:     hash(r + bid + pseudo[:PHL-BL-OL]),
		Original:   r,
		Difficulty: difficulty,
	}
	return c, nil
}

var (
	ErrMinDifficulty = fmt.Errorf("difficulty needs to be bigger than 1")
	ErrMaxDifficulty = fmt.Errorf("difficulty needs to be smaller than the alphabet size")
)

type Batch struct {
	BatchId    string      `json:"bid"`
	Challenges []Challenge `json:"challenges"`
}

func CreateBatch(difficulty, CPB int) (*Batch, error) {
	if difficulty < 2 {
		return nil, ErrMinDifficulty
	}
	if len(alphabet) <= difficulty {
		return nil, ErrMaxDifficulty
	}
	bid, err := randstring(alphabet, BL)
	if err != nil {
		return nil, fmt.Errorf("creating batch id: %w", err)
	}
	b := &Batch{
		BatchId: bid,
	}
	for range CPB {
		c, err := CreateChallenge(bid, difficulty)
		if err != nil {
			return nil, fmt.Errorf("CreateChallenge: %w", err)
		}
		b.Challenges = append(b.Challenges, c)
	}
	return b, nil
}
