/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_download_results.go - get file check results
*/

package vone

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
)

type sandboxDownloadResultsRequest struct {
	baseRequest
	id       string
	response io.ReadCloser
}

var _ vOneRequest = &sandboxDownloadResultsRequest{}

func (v *VOne) SandboxDownloadResults(id string) *sandboxDownloadResultsRequest {
	f := &sandboxDownloadResultsRequest{id: id}
	f.baseRequest.init(v)
	return f
}

// Do performs the request and returns a ReadCloser.
// Only call Do() or Store() once; the returned stream must be consumed immediately.
func (f *sandboxDownloadResultsRequest) Do(ctx context.Context) (io.ReadCloser, error) {
	if err := uuid.Validate(f.id); err != nil {
		return nil, fmt.Errorf("download result/investigation package: %w", err)
	}
	if err := f.checkUsed(); err != nil {
		return nil, fmt.Errorf("download result/investigation package: %w", err)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return f.response, nil
}

// Store writes result to given filename
// Only call Do() or Store() once; the returned stream must be consumed immediately.
func (f *sandboxDownloadResultsRequest) Store(ctx context.Context, filePath string) error {
	if _, err := f.Do(ctx); err != nil {
		return err
	}
	if f.response != nil {
		defer f.response.Close()
	}
	output, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer output.Close()
	_, err = io.Copy(output, f.response)
	return err
}

func (f *sandboxDownloadResultsRequest) url() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s/report", f.id)
}

func (f *sandboxDownloadResultsRequest) responseBody(body io.ReadCloser) {
	f.response = body
}
