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

	"github.com/google/uuid"
)

const (
	InternalServerError = "InternalServerError"
	Unsupported         = "Unsupported"
)

type (
	SandboxSubmissionStatusResponse struct {
		ID                 string        `json:"id"`
		Action             Action        `json:"action"`
		Status             Status        `json:"status"`
		Error              Error         `json:"error"`
		CreatedDateTime    VisionOneTime `json:"createdDateTime"`
		LastActionDateTime VisionOneTime `json:"lastActionDateTime"`
		ResourceLocation   string        `json:"resourceLocation"`
		IsCached           bool          `json:"isCached"`
		Digest             Digest        `json:"digest"`
		Arguments          string        `json:"arguments"`
	}
)

var ErrSubmission = errors.New("submission error")

func (s *SandboxSubmissionStatusResponse) GetError() error {
	if s.Error.Code == ErrorCodeNoError {
		return nil
	}
	return fmt.Errorf("%w: %v: %s", ErrSubmission, s.Error.Code, s.Error.Message)
}

type sandboxSubmissionStatusRequest struct {
	baseRequest
	id       string
	response SandboxSubmissionStatusResponse
}

func (v *VOne) SandboxSubmissionStatus(id string) *sandboxSubmissionStatusRequest {
	f := &sandboxSubmissionStatusRequest{
		id: id,
	}
	f.baseRequest.init(v)
	return f
}

func (f *sandboxSubmissionStatusRequest) Do(ctx context.Context) (*SandboxSubmissionStatusResponse, error) {
	if err := uuid.Validate(f.id); err != nil {
		return nil, fmt.Errorf("submission status: %w", err)
	}
	if err := f.checkUsed(); err != nil {
		return nil, fmt.Errorf("submissions status: %w", err)
	}
	if f.vone.mockup != nil {
		return f.vone.mockup.SubmissionStatus(f)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.response, nil
}

func (f *sandboxSubmissionStatusRequest) url() string {
	return fmt.Sprintf("/v3.0/sandbox/tasks/%s", f.id)
}

func (f *sandboxSubmissionStatusRequest) responseStruct() any {
	return &f.response
}
