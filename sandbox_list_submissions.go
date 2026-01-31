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
	"strings"
	"time"
)

type (
	SandboxSubmissionsResponse struct {
		Items    []SandboxSubmissionStatusResponse `json:"items"`
		NextLink string                            `json:"nextLink"`
	}
)

type sandboxSubmissionsRequest struct {
	baseFunc
	response SandboxSubmissionsResponse
}

func (v *VOne) SandboxListSubmissions() *sandboxSubmissionsRequest {
	f := &sandboxSubmissionsRequest{}
	f.baseFunc.init(v)
	return f
}

func (f *sandboxSubmissionsRequest) StartDateTime(t time.Time) *sandboxSubmissionsRequest {
	f.setParameter("startDateTime", t.Format(timeFormat))
	return f
}

func (f *sandboxSubmissionsRequest) EndDateTime(t time.Time) *sandboxSubmissionsRequest {
	f.setParameter("endDateTime", t.Format(timeFormat))
	return f
}

type DateTimeTarget int

const (
	CreatedDateTime DateTimeTarget = iota
	LastActionDateTime
)

func (t DateTimeTarget) String() string {
	return [...]string{"createdDateTime", "lastActionDateTime"}[t]
}

func (f *sandboxSubmissionsRequest) DateTimeTarget(t DateTimeTarget) *sandboxSubmissionsRequest {
	f.setParameter("sateTimeTarget", t.String())
	return f
}

type Order int

const (
	Desc Order = iota
	Asc
)

func (o Order) String() string {
	return [...]string{"desc", "asc"}[o]
}

func (f *sandboxSubmissionsRequest) OrderBy(t DateTimeTarget, o Order) *sandboxSubmissionsRequest {
	f.setParameter("orderBy", strings.Join([]string{t.String(), o.String()}, " "))
	return f
}

func (f *sandboxSubmissionsRequest) Filter(s string) *sandboxSubmissionsRequest {
	f.setParameter("filter", s)
	return f
}

func (f *sandboxSubmissionsRequest) Top(t Top) *sandboxSubmissionsRequest {
	f.setParameter("top", t.String())
	return f
}

func (f *sandboxSubmissionsRequest) Do(ctx context.Context) (*SandboxSubmissionsResponse, error) {
	if f.vone.mockup != nil {
		return f.vone.mockup.ListSubmissions(f)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.response, nil
}

func (f *sandboxSubmissionsRequest) Method() string {
	return methodGet
}

func (*sandboxSubmissionsRequest) url() string {
	return "/v3.0/sandbox/tasks"
}

func (f *sandboxSubmissionsRequest) uri() string {
	return f.response.NextLink
}

func (f *sandboxSubmissionsRequest) responseStruct() any {
	return &f.response
}

func (f *sandboxSubmissionsRequest) Next(ctx context.Context) (*SandboxSubmissionsResponse, error) {
	if f.response.NextLink == "" {
		return nil, io.EOF
	}
	return f.Do(ctx)
}

func (f *sandboxSubmissionsRequest) IterateListSubmissions(ctx context.Context, callback func(*SandboxSubmissionStatusResponse) error) error {
	if f.vone.mockup != nil {
		response, err := f.vone.mockup.ListSubmissions(f)
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
func (f *sandboxSubmissionsRequest) Range(ctx context.Context) iter.Seq2[*SandboxSubmissionStatusResponse, error] {
	return func(yield func(*SandboxSubmissionStatusResponse, error) bool) {
		if f.vone.mockup != nil {
			response, err := f.vone.mockup.ListSubmissions(f)
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
