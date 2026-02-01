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
	baseRequest
	response SandboxSubmissionsResponse
}

func (v *VOne) SandboxListSubmissions() *sandboxSubmissionsRequest {
	f := &sandboxSubmissionsRequest{}
	f.baseRequest.init(v)
	return f
}

func (f *sandboxSubmissionsRequest) StartDateTime(t time.Time) *sandboxSubmissionsRequest {
	f.setParameter("startDateTime", t.Format(timeFormatZ))
	return f
}

func (f *sandboxSubmissionsRequest) EndDateTime(t time.Time) *sandboxSubmissionsRequest {
	f.setParameter("endDateTime", t.Format(timeFormatZ))
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
	f.setParameter("dateTimeTarget", t.String())
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

func (*sandboxSubmissionsRequest) url() string {
	return "/v3.0/sandbox/tasks"
}

func (f *sandboxSubmissionsRequest) uri() string {
	return f.response.NextLink
}

func (f *sandboxSubmissionsRequest) responseStruct() any {
	return &f.response
}

func (f *sandboxSubmissionsRequest) nextLink() string {
	return f.response.NextLink
}

func (f *sandboxSubmissionsRequest) resetPagination() {
	f.response.NextLink = ""
}

func (f *sandboxSubmissionsRequest) Next(ctx context.Context) (*SandboxSubmissionsResponse, error) {
	if f.response.NextLink == "" {
		return nil, io.EOF
	}
	return f.Do(ctx)
}

/*
// Range - iterator for all submissions (go 1.23 and later)

	func (f *sandboxSubmissionsRequest) Range(ctx context.Context) iter.Seq2[*SandboxSubmissionStatusResponse, error] {
		return func(yield func(*SandboxSubmissionStatusResponse, error) bool) {
			if err := f.checkUsed(); err != nil {
				yield(nil, fmt.Errorf("submissions: %w", err))
				return
			}
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
*/
func (f *sandboxSubmissionsRequest) Paginator() *Paginator[
	SandboxSubmissionsResponse,
	SandboxSubmissionStatusResponse,
] {
	return NewPaginator(
		f,
		func(r *SandboxSubmissionsResponse) []SandboxSubmissionStatusResponse {
			return r.Items
		},
	)
}
