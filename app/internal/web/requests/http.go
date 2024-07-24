// rq -> http.Request
// rs -> http.Response
// bq -> Request
// bs -> Response
package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

func separateParams(in any) (map[string]string, map[string]any) {
	url := map[string]string{}
	body := map[string]any{}

	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)	
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
	}
	fields := v.NumField()
	var tags reflect.StructTag

	for i := 0; i < fields; i++ {
		tags = t.Field(i).Tag
		if key, exists := tags.Lookup("url"); exists {
			url[key] = v.Field(i).String()
		} else if key, exists := tags.Lookup("json"); exists {
			body[key] = v.Field(i).Interface()
		}
	}
	return url, body
}

func NewRequest[Request any](url, method string, params *Request) (*http.Request, error) {
	var err error
	var urlParams, bodyParams = separateParams(params)
	var buffer = bytes.NewBuffer([]byte{})
	if bodyParams != nil {
		err = json.NewEncoder(buffer).Encode(bodyParams)
		if err != nil {
			return nil, fmt.Errorf("serializing the body: %w", err)
		}
	}
	var r *http.Request
	r, err = http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, fmt.Errorf("creating request object: %w", err)
	}
	if bodyParams != nil {
		r.Header.Set("Content-Type", mime.TypeByExtension("json"))
		r.Header.Set("Content-Length", fmt.Sprintf("%d", buffer.Len()))
	}
	if urlParams != nil {
		r = mux.SetURLVars(r, urlParams)
	}
	return r, nil
}

func WriteJsonResponse[Response any](bs Response, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(bs)
	if err != nil {
		return fmt.Errorf("serializing the body: %w", err)
	}
	return nil
}

func parseJsonResponse[Response any](rs *http.Response, bs *Response) error {
	if ct := rs.Header.Get("Content-Type"); ct != mime.TypeByExtension(".json") {
		return fmt.Errorf("unsupported Content-Type: %s", ct)
	}
	err := json.NewDecoder(rs.Body).Decode(bs)
	if err != nil {
		return fmt.Errorf("parsing the response body: %w", err)
	}
	return nil
}

func SendRaw[Request any](url, method string, bq *Request) (*http.Response, error) {
	rq, err := NewRequest(url, method, bq)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	rs, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, fmt.Errorf("sending the request: %w", err)
	}
	return rs, nil
}

func Send[Request any, Response any](url, method string, bq *Request, bs *Response) error {
	rq, err := NewRequest(url, method, bq)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	rs, err := http.DefaultClient.Do(rq)
	if err != nil {
		return fmt.Errorf("sending the request: %w", err)
	}
	err = parseJsonResponse[Response](rs, bs)
	if err != nil {
		return fmt.Errorf("binding response: %w", err)
	}
	return nil
}
