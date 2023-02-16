package vone

import (
	"fmt"
	"time"
)

type SandboxSuspiciousObjectsResponse struct {
	Items []struct {
		RiskLevel                  string    `json:"riskLevel"`
		AnalysisCompletionDateTime time.Time `json:"analysisCompletionDateTime"`
		ExpiredDateTime            time.Time `json:"expiredDateTime"`
		RootSha1                   string    `json:"rootSha1"`
		IP                         string    `json:"ip"`
	} `json:"items"`
}

type SandboxSuspiciousObjectsFunc struct {
	BaseFunc
	id       string
	Response SandboxSuspiciousObjectsResponse
}

func (v *VOne) SandboxSuspiciousObjects(id string) *SandboxSuspiciousObjectsFunc {
	f := &SandboxSuspiciousObjectsFunc{id: id}
	f.BaseFunc.Init(v)
	return f
}

func (f *SandboxSuspiciousObjectsFunc) Do() (*SandboxSuspiciousObjectsResponse, error) {
	if err := f.vone.Call(f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *SandboxSuspiciousObjectsFunc) Method() string {
	return "GET"
}

func (f *SandboxSuspiciousObjectsFunc) URL() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s/suspiciousObjects", f.id)
}

func (f *SandboxSuspiciousObjectsFunc) ResponseStruct() any {
	return &f.Response
}
