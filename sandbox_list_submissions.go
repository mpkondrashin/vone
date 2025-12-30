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
	ListSubmissionsItem struct {
		ID                 string        `json:"id"`
		Action             Action        `json:"action"`
		Status             Status        `json:"status"`
		Error              Error         `json:"error"`
		CreatedDateTime    VisionOneTime `json:"createdDateTime"`
		LastActionDateTime VisionOneTime `json:"lastActionDateTime"`
		ResourceLocation   string        `json:"resourceLocation"`
		IsCached           bool          `json:"isCached"`
		Digest             Digest        `json:"digest"`
		Arguments          string        `json:"arguments"`
	}

	SandboxSubmissionsResponse struct {
		Items    []ListSubmissionsItem `json:"items"`
		NextLink string                `json:"nextLink"`
	}
)

type SandboxSubmissionsFunc struct {
	baseFunc
	Response SandboxSubmissionsResponse
}

func (v *VOne) SandboxListSubmissions() *SandboxSubmissionsFunc {
	f := &SandboxSubmissionsFunc{}
	f.baseFunc.init(v)
	return f
}

func (f *SandboxSubmissionsFunc) StartDateTime(t time.Time) *SandboxSubmissionsFunc {
	f.setParameter("startDateTime", t.Format(timeFormat))
	return f
}

func (f *SandboxSubmissionsFunc) EndDateTime(t time.Time) *SandboxSubmissionsFunc {
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

func (f *SandboxSubmissionsFunc) DateTimeTarget(t DateTimeTarget) *SandboxSubmissionsFunc {
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

func (f *SandboxSubmissionsFunc) OrderBy(t DateTimeTarget, o Order) *SandboxSubmissionsFunc {
	f.setParameter("orderBy", strings.Join([]string{t.String(), o.String()}, " "))
	return f
}

func (f *SandboxSubmissionsFunc) Filter(s string) *SandboxSubmissionsFunc {
	f.setParameter("filter", s)
	return f
}

func (f *SandboxSubmissionsFunc) Top(t Top) *SandboxSubmissionsFunc {
	f.setParameter("top", t.String())
	return f
}

func (f *SandboxSubmissionsFunc) Do(ctx context.Context) (*SandboxSubmissionsResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *SandboxSubmissionsFunc) Method() string {
	return methodGet
}

func (*SandboxSubmissionsFunc) url() string {
	return "/v3.0/sandbox/tasks"
}

func (f *SandboxSubmissionsFunc) uri() string {
	return f.Response.NextLink
}

func (f *SandboxSubmissionsFunc) responseStruct() any {
	return &f.Response
}

func (f *SandboxSubmissionsFunc) Next(ctx context.Context) (*SandboxSubmissionsResponse, error) {
	if f.Response.NextLink == "" {
		return nil, io.EOF
	}
	return f.Do(ctx)
}

func (f *SandboxSubmissionsFunc) IterateListSubmissions(ctx context.Context, callback func(*ListSubmissionsItem) error) error {
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
func (f *SandboxSubmissionsFunc) Range(ctx context.Context) iter.Seq2[*ListSubmissionsItem, error] {
	return func(yield func(*ListSubmissionsItem, error) bool) {
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
