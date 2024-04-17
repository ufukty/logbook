// rq -> http.Request
// rs -> http.Response
// bq -> Request
// bs -> Response
package reqs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"logbook/config/api"
	"mime"
	"net/http"
	"path/filepath"
	"reflect"

	"github.com/gorilla/mux"
)

func separateParams(in any) (map[string]string, map[string]any) {
	var t = reflect.TypeOf(in)
	var v = reflect.ValueOf(in)
	var fields = v.NumField()
	var url = map[string]string{}
	var body = map[string]any{}
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

func newRequest(url, method string, params any) (*http.Request, error) {
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

func WriteJsonResponse(bs any, rsw http.ResponseWriter) error {
	rsw.Header().Set("Content-Type", mime.TypeByExtension(".json"))
	rsw.WriteHeader(http.StatusOK)
	err := json.NewEncoder(rsw).Encode(bs)
	if err != nil {
		return fmt.Errorf("serializing the body: %w", err)
	}
	return nil
}

func parseJsonResponse[Response any](rs *http.Response) (bs *Response, err error) {
	rs.Header.Get("Content-Type")
	bs = new(Response)
	err = json.NewDecoder(rs.Body).Decode(bs)
	if err != nil {
		return nil, fmt.Errorf("parsing the response body: %w", err)
	}
	return
}

func Send[Request any, Response any](url, method string, bq *Request) (*Response, error) {
	var (
		rq  *http.Request
		rs  *http.Response
		bs  = new(Response)
		err error
	)
	rq, err = newRequest(url, method, bq)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	rs, err = http.DefaultClient.Do(rq)
	if err != nil {
		return nil, fmt.Errorf("sending the request: %w", err)
	}
	bs, err = parseJsonResponse[Response](rs)
	if err != nil {
		return nil, fmt.Errorf("binding response: %w", err)
	}
	return bs, nil
}

func SendTo[Request any, Response any](path string, dst api.Endpoint, bq *Request) (*Response, error) {
	return Send[Request, Response](filepath.Join(path, string(dst.Method)), dst.Method, bq)
}
