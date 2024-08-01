/*
	Trend Micro Vision One API SDK
	(c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	it_add_exception.go - add exceptions for Threat Intelligance
*/

package vone

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

type (
	// TIAddException - Add exceptions payload json
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

	// TIAddExceptionResponce - Add exceptions response json struct
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
	Request  TIAddException
	Response TIAddExceptionResponce
}

var _ vOneFunc = &TIAddExceptionFunc{}

// AddExceptions - return new TIAddExceptionFunc struct
func (v *VOne) AddExceptions() *TIAddExceptionFunc {
	f := &TIAddExceptionFunc{}
	f.baseFunc.init(v)
	return f
}

// Add - add new exception
func (f *TIAddExceptionFunc) Add(exception TIException) *TIAddExceptionFunc {
	f.Request = append(f.Request, exception)
	return f
}

// AddSO - add new exception of certain type
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *TIAddExceptionFunc) AddSO(so SO, value string, description string) *TIAddExceptionFunc {
	e := TIException{
		Description: description,
	}
	switch so {
	case SODomain:
		e.Domain = value
	case SOIP:
		e.IP = value
	case SOSenderMailAddress:
		e.SenderMailAddress = value
	case SOFileSha1:
		e.FileSha1 = value
	case SOFileSha256:
		e.FileSha256 = value
	}
	f.Add(e)
	return f
}

// AddURL - add new URL exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *TIAddExceptionFunc) AddURL(url string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		URL:         url,
		Description: description,
	})
	return f
}

// AddDomain - add new Domain exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *TIAddExceptionFunc) AddDomain(domain string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		Domain:      domain,
		Description: description,
	})
	return f
}

// AddDomain - add new IP exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *TIAddExceptionFunc) AddIP(ip string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		IP:          ip,
		Description: description,
	})
	return f
}

// AddSenderMailAddress - add new SenderMailAddress exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *TIAddExceptionFunc) AddSenderMailAddress(senderMailAddress string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		SenderMailAddress: senderMailAddress,
		Description:       description,
	})
	return f
}

// AddFileSHA1 - add new SHA1 exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *TIAddExceptionFunc) AddFileSHA1(fileSHA1 string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		FileSha1:    fileSHA1,
		Description: description,
	})
	return f
}

// AddFileSHA256 - add new SHA256 exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *TIAddExceptionFunc) AddFileSHA256(fileSHA256 string, description string) *TIAddExceptionFunc {
	f.Add(TIException{
		FileSha256:  fileSHA256,
		Description: description,
	})
	return f
}

// Do - execute the API call
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
