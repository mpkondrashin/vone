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

	SandboxSubmissionsResponse struct {
		Items    []ListSubmissionsItem `json:"items"`
		NextLink string                `json:"nextLink"`
	}
)

type SandboxSubmissionsFunc struct {
	BaseFunc
	Response SandboxSubmissionsResponse
}

func (v *VOne) SandboxListSubmissions() *SandboxSubmissionsFunc {
	f := &SandboxSubmissionsFunc{}
	f.BaseFunc.Init(v)
	//log.Println("AAA", f)
	return f
}

func (f *SandboxSubmissionsFunc) StartDateTime(t time.Time) *SandboxSubmissionsFunc {
	f.SetParameter("startDateTime", t.Format(timeFormat))
	return f
}

func (f *SandboxSubmissionsFunc) EndDateTime(t time.Time) *SandboxSubmissionsFunc {
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

func (f *SandboxSubmissionsFunc) DateTimeTarget(t DateTimeTarget) *SandboxSubmissionsFunc {
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

func (f *SandboxSubmissionsFunc) OrderBy(t DateTimeTarget, o Order) *SandboxSubmissionsFunc {
	f.SetParameter("orderBy", strings.Join([]string{t.String(), o.String()}, " "))
	//log.Println("BBB", f)
	return f
}

func (f *SandboxSubmissionsFunc) Filter(s string) *SandboxSubmissionsFunc {
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

func (f *SandboxSubmissionsFunc) Top(t Top) *SandboxSubmissionsFunc {
	f.SetParameter("top", t.String())
	return f
}

func (f *SandboxSubmissionsFunc) Do() (*SandboxSubmissionsResponse, error) {
	if err := f.vone.Call(f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *SandboxSubmissionsFunc) Method() string {
	return "GET"
}

func (*SandboxSubmissionsFunc) URL() string {
	return "/v3.0/sandbox/tasks"
}

func (f *SandboxSubmissionsFunc) URI() string {
	return f.Response.NextLink
}

func (f *SandboxSubmissionsFunc) ResponseStruct() any {
	return &f.Response
}

func (f *SandboxSubmissionsFunc) Next() (*SandboxSubmissionsResponse, error) {
	if f.Response.NextLink == "" {
		return nil, io.EOF
	}
	return f.Do()
}

func (f *SandboxSubmissionsFunc) IterateListSubmissions(callback func(*ListSubmissionsItem) error) error {
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
