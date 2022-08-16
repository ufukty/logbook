package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/pkg/errors"
)

func PrepareJSONRequest(method string, target string, requestParams interface{}) *http.Request {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(requestParams)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Could not convert struct to json formed response body"))
	}
	r := httptest.NewRequest(method, target, &buf)
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	return r
}

func DecodeJSONRequest(param interface{}, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		return errors.Wrap(err, "decodeContentTypeJSON")
	}
	return nil
}

func DecodeJSONResponse(param interface{}, w *httptest.ResponseRecorder) error {
	res := w.Result()
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&param)
	if err != nil {
		return errors.Wrap(err, "decodeContentTypeJSON")
	}
	return nil
}
