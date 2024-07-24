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
	"regexp"
	"strings"
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

var segmentpattern = regexp.MustCompile("^{(.*)}$")

func substitudeUrlWithParams(url string, urlparams map[string]string) (string, error) {
	segments := strings.Split(url, "/")
	for i := 0; i < len(segments); i++ {
		matches := segmentpattern.FindStringSubmatch(segments[i])
		if len(matches) == 2 {
			if value, ok := urlparams[matches[1]]; ok {
				segments[i] = value
			} else {
				return "", fmt.Errorf("url contains additional parameters that mapping type doesn't provide values for: %s", matches[1])
			}
		}
	}
	return strings.Join(segments, "/"), nil
}

func NewRequest[Request any](url, method string, params *Request) (*http.Request, error) {
	var err error
	var urlparams, bodyparams = separateParams(params)
	var buffer = bytes.NewBuffer([]byte{})
	if bodyparams != nil {
		err = json.NewEncoder(buffer).Encode(bodyparams)
		if err != nil {
			return nil, fmt.Errorf("serializing the body: %w", err)
		}
	}
	if urlparams != nil {
		url, err = substitudeUrlWithParams(url, urlparams)
		if err != nil {
			return nil, fmt.Errorf("substituting url parameters: %w", err)
		}
	}
	var r *http.Request
	r, err = http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, fmt.Errorf("creating request object: %w", err)
	}
	if bodyparams != nil {
		r.Header.Set("Content-Type", mime.TypeByExtension("json"))
		r.Header.Set("Content-Length", fmt.Sprintf("%d", buffer.Len()))
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
	if rs.StatusCode != 200 {
		return fmt.Errorf("server responded with status code: %d", rs.StatusCode)
	}
	if isBodyNeeded(bs) {
		err = parseJsonResponse[Response](rs, bs)
		if err != nil {
			return fmt.Errorf("binding response: %w", err)
		}
	}
	return nil
}
