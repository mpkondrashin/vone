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

type SandboxSuspiciousObject struct {
	RiskLevel                  RiskLevel     `json:"riskLevel"`
	AnalysisCompletionDateTime VisionOneTime `json:"analysisCompletionDateTime"`
	ExpiredDateTime            VisionOneTime `json:"expiredDateTime"`
	RootSHA1                   string        `json:"rootSha1"`
	IP                         string        `json:"ip"`
	URL                        string        `json:"url"`
	FileSHA1                   string        `json:"fileSha1"`
	Domain                     string        `json:"domain"`
}

type SandboxSuspiciousObjectsResponse struct {
	Items []SandboxSuspiciousObject `json:"items"`
}

type sandboxSuspiciousObjectsRequest struct {
	baseRequest
	id       string
	response SandboxSuspiciousObjectsResponse
}

func (v *VOne) SandboxSuspiciousObjects(id string) *sandboxSuspiciousObjectsRequest {
	f := &sandboxSuspiciousObjectsRequest{id: id}
	f.baseRequest.init(v)
	return f
}

func (f *sandboxSuspiciousObjectsRequest) Do(ctx context.Context) (*SandboxSuspiciousObjectsResponse, error) {
	if err := f.checkUsed(); err != nil {
		return nil, nil, fmt.Errorf("syspicious objects: %w", err)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.response, nil
}

func (f *sandboxSuspiciousObjectsRequest) url() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s/suspiciousObjects", f.id)
}

func (f *sandboxSuspiciousObjectsRequest) responseStruct() any {
	return &f.response
}
