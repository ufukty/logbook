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

func TestPathFromInternet(t *testing.T) {
	c, err := ReadConfig("../../api.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("prep: %w", err))
	}

	tcs := map[Addressable]string{
		c.Public:                  "/api/public",
		c.Public.Services.Account: "/api/public/account",
		c.Public.Services.Account.Endpoints.CreateUser: "/api/public/account/user",
		c.Public.Services.Document:                     "/api/public/document",
		c.Public.Services.Document.Endpoints.List:      "/api/public/document/list/{root}",
	}

	for in, want := range tcs {
		t.Run(want, func(t *testing.T) {
			got := PathFromInternet(in)
			if want != got {
				t.Errorf("assert\nwant: %q\ngot : %q", want, got)
			}
		})
	}

}
