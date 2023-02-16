package vone

import (
	"bytes"
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

	SubmitURLsToSandboxResponse []SubmitURLsToSandboxStruct
)

type SubmitURLsToSandboxFunc struct {
	BaseFunc
	Request  SubmitURLsToSandboxRequest
	Response SubmitURLsToSandboxResponse
}

var _ Func = &SubmitURLsToSandboxFunc{}

/*
func (v *VOne) SubmitURLsToSandbox(urls []string) (SubmitURLsToSandboxResponse, error) {
	f, err := NewSubmitURLsToSandboxFunc()
	if err != nil {
		return nil, fmt.Errorf("NewSubmitURLsToSandboxFunc: %w", err)
	}
	for _, url := range urls {
		f.AddURL(url)
	}
	if err := v.Call(f); err != nil {
		return nil, fmt.Errorf("Call: %w", err)
	}
	return f.Response, nil
}
*/
/*
func (v *VOne) SubmitURLsToSandboxData(urls []string) (SubmitURLsToSandboxDataResponse, error) {
	var data SubmitURLsToSandboxRequest
	for _, url := range urls {
		data = append(data, SubmitURLsToSandboxURL{URL: url})
	}

	decoder := json.NewDecoder(body)
	var respData SubmitURLsToSandboxDataResponse
	if err := decoder.Decode(&respData); err != nil && err != io.EOF {
		return nil, fmt.Errorf("response error: %w", err)
	}
	return respData, nil
}
*/

func (v *VOne) NewSubmitURLsToSandbox() *SubmitURLsToSandboxFunc {
	f := &SubmitURLsToSandboxFunc{}
	f.BaseFunc.Init(v)
	return f
}

func (s *SubmitURLsToSandboxFunc) AddURL(url string) *SubmitURLsToSandboxFunc {
	s.Request = append(s.Request, SubmitURLsToSandboxURL{URL: url})
	return s
}

func (f *SubmitURLsToSandboxFunc) AddURLs(urls []string) *SubmitURLsToSandboxFunc {
	for _, url := range urls {
		f.AddURL(url)
	}
	return f
}

func (f *SubmitURLsToSandboxFunc) Do() (SubmitURLsToSandboxResponse, error) {
	if err := f.vone.Call(f); err != nil {
		return nil, err
	}
	return f.Response, nil
}

func (f *SubmitURLsToSandboxFunc) Method() string {
	return "POST"
}

func (s *SubmitURLsToSandboxFunc) URL() string {
	return "/v3.0/sandbox/urls/analyze"
}

func (f *SubmitURLsToSandboxFunc) RequestBody() io.Reader {
	data, err := json.Marshal(f.Request)
	if err != nil {

	}
	return bytes.NewBuffer(data)
}

func (f *SubmitURLsToSandboxFunc) ResponseStruct() any {
	return &f.Response
}
