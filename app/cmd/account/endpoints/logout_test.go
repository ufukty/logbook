package endpoints

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoutWithInvalidToken(t *testing.T) {
	ep, err := getTestDependencies()
	if err != nil {
		t.Fatal(fmt.Errorf("getting dependencies: %w", err))
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "invalid token",
	})
	w := httptest.NewRecorder()

	ep.Logout(w, r)

	if w.Code == 200 {
		t.Error("Failure expected.")
	}
}
