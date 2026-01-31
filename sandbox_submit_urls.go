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
	"io"
)

type (
	SubmitURLsToSandboxURL struct {
		URL string `json:"url"`
	}

	SubmitURLsToSandboxRequest []SubmitURLsToSandboxURL

	BodyStruct struct {
		URL    string `json:"url"`
		ID     string `json:"id"`
		Digest struct {
			MD5    string `json:"md5"`
			SHA1   string `json:"sha1"`
			SHA256 string `json:"sha256"`
		} `json:"digest"`
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
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

type sandboxSubmitURLsFunc struct {
	baseFunc
	request         SubmitURLsToSandboxRequest
	response        SandboxSubmitURLsToSandboxResponse
	responseHeaders SandboxSubmitFileResponseHeaders
}

func (v *VOne) SandboxSubmitURLs() *sandboxSubmitURLsFunc {
	f := &sandboxSubmitURLsFunc{}
	f.baseFunc.init(v)
	return f
}

func (s *sandboxSubmitURLsFunc) AddURL(url string) *sandboxSubmitURLsFunc {
	s.request = append(s.request, SubmitURLsToSandboxURL{URL: url})
	return s
}

func (f *sandboxSubmitURLsFunc) AddURLs(urls []string) *sandboxSubmitURLsFunc {
	for _, url := range urls {
		f.AddURL(url)
	}
	return f
}

func (f *sandboxSubmitURLsFunc) Do(ctx context.Context) (SandboxSubmitURLsToSandboxResponse, *SandboxSubmitFileResponseHeaders, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, nil, err
	}
	return f.response, &f.responseHeaders, nil
}

func (f *sandboxSubmitURLsFunc) method() string {
	return methodPost
}

func (s *sandboxSubmitURLsFunc) url() string {
	return "/v3.0/sandbox/urls/analyze"
}

func (f *sandboxSubmitURLsFunc) requestBody() io.Reader {
	data, _ := json.Marshal(f.request)
	return bytes.NewBuffer(data)
}

func (f *sandboxSubmitURLsFunc) responseStruct() any {
	return &f.response
}

func (f *sandboxSubmitURLsFunc) responseHeader() any {
	return &f.responseHeaders
}
