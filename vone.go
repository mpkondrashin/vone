package vone

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	GET = "GET"
)

type VOneError struct {
	HttpError int
}

func VOneErr(httpError int) VOneError {
	return VOneError{
		HttpError: httpError,
	}
}

func (e VOneError) Error() string {
	return strconv.Itoa(e.HttpError)
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

func (v *VOne) Request(method, url string, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, v.urlBase+url, body)
	if err != nil {
		return nil, fmt.Errorf("VOne: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+v.bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return nil, VOneErr(resp.StatusCode)
	}
	return resp.Body, nil
}

func (v *VOne) Get(url string, body io.Reader) (io.ReadCloser, error) {
	return v.Request("GET", url, body)
}
