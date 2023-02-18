/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	VOne - Web API SDK

	vone.go - main SDK declarations
*/

package vone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	GET = "GET"

	application_json = "application/json"

	timeFormat = "2006-1-02T15:04:05Z"
)

type Error struct {
	ErrorData struct {
		Message    string `json:"message"`
		Code       string `json:"code"`
		Innererror struct {
			Service string `json:"service"`
			Code    string `json:"code"`
		} `json:"innererror"`
	} `json:"error"`
}

func (e *Error) Error() string {
	return e.ErrorData.Code + ". " + e.ErrorData.Message
}

type VOne struct {
	urlBase string
	bearer  string
}

func NewVOne(url string, bearer string) *VOne {
	return &VOne{
		urlBase: url,
		bearer:  bearer,
	}
}

func (v *VOne) RequestJSON(method, url string, bodyData any) (io.ReadCloser, error) {
	var body io.Reader
	if bodyData != nil {
		buffer, err := json.Marshal(bodyData)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(buffer)
	}
	return v.Request(method, url, body, application_json)
}

func (v *VOne) Request(method, url string, body io.Reader, contentType string) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, v.urlBase+url, body)
	if err != nil {
		return nil, fmt.Errorf("VOne: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+v.bearer)
	if body != nil {
		req.Header.Set("Content-Type", contentType)
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
		vOneErr := new(Error)
		if err := json.Unmarshal(data.Bytes(), vOneErr); err != nil {
			return nil, err
		}
		return nil, vOneErr
	}
	return resp.Body, nil
}

func (v *VOne) Call(f Func) error {
	uri := f.URI()
	if uri == "" {
		uri = v.urlBase + f.URL()
	}
	return v.CallURL(f, uri)
}

func (v *VOne) CallURL(f Func, uri string) error {
	req, err := http.NewRequest(f.Method(), uri, f.RequestBody())
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+v.bearer)
	if f.RequestBody() != nil {
		req.Header.Set("Content-Type", f.ContentType())
	}
	f.Populate(req)
	//log.Println("EEE", req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http.Client.Do: %v", err)
	}
	//log.Println("EEE RESP", resp)
	if resp.StatusCode > 299 {
		var data bytes.Buffer
		if _, err := io.Copy(&data, resp.Body); err != nil {
			return fmt.Errorf("io.Copy: %v", err)
		}
		//log.Printf("respond: %v\n", data.String())
		vOneErr := new(Error)
		if err := json.Unmarshal(data.Bytes(), vOneErr); err != nil {
			return fmt.Errorf("json.Unmarshal: %v", err)
		}
		return fmt.Errorf("vOneErr: %w", vOneErr)
	}
	if f.ResponseStruct() == nil {
		f.ResponseBody(resp.Body)
		return nil
	}
	return v.DecodeBody(f, resp.Body)
}

func (v *VOne) DecodeBody(f Func, body io.ReadCloser) error {
	defer body.Close()
	err := json.NewDecoder(body).Decode(f.ResponseStruct())
	if err != nil && err != io.EOF {
		return fmt.Errorf("response parse error: %w", err)
	}
	return nil
}

/*
func (v *VOne) Get(url string, body any, contentType string) (io.ReadCloser, error) {
	return v.Request("GET", url, body, contentType)
}

func (v *VOne) Post(url string, body any, contentType string) (io.ReadCloser, error) {
	return v.Request("POST", url, body, contentType)
}
*/
