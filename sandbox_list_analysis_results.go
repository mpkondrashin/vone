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

type sandboxListAnalysisResultsRequest struct {
	baseFunc
	response SandboxListAnalysisResultResponse
}

func (v *VOne) SandboxListAnalysisResults() *sandboxListAnalysisResultsRequest {
	f := &sandboxListAnalysisResultsRequest{}
	f.baseFunc.init(v)
	return f
}

func (f *sandboxListAnalysisResultsRequest) StartDateTime(t time.Time) *sandboxListAnalysisResultsRequest {
	f.setParameter("startDateTime", t.Format(timeFormatZ))
	return f
}

func (f *sandboxListAnalysisResultsRequest) EndDateTime(t time.Time) *sandboxListAnalysisResultsRequest {
	f.setParameter("endDateTime", t.Format(timeFormatZ))
	return f
}

func (f *sandboxListAnalysisResultsRequest) OrderBy(orderBy string) *sandboxListAnalysisResultsRequest {
	f.setParameter("orderBy", orderBy)
	return f
}

func (f *sandboxListAnalysisResultsRequest) Filter(s string) *sandboxListAnalysisResultsRequest {
	f.setParameter("filter", s)
	return f
}

func (f *sandboxListAnalysisResultsRequest) Top(t Top) *sandboxListAnalysisResultsRequest {
	f.setParameter("top", t.String())
	return f
}

func (f *sandboxListAnalysisResultsRequest) Do(ctx context.Context) (*SandboxListAnalysisResultResponse, error) {
	if f.vone.mockup != nil {
		return f.vone.mockup.ListAnalysisResults(f)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.response, nil
}

func (f *sandboxListAnalysisResultsRequest) Method() string {
	return methodGet
}

func (*sandboxListAnalysisResultsRequest) url() string {
	return "/v3.0/sandbox/analysisResults"
}

func (f *sandboxListAnalysisResultsRequest) uri() string {
	return f.response.NextLink
}

func (f *sandboxListAnalysisResultsRequest) responseStruct() any {
	return &f.response
}

func (f *sandboxListAnalysisResultsRequest) Next(ctx context.Context) (*SandboxListAnalysisResultResponse, error) {
	if f.response.NextLink == "" {
		return nil, io.EOF
	}
	return f.Do(ctx)
}

func (f *sandboxListAnalysisResultsRequest) IterateListSubmissions(ctx context.Context, callback func(*SandboxAnalysisResultsResponse) error) error {
	if f.vone.mockup != nil {
		response, err := f.vone.mockup.ListAnalysisResults(f)
		if err != nil {
			return err
		}
		for i := range response.Items {
			if err := callback(&response.Items[i]); err != nil {
				return err
			}
		}
		return nil
	}
	for {
		if err := f.vone.call(ctx, f); err != nil {
			return err
		}
		for i := range f.response.Items {
			if err := callback(&f.response.Items[i]); err != nil {
				return err
			}
		}
		if f.response.NextLink == "" {
			return nil
		}
	}
}

// Range - iterator for all submissions (go 1.23 and later)
func (f *sandboxListAnalysisResultsRequest) Range(ctx context.Context) iter.Seq2[*SandboxAnalysisResultsResponse, error] {
	return func(yield func(*SandboxAnalysisResultsResponse, error) bool) {
		if f.vone.mockup != nil {
			response, err := f.vone.mockup.ListAnalysisResults(f)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range response.Items {
				if !yield(&response.Items[i], nil) {
					return
				}
			}
			return
		}
		for {
			if err := f.vone.call(ctx, f); err != nil {
				yield(nil, err)
				return
			}
			for i := range f.response.Items {
				if !yield(&f.response.Items[i], nil) {
					return
				}
			}
			if f.response.NextLink == "" {
				return
			}
		}
	}
}
