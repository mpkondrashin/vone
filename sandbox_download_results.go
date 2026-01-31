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

type sandboxDownloadResultsFunc struct {
	baseFunc
	id       string
	response io.ReadCloser
}

var _ vOneFunc = &sandboxDownloadResultsFunc{}

func (v *VOne) SandboxDownloadResults(id string) *sandboxDownloadResultsFunc {
	f := &sandboxDownloadResultsFunc{id: id}
	f.baseFunc.init(v)
	return f
}

func (f *sandboxDownloadResultsFunc) Do(ctx context.Context) (io.ReadCloser, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return f.response, nil
}

func (f *sandboxDownloadResultsFunc) Store(ctx context.Context, filePath string) error {
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

func (s *sandboxDownloadResultsFunc) url() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s/report", s.id)
}

func (f *sandboxDownloadResultsFunc) responseStruct() any {
	return nil
}

func (s *sandboxDownloadResultsFunc) responseBody(body io.ReadCloser) {
	s.response = body
}
