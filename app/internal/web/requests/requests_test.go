package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func createTestRequest(body interface{}, cookies map[string]string, urlVars map[string]string) *http.Request {
	// Encode the body as JSON
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyBytes))

	// Add cookies
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	// Add URL variables using gorilla/mux
	req = mux.SetURLVars(req, urlVars)

	return req
}

func TestValidRequest(t *testing.T) {
	body := map[string]string{"name": "John"}
	cookies := map[string]string{
		"session_token": "123",
	}
	urlVars := map[string]string{
		"user_id": "42",
	}

	req := createTestRequest(body, cookies, urlVars)
	bq := &TestRequest{}
	err := ParseRequest(req, bq)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if bq.Name != MockStringType("John") {
		t.Errorf("expected name to be 'John', got %s", bq.Name)
	}

	if bq.UserID != MockIntType(42) {
		t.Errorf("expected user_id to be 42, got %d", bq.UserID)
	}

	if bq.SessionToken.Value != 123 {
		t.Errorf("expected session_token to be '123', got %d", bq.SessionToken.Value)
	}

	if bq.SessionToken.Cookie == nil {
		t.Error("expected session_token details to be set, got nil")
	}
}

func TestMissingCookie(t *testing.T) {
	body := map[string]string{"name": "John"}
	urlVars := map[string]string{
		"user_id": "42",
	}

	req := createTestRequest(body, nil, urlVars)
	bq := &TestRequest{}
	err := ParseRequest(req, bq)
	if err == nil {
		t.Fatal("expected error due to missing cookie, got nil")
	}
	fmt.Println(err)
}

func TestMissingURLFragment(t *testing.T) {
	body := map[string]string{"name": "John"}
	cookies := map[string]string{
		"session_token": "123",
	}

	req := createTestRequest(body, cookies, nil)
	bq := &TestRequest{}
	err := ParseRequest(req, bq)
	if err == nil {
		t.Fatal("expected error due to missing URL fragment, got nil")
	}
}

func TestInvalidBody(t *testing.T) {
	body := "invalid json"
	cookies := map[string]string{
		"session_token": "123",
	}
	urlVars := map[string]string{
		"user_id": "42",
	}

	req := createTestRequest(body, cookies, urlVars)
	bq := &TestRequest{}
	err := ParseRequest(req, bq)
	if err == nil {
		t.Fatal("expected error due to invalid JSON body, got nil")
	}
}

func TestVerificationError(t *testing.T) {
	body := map[string]string{"name": "John"}
	cookies := map[string]string{
		"session_token": "123",
	}
	urlVars := map[string]string{
		"user_id": "42",
	}

	req := createTestRequest(body, cookies, urlVars)
	bq := &TestRequest{}
	err := ParseRequest(req, bq)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Intentionally causing a verification error
	if bq.Name != MockStringType("Doe") {
		t.Errorf("expected name to be 'Doe', got %s", bq.Name)
	}
}
