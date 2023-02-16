package vone

import (
	"io"
	"net/http"
)

type Func interface {
	Method() string         // GET, POST, ...
	URL() string            // last part of URI
	URI() string            // full URI (with https://xdr...)
	RequestBody() io.Reader // Body
	ContentType() string    // application/json by default
	ResponseStruct() any    // Pointer to struct/slice to parse JSON
	Populate(*http.Request) // Pupulate request path and headers
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
	return &VOneError{}
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

/*
func (f *BaseFunc) Perform(method, url string, body io.Reader, contentType string) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, f.vone.urlBase+url, body)
	if err != nil {
		return nil, fmt.Errorf("VOne: %w", err)
	}
	if body != nil {
		f.SetHeader("Content-Type", contentType)
	}
	f.Populate(req)
	for key, h := range req.Header {
		log.Println("H", key, h)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("VOne: %v", err)
	}
	if resp.StatusCode > 299 {
		var data bytes.Buffer
		if _, err := io.Copy(&data, resp.Body); err != nil {
			return nil, err
		}
		vOneErr := new(VOneError)
		if err := json.Unmarshal(data.Bytes(), vOneErr); err != nil {
			return nil, err
		}
		return nil, vOneErr
	}
	return resp.Body, nil
}
*/
