// this utility is used to create test cases for frontend implementation in TypeScript
package main

import (
	"encoding/json"
	"fmt"
	"logbook/internal/crypto/challenge"
	"os"
)

func create(difficulty int) error {
	cs, err := challenge.CreateBatch(difficulty, 100)
	if err != nil {
		return fmt.Errorf("CreateBatch: %w", err)
	}
	f, err := os.Create(fmt.Sprintf("d%d.json", difficulty))
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(cs)
	if err != nil {
		return fmt.Errorf("printing json: %w", err)
	}
	return nil
}

func Main() error {
	for d := 10; d <= 60; d += 10 {
		err := create(d)
		if err != nil {
			return fmt.Errorf("d%d: %w", d, err)
		}
	}
	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
