package vone

import (
	"encoding/json"
	"fmt"
)

type SandboxGetDailyReserveData struct {
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

func (v *VOne) SandboxGetDailyReserve() (*SandboxGetDailyReserveData, error) {
	url := "/v3.0/sandbox/submissionUsage"
	body, err := v.RequestJSON("GET", url, nil)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(body)
	data := new(SandboxGetDailyReserveData)
	if err := decoder.Decode(data); err != nil {
		return nil, fmt.Errorf("%s response error: %w", url, err)
	}
	return data, nil
}

type SanboxDailyReserveFunc struct {
	BaseFunc
	Result *SandboxGetDailyReserveData
}

var _ Func = &SanboxDailyReserveFunc{}

func (v *VOne) SanboxDailyReserve() (*SandboxGetDailyReserveData, error) {
	f, err := NewSandboxDailyReserve()
	if err != nil {
		return nil, fmt.Errorf("SanboxDailyReserve: %w", err)
	}
	if err := v.Call(f); err != nil {
		return nil, fmt.Errorf("Call: %w", err)
	}
	return f.Result, nil
}

func NewSandboxDailyReserve() (*SanboxDailyReserveFunc, error) {
	f := &SanboxDailyReserveFunc{}
	f.BaseFunc.Init()
	return f, nil
}

func (f *SanboxDailyReserveFunc) Method() string {
	return "POST"
}

func (s *SanboxDailyReserveFunc) URL() string {
	return "/v3.0/sandbox/submissionUsage"
}
