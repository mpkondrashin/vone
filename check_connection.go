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

type sanboxCheckConnectionRequest struct {
	baseFunc
	response CheckConnectionResponse
}

var _ vOneFunc = &sanboxCheckConnectionRequest{}

func (f *sanboxCheckConnectionRequest) Do(ctx context.Context) (*CheckConnectionResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.response, nil
}

func (v *VOne) CheckConnection() *sanboxCheckConnectionRequest {
	f := &sanboxCheckConnectionRequest{}
	f.baseFunc.init(v)
	return f
}

func (s *sanboxCheckConnectionRequest) url() string {
	return "/v3.0/healthcheck/connectivity"
}

func (f *sanboxCheckConnectionRequest) responseStruct() any {
	return &f.response
}
