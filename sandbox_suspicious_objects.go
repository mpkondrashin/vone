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
)

type SandboxSuspiciousObjectsResponse struct {
	Items []struct {
		RiskLevel                  RiskLevel     `json:"riskLevel"`
		AnalysisCompletionDateTime VisionOneTime `json:"analysisCompletionDateTime"`
		ExpiredDateTime            VisionOneTime `json:"expiredDateTime"`
		RootSHA1                   string        `json:"rootSha1"`
		IP                         string        `json:"ip"`
		URL                        string        `json:"url"`
		FileSHA1                   string        `json:"fileSha1"`
		Domain                     string        `json:"domain"`
	} `json:"items"`
}

type sandboxSuspiciousObjectsFunc struct {
	baseFunc
	id       string
	response SandboxSuspiciousObjectsResponse
}

func (v *VOne) SandboxSuspiciousObjects(id string) *sandboxSuspiciousObjectsFunc {
	f := &sandboxSuspiciousObjectsFunc{id: id}
	f.baseFunc.init(v)
	return f
}

func (f *sandboxSuspiciousObjectsFunc) Do(ctx context.Context) (*SandboxSuspiciousObjectsResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.response, nil
}

func (f *sandboxSuspiciousObjectsFunc) method() string {
	return methodGet
}

func (f *sandboxSuspiciousObjectsFunc) url() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s/suspiciousObjects", f.id)
}

func (f *sandboxSuspiciousObjectsFunc) responseStruct() any {
	return &f.response
}
