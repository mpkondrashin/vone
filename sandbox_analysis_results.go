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

// SandboxAnalysisResultsResponse - structure of VisionOne sandbox analysis results JSON
type SandboxAnalysisResultsResponse struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Digest struct {
		MD5    string `json:"md5"`
		SHA1   string `json:"sha1"`
		SHA256 string `json:"sha256"`
	} `json:"digest"`
	Arguments                  string        `json:"arguments"`
	AnalysisCompletionDateTime VisionOneTime `json:"analysisCompletionDateTime"`
	RiskLevel                  RiskLevel     `json:"riskLevel"`
	DetectionNames             []string      `json:"detectionNames"`
	ThreatTypes                []string      `json:"threatTypes"`
	TrueFileType               string        `json:"trueFileType"`
}

type SandboxAnalysisResultsFunc struct {
	baseFunc
	id       string
	Response SandboxAnalysisResultsResponse
}

var _ vOneFunc = &SandboxAnalysisResultsFunc{}

// SandboxAnalysisResults - get function that downloads sanbox analysis results
func (v *VOne) SandboxAnalysisResults(id string) *SandboxAnalysisResultsFunc {
	f := &SandboxAnalysisResultsFunc{id: id}
	f.baseFunc.init(v)
	return f
}

// Do - get sanbox analysis results
func (f *SandboxAnalysisResultsFunc) Do(ctx context.Context) (*SandboxAnalysisResultsResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *SandboxAnalysisResultsFunc) method() string {
	return methodGet
}

func (s *SandboxAnalysisResultsFunc) url() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s", s.id)
}

func (f *SandboxAnalysisResultsFunc) responseStruct() any {
	return &f.Response
}
