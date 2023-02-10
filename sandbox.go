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
	body, err := v.Get(url, nil)
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
