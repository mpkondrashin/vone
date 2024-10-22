/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_submission_status.go - status of submission
*/

package vone

import (
	"context"
	"errors"
	"fmt"
)

const (
	InternalServerError = "InternalServerError"
	Unsupported         = "Unsupported"
)

type (
	SandboxSubmissionStatusResponse struct {
		ID     string `json:"id"`
		Action Action `json:"action"`
		Status Status `json:"status"`
		Error  struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		CreatedDateTime    VisionOneTime `json:"createdDateTime"`
		LastActionDateTime VisionOneTime `json:"lastActionDateTime"`
		ResourceLocation   string        `json:"resourceLocation"`
		IsCached           bool          `json:"isCached"`
		Digest             struct {
			MD5    string `json:"md5"`
			SHA1   string `json:"sha1"`
			SHA256 string `json:"sha256"`
		} `json:"digest"`
		Arguments string `json:"arguments"`
	}
)

var ErrSubmission = errors.New("submission error")

func (s *SandboxSubmissionStatusResponse) GetError() error {
	if s.Error.Code == "" {
		return nil
	}
	return fmt.Errorf("%w: %v: %s", ErrSubmission, s.Error.Code, s.Error.Message)
}

type SandboxSubmissionStatusFunc struct {
	baseFunc
	id       string
	Response SandboxSubmissionStatusResponse
}

func (v *VOne) SandboxSubmissionStatus(id string) *SandboxSubmissionStatusFunc {
	f := &SandboxSubmissionStatusFunc{
		id: id,
	}
	f.baseFunc.init(v)
	return f
}

func (f *SandboxSubmissionStatusFunc) Do(ctx context.Context) (*SandboxSubmissionStatusResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *SandboxSubmissionStatusFunc) url() string {
	return fmt.Sprintf("/v3.0/sandbox/tasks/%s", f.id)
}

func (f *SandboxSubmissionStatusFunc) responseStruct() any {
	return &f.Response
}
