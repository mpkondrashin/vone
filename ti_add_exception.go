/*
	Trend Micro Vision One API SDK
	(c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	it_add_exception.go - upload file for analysis
*/

package vone

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

// TIAddException - Add exceptions payload struct
type (
	TIAddException []TIException

	TIException struct {
		URL               string `json:"url,omitempty"`
		Description       string `json:"description,omitempty"`
		Domain            string `json:"domain,omitempty"`
		IP                string `json:"ip,omitempty"`
		SenderMailAddress string `json:"senderMailAddress,omitempty"`
		FileSha1          string `json:"fileSha1,omitempty"`
		FileSha256        string `json:"fileSha256,omitempty"`
	}

	TIAddExceptionResponce []struct {
		Status int `json:"status"`
		Body   struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		} `json:"body,omitempty"`
	}
)

// TIAddExceptionFunc - function to add exceptions
type TIAddExceptionFunc struct {
	baseFunc
	//filePath            string
	Request TIAddException
	//formDataContentType string
	Response TIAddExceptionResponce
	// ResponseHeaders     SandboxSubmitFileResponseHeaders
}

var _ vOneFunc = &TIAddExceptionFunc{}

// SandboxSubmitFile - return new XXXXX
func (v *VOne) AddExceptions() *TIAddExceptionFunc {
	f := &TIAddExceptionFunc{}
	f.baseFunc.init(v)
	return f
}

func (f *TIAddExceptionFunc) Add(exception TIException) *TIAddExceptionFunc {
	f.Request = append(f.Request, exception)
	return f
}

func (f *TIAddExceptionFunc) AddURL(url string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		URL:         url,
		Description: description,
	})
	return f
}

func (f *TIAddExceptionFunc) AddDomain(domain string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		Domain:      domain,
		Description: description,
	})
	return f
}

func (f *TIAddExceptionFunc) AddIP(ip string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		IP:          ip,
		Description: description,
	})
	return f
}

func (f *TIAddExceptionFunc) AddSenderMailAddress(senderMailAddress string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		SenderMailAddress: senderMailAddress,
		Description:       description,
	})
	return f
}

func (f *TIAddExceptionFunc) AddFileSHA1(fileSHA1 string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		FileSha1:    fileSHA1,
		Description: description,
	})
	return f
}

func (f *TIAddExceptionFunc) AddFileSHA256(fileSHA256 string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		FileSha256:  fileSHA256,
		Description: description,
	})
	return f
}

func (f *TIAddExceptionFunc) Do(ctx context.Context) (*TIAddExceptionResponce, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *TIAddExceptionFunc) method() string {
	return methodPost
}

func (s *TIAddExceptionFunc) url() string {
	return "/v3.0/threatintel/suspiciousObjectExceptions"
}

func (f *TIAddExceptionFunc) requestBody() io.Reader {
	jsonData, err := json.Marshal(f.Request)
	if err != nil {
		return nil
	}
	return bytes.NewReader(jsonData)
}

func (f *TIAddExceptionFunc) responseStruct() any {
	return &f.Response
}

/*
func (f *SandboxSubmitFileToSandboxFunc) responseHeader() any {
	return &f.ResponseHeaders
}
*/
