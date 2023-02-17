/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_submission_status.go - status of submission
*/

package vone

import (
	"fmt"
	"time"
)

type (
	SandboxSubmissionStatusResponse struct {
		ID     string `json:"id"`
		Action string `json:"action"`
		Status string `json:"status"`
		Error  struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		CreatedDateTime    time.Time `json:"createdDateTime"`
		LastActionDateTime time.Time `json:"lastActionDateTime"`
		ResourceLocation   string    `json:"resourceLocation"`
		IsCached           bool      `json:"isCached"`
		Digest             struct {
			MD5    string `json:"md5"`
			SHA1   string `json:"sha1"`
			SHA256 string `json:"sha256"`
		} `json:"digest"`
		Arguments string `json:"arguments"`
	}
)

type SandboxSubmissionStatusFunc struct {
	BaseFunc
	id       string
	Response SandboxSubmissionStatusResponse
}

func (v *VOne) SandboxSubmissionStatus(id string) *SandboxSubmissionStatusFunc {
	f := &SandboxSubmissionStatusFunc{
		id: id,
	}
	f.BaseFunc.Init(v)
	return f
}

func (f *SandboxSubmissionStatusFunc) Do() (*SandboxSubmissionStatusResponse, error) {
	if err := f.vone.Call(f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *SandboxSubmissionStatusFunc) URL() string {
	return fmt.Sprintf("/v3.0/sandbox/tasks/%s", f.id)
}

func (f *SandboxSubmissionStatusFunc) ResponseStruct() any {
	return &f.Response
}
