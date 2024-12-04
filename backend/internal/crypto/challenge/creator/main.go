// this utility is used to create test cases for frontend implementation in TypeScript
package main

import (
	"encoding/json"
	"fmt"
	"logbook/internal/crypto/challenge"
	"os"
)

func Main() error {
	cs := []challenge.Challange{}

	for i := 0; i < 100; i++ {
		c, err := challenge.NewChallenge(500, 3)
		if err != nil {
			return fmt.Errorf("creating a challenge: %w", err)
		}
		cs = append(cs, c)
	}

	err := json.NewEncoder(os.Stdout).Encode(cs)
	if err != nil {
		return fmt.Errorf("printing json: %w", err)
	}

	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
