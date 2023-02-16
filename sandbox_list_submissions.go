package vone

import (
	"io"
	"strings"
	"time"
)

type (
	ListSubmissionsItem struct {
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
	}

	ListSubmissionsResponse struct {
		Items    []ListSubmissionsItem `json:"items"`
		NextLink string                `json:"nextLink"`
	}
)

type ListSubmissionsFunc struct {
	BaseFunc
	Response ListSubmissionsResponse
}

var _ Func = &ListSubmissionsFunc{}

func (v *VOne) SandboxListSubmissions() *ListSubmissionsFunc {
	f := &ListSubmissionsFunc{}
	f.BaseFunc.Init(v)
	//log.Println("AAA", f)
	return f
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
	//log.Println("BBB", f)
	return f
}

func (f *ListSubmissionsFunc) Filter(s string) *ListSubmissionsFunc {
	f.SetParameter("filter", s)
	return f
}

type Top int

const (
	Top50 Top = iota
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

func (f *ListSubmissionsFunc) Do() (*ListSubmissionsResponse, error) {
	if err := f.vone.Call(f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *ListSubmissionsFunc) Method() string {
	return "GET"
}

func (*ListSubmissionsFunc) URL() string {
	return "/v3.0/sandbox/tasks"
}

func (f *ListSubmissionsFunc) URI() string {
	return f.Response.NextLink
}

func (f *ListSubmissionsFunc) ResponseStruct() any {
	return &f.Response
}

func (f *ListSubmissionsFunc) Next() (*ListSubmissionsResponse, error) {
	if f.Response.NextLink == "" {
		return nil, io.EOF
	}
	return f.Do()
}

func (f *ListSubmissionsFunc) IterateListSubmissions(callback func(*ListSubmissionsItem) error) error {
	for {
		if err := f.vone.Call(f); err != nil {
			return err
		}
		for _, r := range f.Response.Items {
			if err := callback(&r); err != nil {
				return err
			}
		}
		if f.Response.NextLink == "" {
			return nil
		}
	}
}
