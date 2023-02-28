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

type SandboxInvestigationPackageFunc struct {
	SandboxDownloadResultsFunc
}

func (v *VOne) SandboxInvestigationPackage(id string) *SandboxInvestigationPackageFunc {
	return &SandboxInvestigationPackageFunc{*v.SandboxDownloadResults(id)}
}

func (f *SandboxInvestigationPackageFunc) Do(ctx context.Context) (io.ReadCloser, error) {
	if err := f.vone.Call(ctx, f); err != nil {
		return nil, err
	}
	return f.Response, nil
}

func (f *SandboxInvestigationPackageFunc) Store(ctx context.Context, filePath string) error {
	if _, err := f.Do(ctx); err != nil {
		return nil
	}
	defer f.Response.Close()
	output, err := os.Create(filePath)
	if err != nil {
		return nil
	}
	defer output.Close()
	_, err = io.Copy(output, f.Response)
	return err
}

func (s *SandboxInvestigationPackageFunc) URL() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s/investigationPackage", s.id)
}
