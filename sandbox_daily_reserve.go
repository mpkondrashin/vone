/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_daily_reserve.go - get dalily usage quota
*/

package vone

import "context"

type SandboxDailyReserveResponse struct {
	SubmissionReserveCount   int `json:"submissionReserveCount"`
	SubmissionRemainingCount int `json:"submissionRemainingCount"`
	SubmissionCount          int `json:"submissionCount"`
	SubmissionExemptionCount int `json:"submissionExemptionCount"`
	SubmissionCountDetail    struct {
		FileCount          int `json:"fileCount"`
		FileExemptionCount int `json:"fileExemptionCount"`
		URLCount           int `json:"urlCount"`
		URLExemptionCount  int `json:"urlExemptionCount"`
	} `json:"submissionCountDetail"`
}

type SanboxDailyReserveFunc struct {
	BaseFunc
	Response SandboxDailyReserveResponse
}

func (f *SanboxDailyReserveFunc) Do(ctx context.Context) (*SandboxDailyReserveResponse, error) {
	if err := f.vone.Call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (v *VOne) SandboxDailyReserve() *SanboxDailyReserveFunc {
	f := &SanboxDailyReserveFunc{}
	f.BaseFunc.Init(v)
	return f
}

func (s *SanboxDailyReserveFunc) URL() string {
	return "/v3.0/sandbox/submissionUsage"
}

func (f *SanboxDailyReserveFunc) ResponseStruct() any {
	return &f.Response
}
