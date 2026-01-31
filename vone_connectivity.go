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

// connectivityFunc - check Vision One connectivity
type connectivityFunc struct {
	baseFunc
}

var _ vOneFunc = &connectivityFunc{}

func (v *VOne) Connectivity() *connectivityFunc {
	f := &connectivityFunc{}
	f.baseFunc.init(v)
	return f
}

func (f *connectivityFunc) Do(ctx context.Context) error {
	if err := f.vone.call(ctx, f); err != nil {
		return err
	}
	return nil
}

func (f *connectivityFunc) url() string {
	return "/v3.0/healthcheck/connectivity"
}

func (f *connectivityFunc) responseStruct() any {
	return &ErrorData{}
}

func (f *connectivityFunc) responseBody(b io.ReadCloser) {
	io.Copy(os.Stdout, b)
	b.Close()
}
