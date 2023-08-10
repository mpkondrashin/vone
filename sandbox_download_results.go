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
)

type SandboxDownloadResultsFunc struct {
	baseFunc
	id       string
	Response io.ReadCloser
}

var _ vOneFunc = &SandboxDownloadResultsFunc{}

func (v *VOne) SandboxDownloadResults(id string) *SandboxDownloadResultsFunc {
	f := &SandboxDownloadResultsFunc{id: id}
	f.baseFunc.init(v)
	return f
}

func (f *SandboxDownloadResultsFunc) Do(ctx context.Context) (io.ReadCloser, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return f.Response, nil
}

func (f *SandboxDownloadResultsFunc) Store(ctx context.Context, filePath string) error {
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

func (s *SandboxDownloadResultsFunc) url() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s/report", s.id)
}

func (f *SandboxDownloadResultsFunc) responseStruct() any {
	return nil
}

func (s *SandboxDownloadResultsFunc) responseBody(body io.ReadCloser) {
	s.Response = body
}
