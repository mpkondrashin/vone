package vone

import (
	"fmt"
	"time"
)

type SandboxAnalysisResultsResponse struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Digest struct {
		MD5    string `json:"md5"`
		SHA1   string `json:"sha1"`
		SHA256 string `json:"sha256"`
	} `json:"digest"`
	Arguments                  string    `json:"arguments"`
	AnalysisCompletionDateTime time.Time `json:"analysisCompletionDateTime"`
	RiskLevel                  string    `json:"riskLevel"`
	DetectionNames             []string  `json:"detectionNames"`
	ThreatTypes                []string  `json:"threatTypes"`
	TrueFileType               string    `json:"trueFileType"`
}

type SandboxAnalysisResultsFunc struct {
	BaseFunc
	id       string
	Response SandboxAnalysisResultsResponse
}

func (v *VOne) SandboxAnalysisResults(id string) *SandboxAnalysisResultsFunc {
	f := &SandboxAnalysisResultsFunc{id: id}
	f.BaseFunc.Init(v)
	return f
}

func (f *SandboxAnalysisResultsFunc) Do() (*SandboxAnalysisResultsResponse, error) {
	if err := f.vone.Call(f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *SandboxAnalysisResultsFunc) Method() string {
	return "GET"
}

func (s *SandboxAnalysisResultsFunc) URL() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s", s.id)
}

func (f *SandboxAnalysisResultsFunc) ResponseStruct() any {
	return &f.Response
}
