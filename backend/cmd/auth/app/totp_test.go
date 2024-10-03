package app

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"testing"
)

func TestRequestTotp(t *testing.T) {
	a := New()

	p, err := a.GenerateTotpKey(GenerateTotpKeyRequest{
		User: "",
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act, GenerateTotpKey: %w", err))
	}

	err = os.MkdirAll("testresults", fs.ModePerm)
	if err != nil {
		t.Fatal(fmt.Errorf("MkdirAll: %w", err))
	}

	f, err := os.Create("testresults/qr.json")
	if err != nil {
		t.Fatal(fmt.Errorf("Create: %w", err))
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(p)
	if err != nil {
		t.Fatal(fmt.Errorf("encode: %w", err))
	}
}
