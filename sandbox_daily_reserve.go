/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_daily_reserve.go - get dalily usage quota
*/

package vone

import (
	"context"
	"fmt"
)

type SubmissionCountDetail struct {
	FileCount          int `json:"fileCount"`
	FileExemptionCount int `json:"fileExemptionCount"`
	URLCount           int `json:"urlCount"`
	URLExemptionCount  int `json:"urlExemptionCount"`
}

type SandboxDailyReserveResponse struct {
	SubmissionReserveCount   int                   `json:"submissionReserveCount"`
	SubmissionRemainingCount int                   `json:"submissionRemainingCount"`
	SubmissionCount          int                   `json:"submissionCount"`
	SubmissionExemptionCount int                   `json:"submissionExemptionCount"`
	SubmissionCountDetail    SubmissionCountDetail `json:"submissionCountDetail"`
}

type sandboxDailyReserveRequest struct {
	baseRequest
	response SandboxDailyReserveResponse
}

var _ vOneRequest = &sandboxDailyReserveRequest{}

func (f *sandboxDailyReserveRequest) Do(ctx context.Context) (*SandboxDailyReserveResponse, error) {
	if err := f.checkUsed(); err != nil {
		return nil, fmt.Errorf("daily reserve: %w", err)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, fmt.Errorf("daily reserve: %w", err)
	}
	return &f.response, nil
}

func (v *VOne) SandboxDailyReserve() *sandboxDailyReserveRequest {
	f := &sandboxDailyReserveRequest{}
	f.baseRequest.init(v)
	return f
}

func (s *sandboxDailyReserveRequest) url() string {
	return "/v3.0/sandbox/submissionUsage"
}

func (f *sandboxDailyReserveRequest) responseStruct() any {
	return &f.response
}
