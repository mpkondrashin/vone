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

type sanboxDailyReserveFunc struct {
	baseFunc
	response SandboxDailyReserveResponse
}

var _ vOneFunc = &sanboxDailyReserveFunc{}

func (f *sanboxDailyReserveFunc) Do(ctx context.Context) (*SandboxDailyReserveResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.response, nil
}

func (v *VOne) SandboxDailyReserve() *sanboxDailyReserveFunc {
	f := &sanboxDailyReserveFunc{}
	f.baseFunc.init(v)
	return f
}

func (s *sanboxDailyReserveFunc) url() string {
	return "/v3.0/sandbox/submissionUsage"
}

func (f *sanboxDailyReserveFunc) responseStruct() any {
	return &f.response
}
