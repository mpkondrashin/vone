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
	BaseFunc
	id       string
	Response io.ReadCloser
}

func (v *VOne) SandboxDownloadResults(id string) *SandboxDownloadResultsFunc {
	f := &SandboxDownloadResultsFunc{id: id}
	f.BaseFunc.Init(v)
	return f
}

func (f *SandboxDownloadResultsFunc) Do(ctx context.Context) (io.ReadCloser, error) {
	if err := f.vone.Call(ctx, f); err != nil {
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

func (s *SandboxDownloadResultsFunc) URL() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s/report", s.id)
}

func (f *SandboxDownloadResultsFunc) ResponseStruct() any {
	return nil
}

func (s *SandboxDownloadResultsFunc) ResponseBody(body io.ReadCloser) {
	s.Response = body
}
