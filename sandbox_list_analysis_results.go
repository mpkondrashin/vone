/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_list_submissions.go - list all submissions
*/

package vone

import (
	"context"
	"io"
	"iter"
	"time"
)

type (
	SandboxListAnalysisResultResponse struct {
		Items    []SandboxAnalysisResultsResponse `json:"items"`
		NextLink string                           `json:"nextLink"`
	}
)

type SandboxListAnalysisResultsFunc struct {
	baseFunc
	Response SandboxListAnalysisResultResponse
}

func (v *VOne) SandboxListAnalysisResults() *SandboxListAnalysisResultsFunc {
	f := &SandboxListAnalysisResultsFunc{}
	f.baseFunc.init(v)
	return f
}

func (f *SandboxListAnalysisResultsFunc) StartDateTime(t time.Time) *SandboxListAnalysisResultsFunc {
	f.setParameter("startDateTime", t.Format(timeFormatZ))
	return f
}

func (f *SandboxListAnalysisResultsFunc) EndDateTime(t time.Time) *SandboxListAnalysisResultsFunc {
	f.setParameter("endDateTime", t.Format(timeFormatZ))
	return f
}

func (f *SandboxListAnalysisResultsFunc) OrderBy(orderBy string) *SandboxListAnalysisResultsFunc {
	f.setParameter("orderBy", orderBy)
	return f
}

func (f *SandboxListAnalysisResultsFunc) Filter(s string) *SandboxListAnalysisResultsFunc {
	f.setParameter("filter", s)
	return f
}

func (f *SandboxListAnalysisResultsFunc) Top(t Top) *SandboxListAnalysisResultsFunc {
	f.setParameter("top", t.String())
	return f
}

func (f *SandboxListAnalysisResultsFunc) Do(ctx context.Context) (*SandboxListAnalysisResultResponse, error) {
	if f.vone.mockup != nil {
		return f.vone.mockup.ListAnalysisResults(f)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *SandboxListAnalysisResultsFunc) Method() string {
	return methodGet
}

func (*SandboxListAnalysisResultsFunc) url() string {
	return "/v3.0/sandbox/analysisResults"
}

func (f *SandboxListAnalysisResultsFunc) uri() string {
	return f.Response.NextLink
}

func (f *SandboxListAnalysisResultsFunc) responseStruct() any {
	return &f.Response
}

func (f *SandboxListAnalysisResultsFunc) Next(ctx context.Context) (*SandboxListAnalysisResultResponse, error) {
	if f.Response.NextLink == "" {
		return nil, io.EOF
	}
	return f.Do(ctx)
}

func (f *SandboxListAnalysisResultsFunc) IterateListSubmissions(ctx context.Context, callback func(*SandboxAnalysisResultsResponse) error) error {
	for {
		if err := f.vone.call(ctx, f); err != nil {
			return err
		}
		for i := range f.Response.Items {
			if err := callback(&f.Response.Items[i]); err != nil {
				return err
			}
		}
		if f.Response.NextLink == "" {
			return nil
		}
	}
}

// Range - iterator for all submissions (go 1.23 and later)
func (f *SandboxListAnalysisResultsFunc) Range(ctx context.Context) iter.Seq2[*SandboxAnalysisResultsResponse, error] {
	return func(yield func(*SandboxAnalysisResultsResponse, error) bool) {
		for {
			if err := f.vone.call(ctx, f); err != nil {
				yield(nil, err)
				return
			}
			for i := range f.Response.Items {
				if !yield(&f.Response.Items[i], nil) {
					return
				}
			}
			if f.Response.NextLink == "" {
				return
			}
		}
	}
}
