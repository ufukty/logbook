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

	tcs := map[addressable]string{
		c.Public:                  "/api/public",
		c.Public.Services.Account: "/api/public/account/public",
		c.Public.Services.Account.Endpoints.CreateUser: "/api/public/account/public/user",
		c.Public.Services.Document:                     "/api/public/document/public",
		c.Public.Services.Document.Endpoints.List:      "/api/public/document/public/list/{root}",
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
