/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_investigation_package.go - download investigation package
*/

package vone

import (
	"context"
	"fmt"
	"io"
	"os"
)

type sandboxInvestigationPackageRequest struct {
	sandboxDownloadResultsRequest
}

func (v *VOne) SandboxInvestigationPackage(id string) *sandboxInvestigationPackageRequest {
	return &sandboxInvestigationPackageRequest{*v.SandboxDownloadResults(id)}
}

func (f *sandboxInvestigationPackageRequest) Do(ctx context.Context) (io.ReadCloser, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return f.response, nil
}

func (f *sandboxInvestigationPackageRequest) Store(ctx context.Context, filePath string) error {
	if _, err := f.Do(ctx); err != nil {
		return nil
	}
	defer f.response.Close()
	output, err := os.Create(filePath)
	if err != nil {
		return nil
	}
	defer output.Close()
	_, err = io.Copy(output, f.response)
	return err
}

func (s *sandboxInvestigationPackageRequest) url() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s/investigationPackage", s.id)
}
