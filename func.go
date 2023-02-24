package vone

import (
	"io"
	"net/http"
)

type Func interface {
	Method() string             // GET, POST, ...
	URL() string                // last part of URI
	URI() string                // full URI (with https://xdr...)
	RequestBody() io.Reader     // Body
	Populate(*http.Request)     // Pupulate request path and headers
	ContentType() string        // application/json by default
	ResponseStruct() any        // Pointer to struct/slice to parse JSON
	ResponseHeader() any        // Return struct to populate with response headers
	ResponseBody(io.ReadCloser) // process body - is called only if ResponseStruct returns any
}

var _ Func = &BaseFunc{}

type BaseFunc struct {
	vone       *VOne
	parameters map[string]string
	headers    map[string]string
}

func (f *BaseFunc) Method() string {
	return "GET"
}

func (f *BaseFunc) URL() string {
	return ""
}

func (f *BaseFunc) URI() string {
	return ""
}

func (f *BaseFunc) RequestBody() io.Reader {
	return nil
}

func (f *BaseFunc) ContentType() string {
	return "application/json"
}

func (f *BaseFunc) ResponseStruct() any {
	return &Error{}
}

func (f *BaseFunc) ResponseBody(io.ReadCloser) {
}

func (f *BaseFunc) Init(vone *VOne) {
	f.vone = vone
	f.parameters = make(map[string]string)
	f.headers = make(map[string]string)
}

func (f *BaseFunc) SetHeader(name, value string) {
	f.headers[name] = value
}

func (f *BaseFunc) SetParameter(name, value string) {
	f.parameters[name] = value
}

func (f *BaseFunc) Populate(req *http.Request) {
	q := req.URL.Query()
	for key, value := range f.parameters {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	for key, value := range f.parameters {
		req.Header.Add(key, value)
	}
}

func (f *BaseFunc) ResponseHeader() any {
	return nil
}
