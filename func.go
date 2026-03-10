/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	funcs.go - base struct for all Web API functions
*/

package vone

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var ErrAlreadyCalled = errors.New("Do() or Store() already called")

type vOneRequest interface {
	method() string                // GET, POST, ...
	url() string                   // last part of URI
	uri() string                   // full URI (with https://xdr...)
	requestBody() io.Reader        // Body
	populateHeaders(*http.Request) // Populate request headers
	populateParameters(*url.URL)   // Populate request parameters
	contentType() string           // application/json by default
	responseStruct() any           // Pointer to struct/slice to parse JSON
	responseHeader() any           // Return struct to populate with response headers
	responseBody(io.ReadCloser)    // process body - is called only if responseStruct returns any
}

var _ vOneRequest = &baseRequest{}

type baseRequest struct {
	vone       *VOne
	parameters map[string]string
	headers    map[string]string
	used       bool
}

func (f *baseRequest) method() string {
	return methodGet
}

func (f *baseRequest) url() string {
	return ""
}

func (f *baseRequest) uri() string {
	return ""
}

func (f *baseRequest) requestBody() io.Reader {
	return nil
}

func (f *baseRequest) contentType() string {
	return applicationJSON
}

func (f *baseRequest) responseStruct() any {
	return nil
}

func (f *baseRequest) responseBody(io.ReadCloser) {
}

func (f *baseRequest) init(vone *VOne) {
	f.vone = vone
	f.parameters = make(map[string]string)
	f.headers = make(map[string]string)
}

func (f *baseRequest) setHeader(name, value string) {
	f.headers[name] = value
}

func (f *baseRequest) setParameter(name, value string) {
	f.parameters[name] = value
}

func (f *baseRequest) populateHeaders(req *http.Request) {
	for key, value := range f.headers {
		req.Header.Add(key, value)
	}
}

func (f *baseRequest) populateParameters(u *url.URL) {
	q := u.Query()
	for key, value := range f.parameters {
		q.Add(key, value)
	}
	u.RawQuery = strings.ReplaceAll(q.Encode(), "%3A", ":")
}

func (f *baseRequest) responseHeader() any {
	return nil
}

func (f *baseRequest) checkUsed() error {
	if f.used {
		return ErrAlreadyCalled
	}
	f.used = true
	return nil
}
