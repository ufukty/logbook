package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type MockStringType string

func (m *MockStringType) Set(s string) error {
	*m = MockStringType(s)
	return nil
}

type MockIntType int

func (m *MockIntType) Set(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*m = MockIntType(val)
	return nil
}

type TestRequest struct {
	SessionToken Cookie[MockIntType] `cookie:"session_token"`
	UserID       MockIntType         `route:"user_id"`
	Name         MockStringType      `json:"name"`
}

func compose(url string, path map[string]string, cookies map[string]string, body any) (*http.Request, error) {
	b := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		return nil, fmt.Errorf("json encoder: %w", err)
	}

	for k, v := range path {
		url = strings.Replace(url, k, v, 1)
	}

	req := httptest.NewRequest(http.MethodPost, url, b)

	for name, value := range cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	return req, nil
}

func TestValidRequest(t *testing.T) {
	var (
		path    = map[string]string{"{user_id}": "42"}
		cookies = map[string]string{"session_token": "123"}
		body    = map[string]string{"name": "John"}
	)
	r, err := compose("/{user_id}", path, cookies, body)
	if err != nil {
		t.Fatal(fmt.Errorf("compose: %w", err))
	}
	w := httptest.NewRecorder()

	m := http.NewServeMux()
	m.HandleFunc("/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		bq := &TestRequest{}

		err = ParseRequest(httptest.NewRecorder(), r, bq)
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
	})

	m.ServeHTTP(w, r)
}

func TestMissingCookie(t *testing.T) {
	var (
		body = map[string]string{"name": "John"}
		path = map[string]string{"user_id": "42"}
	)
	r, err := compose("/{user_id}", path, nil, body)
	if err != nil {
		t.Fatal(fmt.Errorf("compose: %w", err))
	}
	w := httptest.NewRecorder()

	m := http.NewServeMux()
	m.HandleFunc("/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		err = ParseRequest(httptest.NewRecorder(), r, &TestRequest{})
		if err == nil {
			t.Fatal("expected error due to missing cookie, got nil")
		}
		fmt.Println(err)
	})

	m.ServeHTTP(w, r)
}

func TestMissingURLFragment(t *testing.T) {
	var (
		cookies = map[string]string{"session_token": "123"}
		body    = map[string]string{"name": "John"}
	)

	r, err := compose("/", nil, cookies, body)
	if err != nil {
		t.Fatal(fmt.Errorf("compose: %w", err))
	}

	err = ParseRequest(httptest.NewRecorder(), r, &TestRequest{})
	if err == nil {
		t.Fatal("expected error due to missing URL fragment, got nil")
	}
}

func TestInvalidBody(t *testing.T) {
	var (
		url     = "/{user_id}"
		path    = map[string]string{"{user_id}": "42"}
		cookies = map[string]string{"session_token": "123"}
		body    = "{}"
	)

	r, err := compose(url, path, cookies, body)
	if err != nil {
		t.Fatal(fmt.Errorf("compose: %w", err))
	}
	w := httptest.NewRecorder()

	m := http.NewServeMux()
	m.HandleFunc("/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		err = ParseRequest(httptest.NewRecorder(), r, &TestRequest{})
		if err == nil {
			t.Fatal("expected error due to invalid JSON body, got nil")
		}
		fmt.Println(err)
	})

	m.ServeHTTP(w, r)
}

type TestRequest2 struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func TestParseRequestForEmptyBodies(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    *TestRequest2
		limit   int64
		wantErr bool
	}{
		{
			name: "Body within limit",
			body: `{"field1": "value1", "field2": "value2"}`,
			want: &TestRequest2{
				Field1: "value1",
				Field2: "value2",
			},
			limit:   bodyBufferLimit,
			wantErr: false,
		},
		{
			name:    "Body exceeds limit",
			body:    strings.Repeat("a", int(bodyBufferLimit+1)),
			want:    nil,
			limit:   bodyBufferLimit,
			wantErr: true,
		},
		{
			name:    "Empty body",
			body:    "",
			want:    nil,
			limit:   bodyBufferLimit,
			wantErr: true,
		},
		{
			name:    "Invalid JSON",
			body:    `{"field1": "value1", "field2":}`,
			want:    nil,
			limit:   bodyBufferLimit,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()

			var got TestRequest2
			err := ParseRequest(rec, req, &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (tt.wantErr) != (err != nil) {
				t.Errorf("ParseRequest() wantErr=%t, got %v", tt.wantErr, err)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(&got, tt.want) {
				t.Errorf("ParseRequest() = %v, want %v", &got, tt.want)
			}
		})
	}
}
