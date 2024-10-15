package api

import (
	"fmt"
	"testing"
)

func TestParentRefValues(t *testing.T) {
	c, err := ReadConfig("../../api.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("prep: %w", err))
	}

	if &c.Public.Services != c.Public.Services.Account.Parent {
		t.Error("assert")
	}
}

func TestByGateway(t *testing.T) {
	c, err := ReadConfig("../../api.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("prep: %w", err))
	}

	tcs := map[Addressable]string{
		c.Public:                  "/api",
		c.Public.Services.Account: "/api/account",
		c.Public.Services.Account.Endpoints.CreateAccount: "/api/account/account",
		c.Public.Services.Document:                        "/api/document",
		c.Public.Services.Document.Endpoints.List:         "/api/document/list/{root}",
	}

	for in, want := range tcs {
		t.Run(want, func(t *testing.T) {
			got := ByGateway(in)
			if want != got {
				t.Errorf("assert\nwant: %q\ngot : %q", want, got)
			}
		})
	}

}
