/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_submit_file.go - upload file for analysis
*/

package vone

import (
	"context"
	"io"
	"os"
)

// connectivityRequest - check Vision One connectivity
type connectivityRequest struct {
	baseRequest
}

var _ vOneRequest = &connectivityRequest{}

func (v *VOne) Connectivity() *connectivityRequest {
	f := &connectivityRequest{}
	f.baseRequest.init(v)
	return f
}

func (f *connectivityRequest) Do(ctx context.Context) error {
	if err := f.vone.call(ctx, f); err != nil {
		return err
	}
	return nil
}

func (f *connectivityRequest) url() string {
	return "/v3.0/healthcheck/connectivity"
}

func (f *connectivityRequest) responseStruct() any {
	return &ErrorData{}
}

func (f *connectivityRequest) responseBody(b io.ReadCloser) {
	io.Copy(os.Stdout, b)
	b.Close()
}
