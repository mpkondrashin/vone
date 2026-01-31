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

	// TIException - single expection data. Only one other field beside Description can be non empty string
	TIException struct {
		URL               string `json:"url,omitempty"`
		Description       string `json:"description,omitempty"`
		Domain            string `json:"domain,omitempty"`
		IP                string `json:"ip,omitempty"`
		SenderMailAddress string `json:"senderMailAddress,omitempty"`
		FileSha1          string `json:"fileSha1,omitempty"`
		FileSha256        string `json:"fileSha256,omitempty"`
	}

	// TIAddExceptionResponse - Add exceptions response json struct
	TIAddExceptionResponse []struct {
		Status int `json:"status"`
		Body   struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		} `json:"body,omitempty"`
	}
)

// tiAddExceptionRequest - function to add exceptions
type tiAddExceptionRequest struct {
	baseRequest
	request  TIAddException
	response TIAddExceptionResponse
}

var _ vOneRequest = &tiAddExceptionRequest{}

// AddExceptions - return new TIAddExceptionFunc struct
func (v *VOne) AddExceptions() *tiAddExceptionRequest {
	f := &tiAddExceptionRequest{}
	f.baseRequest.init(v)
	return f
}

// Add - add new exception
func (f *tiAddExceptionRequest) Add(exception TIException) *tiAddExceptionRequest {
	f.request = append(f.request, exception)
	return f
}

// AddSO - add new exception of certain type
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *tiAddExceptionRequest) AddSO(so SO, value string, description string) *tiAddExceptionRequest {
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

// URL - add new URL exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *tiAddExceptionRequest) URL(url string, description string) *tiAddExceptionRequest {
	f.Add(TIException{
		URL:         url,
		Description: description,
	})
	return f
}

// Domain - add new Domain exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *tiAddExceptionRequest) Domain(domain string, description string) *tiAddExceptionRequest {
	f.Add(TIException{
		Domain:      domain,
		Description: description,
	})
	return f
}

// IP - add new IP exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *tiAddExceptionRequest) IP(ip string, description string) *tiAddExceptionRequest {
	f.Add(TIException{
		IP:          ip,
		Description: description,
	})
	return f
}

// SenderMailAddress - add new SenderMailAddress exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *tiAddExceptionRequest) SenderMailAddress(senderMailAddress string, description string) *tiAddExceptionRequest {
	f.Add(TIException{
		SenderMailAddress: senderMailAddress,
		Description:       description,
	})
	return f
}

// FileSHA1 - add new SHA1 exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *tiAddExceptionRequest) FileSHA1(fileSHA1 string, description string) *tiAddExceptionRequest {
	f.Add(TIException{
		FileSha1:    fileSHA1,
		Description: description,
	})
	return f
}

// FileSHA256 - add new SHA256 exception
// The method returns a pointer to the TIAddExceptionFunc, allowing for method chaining.
func (f *tiAddExceptionRequest) FileSHA256(fileSHA256 string, description string) *tiAddExceptionRequest {
	f.Add(TIException{
		FileSha256:  fileSHA256,
		Description: description,
	})
	return f
}

// Do - execute the API call. Example:
//
// err, response := vone.NewVOne(domain, token).AddExceptions().AddIP("8.8.8.8", "Google DNS").Do(context.TODO())
func (f *tiAddExceptionRequest) Do(ctx context.Context) (*TIAddExceptionResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.response, nil
}

func (f *tiAddExceptionRequest) method() string {
	return methodPost
}

func (s *tiAddExceptionRequest) url() string {
	return "/v3.0/threatintel/suspiciousObjectExceptions"
}

func (f *tiAddExceptionRequest) requestBody() io.Reader {
	jsonData, err := json.Marshal(f.request)
	if err != nil {
		return nil
	}
	return bytes.NewReader(jsonData)
}

func (f *tiAddExceptionRequest) responseStruct() any {
	return &f.response
}
