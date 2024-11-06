/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	funcs.go - base struct for all Web API functions
*/

package vone

import (
	"io"
	"net/http"
)

type vOneFunc interface {
	method() string             // GET, POST, ...
	url() string                // last part of URI
	uri() string                // full URI (with https://xdr...)
	requestBody() io.Reader     // Body
	populate(*http.Request)     // Pupulate request path and headers
	contentType() string        // application/json by default
	responseStruct() any        // Pointer to struct/slice to parse JSON
	responseHeader() any        // Return struct to populate with response headers
	responseBody(io.ReadCloser) // process body - is called only if responseStruct returns any
}

var _ vOneFunc = &baseFunc{}

type baseFunc struct {
	vone       *VOne
	parameters map[string]string
	headers    map[string]string
}

func (f *baseFunc) method() string {
	return methodGet
}

func (f *baseFunc) url() string {
	return ""
}

func (f *baseFunc) uri() string {
	return ""
}

func (f *baseFunc) requestBody() io.Reader {
	return nil
}

func (f *baseFunc) contentType() string {
	return applicationJSON
}

func (f *baseFunc) responseStruct() any {
	return nil
}

func (f *baseFunc) responseBody(io.ReadCloser) {
}

func (f *baseFunc) init(vone *VOne) {
	f.vone = vone
	f.parameters = make(map[string]string)
	f.headers = make(map[string]string)
}

func (f *baseFunc) setHeader(name, value string) {
	f.headers[name] = value
}

func (f *baseFunc) setParameter(name, value string) {
	f.parameters[name] = value
}

func (f *baseFunc) populate(req *http.Request) {
	q := req.URL.Query()
	for key, value := range f.parameters {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	for key, value := range f.headers {
		req.Header.Add(key, value)
	}
}

func (f *baseFunc) responseHeader() any {
	return nil
}
