package data

import (
	"encoding/json"
	"fmt"
	"os"
)

type Objective struct {
	Content  string      `json:"content"`
	Children []Objective `json:"children"`
}

func LoadTestData() ([]Objective, error) {
	f, err := os.Open("testdata/company.json")
	if err != nil {
		return nil, fmt.Errorf("opening: %w", err)
	}
	defer f.Close()

	os := &[]Objective{}
	err = json.NewDecoder(f).Decode(os)
	if err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	return *os, nil
}
