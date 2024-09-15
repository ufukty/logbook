package requests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

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
