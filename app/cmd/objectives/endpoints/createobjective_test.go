package endpoints

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateObjective(t *testing.T) {
	ep, err := getTestDependencies()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, creating Endpoints: %w", err))
	}

	r, err := http.NewRequest("POST", "", strings.NewReader(`{
		"parent": {
			"oid": "00000000-0000-0000-0000-000000000000",
			"vid": "00000000-0000-0000-0000-000000000000"
		},
		"content": "Lorem ipsum dolor sit amet consectetur adipscing elit"
	}`))
	if err != nil {
		t.Fatal(fmt.Errorf("prep, creating http request: %w", err))
	}
	// req.Header.Set("Authentication", ")"

	w := httptest.NewRecorder()

	ep.CreateTask(w, r)

	if w.Code != http.StatusOK {
		t.Fatal(fmt.Sprintf("got http error code %v", w.Code))
	}

	// Check the response body
	expected := "SessionId: mock-session-id, Content: test content"
	if w.Body.String() != expected {
		t.Fatal(fmt.Sprintf("handler returned unexpected body: got %v want %v", w.Body.String(), expected))
	}
}
