/*
	Trend Micro Vision One API SDK
	(c) 2026 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Workbench API capabilities

	workbench_modify_status.go - modify workbench alert status
*/

package vone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type (
	// WorkbenchModifyStatusRequest - request body for modifying alert status
	WorkbenchModifyStatusRequest struct {
		Status              AlertStatus `json:"status,omitempty"`
		InvestigationResult InvestigationResult      `json:"investigationResult,omitempty"`
	}
)

type workbenchModifyStatusRequest struct {
	baseRequest
	id       string
	request  WorkbenchModifyStatusRequest
}

var _ vOneRequest = &workbenchModifyStatusRequest{}

// WorkbenchModifyStatus - create a request to modify workbench alert status
func (v *VOne) WorkbenchModifyStatus(id string) *workbenchModifyStatusRequest {
	f := &workbenchModifyStatusRequest{id: id}
	f.baseRequest.init(v)
	return f
}

// Status - set the alert status (e.g., "Open", "Closed", "In Progress")
func (f *workbenchModifyStatusRequest) Status(status AlertStatus) *workbenchModifyStatusRequest {
	f.request.Status = status
	return f
}

// InvestigationResult - set the investigation result (e.g., "True Positive", "False Positive", "No Findings")
func (f *workbenchModifyStatusRequest) InvestigationResult(result InvestigationResult) *workbenchModifyStatusRequest {
	f.request.InvestigationResult = result
	return f
}

// IfMatch - set the If-Match header for optimistic concurrency control
func (f *workbenchModifyStatusRequest) IfMatch(etag string) *workbenchModifyStatusRequest {
	f.setHeader("If-Match", etag)
	return f
}

// Do - execute the modify status request
func (f *workbenchModifyStatusRequest) Do(ctx context.Context) error {
	if f.id == "" {
		return fmt.Errorf("WorkbenchModifyStatus: alert ID is required")
	}
	if f.vone.mockup != nil {
		// Add mockup support if needed
		return fmt.Errorf("WorkbenchModifyStatus: %w", ErrNotImplemented)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return fmt.Errorf("WorkbenchModifyStatus: %w", err)
	}
	return nil
}

func (f *workbenchModifyStatusRequest) method() string {
	return methodPatch
}

func (f *workbenchModifyStatusRequest) url() string {
	return fmt.Sprintf("/v3.0/workbench/alerts/%s", f.id)
}

func (f *workbenchModifyStatusRequest) requestBody() io.Reader {
	data, err := json.Marshal(f.request)
	if err != nil {
		return nil
	}
	return bytes.NewBuffer(data)
}