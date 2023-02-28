/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_suspicious_objects.go - get list of IoC for previously submitted object ID
*/

package vone

import (
	"context"
	"fmt"
	"time"
)

type SandboxSuspiciousObjectsResponse struct {
	Items []struct {
		RiskLevel                  RiskLevel `json:"riskLevel"`
		AnalysisCompletionDateTime time.Time `json:"analysisCompletionDateTime"`
		ExpiredDateTime            time.Time `json:"expiredDateTime"`
		RootSHA1                   string    `json:"rootSha1"`
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

func (f *SandboxSuspiciousObjectsFunc) Do(ctx context.Context) (*SandboxSuspiciousObjectsResponse, error) {
	if err := f.vone.Call(ctx, f); err != nil {
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
