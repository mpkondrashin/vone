/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_download_results.go - get file check results
*/
package vone

import (
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

func (f *SandboxInvestigationPackageFunc) Do() (io.ReadCloser, error) {
	if err := f.vone.Call(f); err != nil {
		return nil, err
	}
	return f.Response, nil
}

func (f *SandboxInvestigationPackageFunc) Store(filePath string) error {
	if _, err := f.Do(); err != nil {
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
