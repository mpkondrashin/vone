/*
	Trend Micro Vision One API SDK
	(c) 2026 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Workbench API capabilities

	workbench_details.go - get workbench alert details
*/

package vone

import (
	"context"
	"fmt"
)

type workbenchDetailsRequest struct {
	baseRequest
	id       string
	response WorkbenchAlert
}

var _ vOneRequest = &workbenchDetailsRequest{}

// WorkbenchAlertDetails - get function that retrieves workbench alert details by ID
func (v *VOne) WorkbenchAlertDetails(id string) *workbenchDetailsRequest {
	f := &workbenchDetailsRequest{id: id}
	f.baseRequest.init(v)
	return f
}

// Do - get workbench alert details
func (f *workbenchDetailsRequest) Do(ctx context.Context) (*WorkbenchAlert, error) {
	if f.id == "" {
		return nil, fmt.Errorf("WorkbenchAlertDetails: alert ID is required")
	}
	if f.vone.mockup != nil {
		// Add mockup support if needed
		return nil, fmt.Errorf("WorkbenchDetailsRequest: %w", ErrNotImplemented)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, fmt.Errorf("WorkbenchAlertDetails: %w", err)
	}
	return &f.response, nil
}

func (f *workbenchDetailsRequest) method() string {
	return methodGet
}

func (f *workbenchDetailsRequest) url() string {
	return fmt.Sprintf("/v3.0/workbench/alerts/%s", f.id)
}

func (f *workbenchDetailsRequest) responseStruct() any {
	return &f.response
}
