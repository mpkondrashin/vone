package vone

import (
	"fmt"
	"strings"
	"time"
)

type (
	ListSubmissionsResponse struct {
		Items []struct {
			ID     string `json:"id"`
			Action string `json:"action"`
			Status string `json:"status"`
			Error  struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
			CreatedDateTime    time.Time `json:"createdDateTime"`
			LastActionDateTime time.Time `json:"lastActionDateTime"`
			ResourceLocation   string    `json:"resourceLocation"`
			IsCached           bool      `json:"isCached"`
			Digest             struct {
				MD5    string `json:"md5"`
				SHA1   string `json:"sha1"`
				SHA256 string `json:"sha256"`
			} `json:"digest"`
			Arguments string `json:"arguments"`
		} `json:"items"`
		NextLink string `json:"nextLink"`
	}
)

type ListSubmissionsFunc struct {
	BaseFunc
	Response *ListSubmissionsResponse
}

var _ Func = &ListSubmissionsFunc{}

//func (v *VOne) ListSubmissions(startDateTime, endDateTime time.Time, dateTimeTarget DateTimeTarget, orderTarget DateTimeTarget, orderBy Order, filter string, top Top) (*ListSubmissionsResponse, error) {
func (v *VOne) ListSubmissions() (*ListSubmissionsResponse, error) {
	f, err := NewListSubmissionsFunc()
	if err != nil {
		return nil, fmt.Errorf("NewSubmitURLsToSandboxFunc: %w", err)
	}
	//	f.StartDateTime(startDateTime).EndDateTime(endDateTime).DateTimeTarget(dateTimeTarget)
	//	f.OrderBy(orderTarget, orderBy)
	//	f.Filter(filter)
	//	f.Top(top)
	if err := v.Call(f); err != nil {
		return nil, fmt.Errorf("Call: %w", err)
	}
	return f.Response, nil
}

func NewListSubmissionsFunc() (*ListSubmissionsFunc, error) {
	f := &ListSubmissionsFunc{}
	f.BaseFunc.Init()
	return f, nil
}

func (f *ListSubmissionsFunc) StartDateTime(t time.Time) *ListSubmissionsFunc {
	f.SetParameter("startDateTime", t.Format(timeFormat))
	return f
}

func (f *ListSubmissionsFunc) EndDateTime(t time.Time) *ListSubmissionsFunc {
	f.SetParameter("endDateTime", t.Format(timeFormat))
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

func (f *ListSubmissionsFunc) DateTimeTarget(t DateTimeTarget) *ListSubmissionsFunc {
	f.SetParameter("sateTimeTarget", t.String())
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

func (f *ListSubmissionsFunc) OrderBy(t DateTimeTarget, o Order) *ListSubmissionsFunc {
	f.SetParameter("orderBy", strings.Join([]string{t.String(), o.String()}, " "))
	return f
}

func (f *ListSubmissionsFunc) Filter(s string) *ListSubmissionsFunc {
	f.SetParameter("filter", s)
	return f
}

type Top int

const (
	Top50 Order = iota
	Top100
	Top200
)

func (t Top) String() string {
	return [...]string{"50", "100", "200"}[t]
}

func (f *ListSubmissionsFunc) Top(t Top) *ListSubmissionsFunc {
	f.SetParameter("top", t.String())
	return f
}

func (f *ListSubmissionsFunc) Method() string {
	return "GET"
}

func (*ListSubmissionsFunc) URL() string {
	return "/v3.0/sandbox/tasks"
}

func (f *ListSubmissionsFunc) ResponseStruct() any {
	return &f.Response
}

type ListSubmissionsNextFunc struct {
	BaseFunc
	previousResponse *ListSubmissionsResponse
	Response         *ListSubmissionsResponse
}

func NewListSubmissionsNextFunc(previousResponse *ListSubmissionsResponse) *ListSubmissionsNextFunc {
	return &ListSubmissionsNextFunc{
		previousResponse: previousResponse,
	}
}

func (f *ListSubmissionsNextFunc) URI() string {
	return f.previousResponse.NextLink
}

func (v *VOne) ListSubmissionsNext(r *ListSubmissionsResponse) (*ListSubmissionsResponse, error) {
	f := NewListSubmissionsNextFunc(r)
	if err := v.Call(f); err != nil {
		return nil, fmt.Errorf("Call: %w", err)
	}
	return f.Response, nil
}

//func (v *VOne) ListSubmissions(startDateTime, endDateTime time.Time, dateTimeTarget DateTimeTarget, orderTarget DateTimeTarget, orderBy Order, filter string, top Top) (*ListSubmissionsResponse, error) {
func (v *VOne) LisdddtSubmissions() (*ListSubmissionsResponse, error) {
	f, err := NewListSubmissionsFunc()
	if err != nil {
		return nil, fmt.Errorf("NewSubmitURLsToSandboxFunc: %w", err)
	}
	//	f.StartDateTime(startDateTime).EndDateTime(endDateTime).DateTimeTarget(dateTimeTarget)
	//	f.OrderBy(orderTarget, orderBy)
	//	f.Filter(filter)
	//	f.Top(top)
	if err := v.Call(f); err != nil {
		return nil, fmt.Errorf("Call: %w", err)
	}
	return f.Response, nil
}
