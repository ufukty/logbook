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

	b := strings.NewReader(`
	{
		"parent": {
			"oid": "00000000-0000-0000-0000-000000000000",
			"vid": "00000000-0000-0000-0000-000000000000"
		},
		"content": "Lorem ipsum dolor sit amet consectetur adipscing elit"
	}
	`)
	req, err := http.NewRequest("POST", "", b)
	if err != nil {
		t.Fatal(fmt.Errorf("prep, creating http request: %w", err))
	}
	// req.Header.Set("Authentication", ")"

	rr := httptest.NewRecorder()

	ep.CreateTask(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatal(fmt.Sprintf("got http error code %v", status))
	}

	// Check the response body
	expected := "SessionId: mock-session-id, Content: test content"
	if rr.Body.String() != expected {
		t.Fatal(fmt.Sprintf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected))
	}
}
