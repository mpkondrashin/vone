/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_analysis_results.go - get analysis result
*/

package vone

import (
	"context"
	"fmt"
)

type Digest struct {
	MD5    string `json:"md5"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
}

// SandboxAnalysisResultsResponse - structure of VisionOne sandbox analysis results JSON
type SandboxAnalysisResultsResponse struct {
	ID                         string        `json:"id"`
	Type                       string        `json:"type"`
	Digest                     Digest        `json:"digest"`
	Arguments                  string        `json:"arguments"`
	AnalysisCompletionDateTime VisionOneTime `json:"analysisCompletionDateTime"`
	RiskLevel                  RiskLevel     `json:"riskLevel"`
	DetectionNames             []string      `json:"detectionNames"`
	ThreatTypes                []string      `json:"threatTypes"`
	TrueFileType               string        `json:"trueFileType"`
}

type sandboxAnalysisResultsFunc struct {
	baseFunc
	id       string
	response SandboxAnalysisResultsResponse
}

var _ vOneFunc = &sandboxAnalysisResultsFunc{}

// SandboxAnalysisResults - get function that downloads sanbox analysis results
func (v *VOne) SandboxAnalysisResults(id string) *sandboxAnalysisResultsFunc {
	f := &sandboxAnalysisResultsFunc{id: id}
	f.baseFunc.init(v)
	return f
}

// Do - get sanbox analysis results
func (f *sandboxAnalysisResultsFunc) Do(ctx context.Context) (*SandboxAnalysisResultsResponse, error) {
	if f.vone.mockup != nil {
		return f.vone.mockup.AnalysisResults(f)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.response, nil
}

func (f *sandboxAnalysisResultsFunc) method() string {
	return methodGet
}

func (s *sandboxAnalysisResultsFunc) url() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s", s.id)
}

func (f *sandboxAnalysisResultsFunc) responseStruct() any {
	return &f.response
}
