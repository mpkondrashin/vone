/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_submission_status.go - status of submission
*/

package vone

import (
	"fmt"
	"time"
)

// g o : g enerate stringer -type=Status
//go:generate enum -package=vone -type=Status -values=succeeded,running,failed

/*type Status int

const (
	Succeed Status = iota
	Running
	Failed
)

var ErrUnknown = errors.New("unknown")

func (s *Status) UnmarshalJSON(data []byte) error {
	m := map[string]Status{
		"succeed": Succeed,
		"running": Running,
		"failed":  Failed,
	}
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	result, ok := m[v]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknown, v)
	}
	*s = result
	return nil
}
*/
type (
	SandboxSubmissionStatusResponse struct {
		ID     string `json:"id"`
		Action string `json:"action"`
		Status Status `json:"status"`
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
)

type SandboxSubmissionStatusFunc struct {
	BaseFunc
	id       string
	Response SandboxSubmissionStatusResponse
}

func (v *VOne) SandboxSubmissionStatus(id string) *SandboxSubmissionStatusFunc {
	f := &SandboxSubmissionStatusFunc{
		id: id,
	}
	f.BaseFunc.Init(v)
	return f
}

func (f *SandboxSubmissionStatusFunc) Do() (*SandboxSubmissionStatusResponse, error) {
	if err := f.vone.Call(f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *SandboxSubmissionStatusFunc) URL() string {
	return fmt.Sprintf("/v3.0/sandbox/tasks/%s", f.id)
}

func (f *SandboxSubmissionStatusFunc) ResponseStruct() any {
	return &f.Response
}
