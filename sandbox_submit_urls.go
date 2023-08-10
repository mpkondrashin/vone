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

type SandboxSubmitURLsToSandboxFunc struct {
	baseFunc
	Request         SubmitURLsToSandboxRequest
	Response        SandboxSubmitURLsToSandboxResponse
	ResponseHeaders SandboxSubmitFileResponseHeaders
}

func (v *VOne) SandboxSubmitURLs() *SandboxSubmitURLsToSandboxFunc {
	f := &SandboxSubmitURLsToSandboxFunc{}
	f.baseFunc.init(v)
	return f
}

func (s *SandboxSubmitURLsToSandboxFunc) AddURL(url string) *SandboxSubmitURLsToSandboxFunc {
	s.Request = append(s.Request, SubmitURLsToSandboxURL{URL: url})
	return s
}

func (f *SandboxSubmitURLsToSandboxFunc) AddURLs(urls []string) *SandboxSubmitURLsToSandboxFunc {
	for _, url := range urls {
		f.AddURL(url)
	}
	return f
}

func (f *SandboxSubmitURLsToSandboxFunc) Do(ctx context.Context) (SandboxSubmitURLsToSandboxResponse, *SandboxSubmitFileResponseHeaders, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, nil, err
	}
	return f.Response, &f.ResponseHeaders, nil
}

func (f *SandboxSubmitURLsToSandboxFunc) method() string {
	return methodPost
}

func (s *SandboxSubmitURLsToSandboxFunc) url() string {
	return "/v3.0/sandbox/urls/analyze"
}

func (f *SandboxSubmitURLsToSandboxFunc) requestBody() io.Reader {
	data, _ := json.Marshal(f.Request)
	return bytes.NewBuffer(data)
}

func (f *SandboxSubmitURLsToSandboxFunc) responseStruct() any {
	return &f.Response
}

func (f *SandboxSubmitURLsToSandboxFunc) responseHeader() any {
	return &f.ResponseHeaders
}
