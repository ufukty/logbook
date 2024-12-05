// this utility is used to create test cases for frontend implementation in TypeScript
package main

import (
	"encoding/json"
	"fmt"
	"logbook/internal/crypto/challenge"
	"os"
)

func Main() error {
	cs, err := challenge.CreateBatch(30, 100)
	if err != nil {
		return fmt.Errorf("CreateBatch: %w", err)
	}

	err = json.NewEncoder(os.Stdout).Encode(cs)
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
