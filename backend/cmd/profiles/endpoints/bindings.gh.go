// Code generated by gohandlers v0.27.3. DO NOT EDIT.

package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func join(segments ...string) string {
	url := ""
	for i, segment := range segments {
		if i != 0 && !strings.HasPrefix(segment, "/") {
			url += "/"
		}
		url += segment
	}
	return url
}

func (bq CreateProfileRequest) Build(host string) (*http.Request, error) {
	uri := "/create-profile"
	body := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(body).Encode(bq); err != nil {
		return nil, fmt.Errorf("json.Encoder.Encode: %w", err)
	}
	r, err := http.NewRequest("POST", join(host, uri), body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Content-Length", fmt.Sprintf("%d", body.Len()))
	return r, nil
}

func (bq *CreateProfileRequest) Parse(rq *http.Request) error {
	if !strings.HasPrefix(rq.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("invalid content type for request: %s", rq.Header.Get("Content-Type"))
	}
	if err := json.NewDecoder(rq.Body).Decode(bq); err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}
	return nil
}
