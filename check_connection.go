/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	check_connection.go - check connectivity and API token
*/

package vone

import "context"

type CheckConnectionResponse struct {
	Status string `json:"status"`
}

type SanboxCheckConnectionFunc struct {
	baseFunc
	Response CheckConnectionResponse
}

var _ vOneFunc = &SanboxCheckConnectionFunc{}

func (f *SanboxCheckConnectionFunc) Do(ctx context.Context) (*CheckConnectionResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (v *VOne) CheckConnection() *SanboxCheckConnectionFunc {
	f := &SanboxCheckConnectionFunc{}
	f.baseFunc.init(v)
	return f
}

func (s *SanboxCheckConnectionFunc) url() string {
	return "/v3.0/healthcheck/connectivity"
}

func (f *SanboxCheckConnectionFunc) responseStruct() any {
	return &f.Response
}
