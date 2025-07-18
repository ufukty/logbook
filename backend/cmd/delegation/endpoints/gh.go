// Code generated by gohandlers v0.37.0. DO NOT EDIT.

package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ufukty/gohandlers/pkg/gohandlers"
)

func (pu *Public) ListHandlers() map[string]gohandlers.HandlerInfo {
	return map[string]gohandlers.HandlerInfo{
		"DelegateObjective":   {Method: "POST", Path: "/delegate-objective", Ref: pu.DelegateObjective},
		"ListDelegationChain": {Method: "GET", Path: "/list-delegation-chain/{subject}", Ref: pu.ListDelegationChain},
		"RemoveDelegation":    {Method: "POST", Path: "/remove-delegation", Ref: pu.RemoveDelegation},
	}
}

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

func (bq DelegateObjectiveRequest) Build(host string) (*http.Request, error) {
	uri := "/delegate-objective"
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

func (bq *DelegateObjectiveRequest) Parse(rq *http.Request) error {
	if !strings.HasPrefix(rq.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("invalid content type for request: %s", rq.Header.Get("Content-Type"))
	}
	if err := json.NewDecoder(rq.Body).Decode(bq); err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}
	return nil
}

func (bq DelegateObjectiveRequest) Validate() (issues map[string]any) {
	issues = map[string]any{}
	if issue := bq.Delegator.Validate(); issue != nil {
		issues["delegator"] = issue
	}
	if issue := bq.Delegee.Validate(); issue != nil {
		issues["delegee"] = issue
	}
	if issue := bq.Objective.Validate(); issue != nil {
		issues["objective"] = issue
	}
	return
}

func (bs DelegateObjectiveResponse) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(bs); err != nil {
		return fmt.Errorf("encoding the body: %w", err)
	}
	return nil
}

func (bs *DelegateObjectiveResponse) Parse(rs *http.Response) error {
	if !strings.HasPrefix(rs.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("invalid content type for request: %s", rs.Header.Get("Content-Type"))
	}
	if err := json.NewDecoder(rs.Body).Decode(bs); err != nil {
		return fmt.Errorf("decoding the body: %w", err)
	}
	return nil
}

func (bq ListDelegationChainRequest) Build(host string) (*http.Request, error) {
	uri := "/list-delegation-chain/{subject}"
	encoded, err := bq.Subject.ToRoute()
	if err != nil {
		return nil, fmt.Errorf("ListDelegationChainRequest.Subject.ToRoute: %w", err)
	}
	uri = strings.Replace(uri, "{subject}", encoded, 1)
	r, err := http.NewRequest("GET", join(host, uri), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	return r, nil
}

func (bq *ListDelegationChainRequest) Parse(rq *http.Request) error {
	if err := bq.Subject.FromRoute(rq.PathValue("subject")); err != nil {
		return fmt.Errorf("ListDelegationChainRequest.Subject.FromRoute: %w", err)
	}
	return nil
}

func (bq ListDelegationChainRequest) Validate() (issues map[string]any) {
	issues = map[string]any{}
	if issue := bq.Subject.Validate(); issue != nil {
		issues["subject"] = issue
	}
	return
}

func (bs ListDelegationChainResponse) Write(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)
	return nil
}

func (bs *ListDelegationChainResponse) Parse(rs *http.Response) error {
	if !strings.HasPrefix(rs.Header.Get("Content-Type"), "") {
		return fmt.Errorf("invalid content type for request: %s", rs.Header.Get("Content-Type"))
	}
	return nil
}

func (bq RemoveDelegationRequest) Build(host string) (*http.Request, error) {
	uri := "/remove-delegation"
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

func (bq *RemoveDelegationRequest) Parse(rq *http.Request) error {
	if !strings.HasPrefix(rq.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("invalid content type for request: %s", rq.Header.Get("Content-Type"))
	}
	if err := json.NewDecoder(rq.Body).Decode(bq); err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}
	return nil
}

func (bq RemoveDelegationRequest) Validate() (issues map[string]any) {
	issues = map[string]any{}
	if issue := bq.Delid.Validate(); issue != nil {
		issues["delid"] = issue
	}
	return
}

func (bs RemoveDelegationResponse) Write(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)
	return nil
}

func (bs *RemoveDelegationResponse) Parse(rs *http.Response) error {
	if !strings.HasPrefix(rs.Header.Get("Content-Type"), "") {
		return fmt.Errorf("invalid content type for request: %s", rs.Header.Get("Content-Type"))
	}
	return nil
}
