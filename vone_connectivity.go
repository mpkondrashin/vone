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

// ConnectivityFunc - check Vision One connectivity
type ConnectivityFunc struct {
	baseFunc
}

var _ vOneFunc = &ConnectivityFunc{}

func (v *VOne) Connectivity() *ConnectivityFunc {
	f := &ConnectivityFunc{}
	f.baseFunc.init(v)
	return f
}

func (f *ConnectivityFunc) Do(ctx context.Context) error {
	if err := f.vone.call(ctx, f); err != nil {
		return err
	}
	return nil
}

func (f *ConnectivityFunc) url() string {
	return "/v3.0/healthcheck/connectivity"
}

func (f *ConnectivityFunc) responseStruct() any {
	return &ErrorData{}
}

func (f *ConnectivityFunc) responseBody(b io.ReadCloser) {
	io.Copy(os.Stdout, b)
	b.Close()
}
