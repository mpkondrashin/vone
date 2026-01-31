/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_submit_urls.go - send urls for analysis
*/

package vone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type (
	SubmitURLsToSandboxURL struct {
		URL string `json:"url"`
	}

	SubmitURLsToSandboxRequest []SubmitURLsToSandboxURL

	BodyStruct struct {
		URL    string     `json:"url"`
		ID     string     `json:"id"`
		Digest Digest     `json:"digest"`
		Error  Innererror `json:"error"`
	}

	SubmitURLsToSandboxStruct struct {
		Status  int `json:"status"`
		Headers []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"headers,omitempty"`
		Body BodyStruct `json:"body,omitempty"`
	}

	SandboxSubmitURLsToSandboxResponse []SubmitURLsToSandboxStruct
)

type sandboxSubmitURLsRequest struct {
	baseRequest
	request         SubmitURLsToSandboxRequest
	response        SandboxSubmitURLsToSandboxResponse
	responseHeaders SandboxSubmitFileResponseHeaders
}

func (v *VOne) SandboxSubmitURLs() *sandboxSubmitURLsRequest {
	f := &sandboxSubmitURLsRequest{}
	f.baseRequest.init(v)
	return f
}

func (s *sandboxSubmitURLsRequest) AddURL(url string) *sandboxSubmitURLsRequest {
	s.request = append(s.request, SubmitURLsToSandboxURL{URL: url})
	return s
}

func (f *sandboxSubmitURLsRequest) AddURLs(urls []string) *sandboxSubmitURLsRequest {
	for _, url := range urls {
		f.AddURL(url)
	}
	return f
}

func (f *sandboxSubmitURLsRequest) Do(ctx context.Context) (SandboxSubmitURLsToSandboxResponse, *SandboxSubmitFileResponseHeaders, error) {
	if err := f.checkUsed(); err != nil {
		return nil, nil, fmt.Errorf("submit URLs: %w", err)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, nil, err
	}
	return f.response, &f.responseHeaders, nil
}

func (f *sandboxSubmitURLsRequest) method() string {
	return methodPost
}

func (s *sandboxSubmitURLsRequest) url() string {
	return "/v3.0/sandbox/urls/analyze"
}

func (f *sandboxSubmitURLsRequest) requestBody() io.Reader {
	data, _ := json.Marshal(f.request)
	return bytes.NewBuffer(data)
}

func (f *sandboxSubmitURLsRequest) responseStruct() any {
	return &f.response
}

func (f *sandboxSubmitURLsRequest) responseHeader() any {
	return &f.responseHeaders
}
