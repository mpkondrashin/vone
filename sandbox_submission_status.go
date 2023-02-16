package vone

import (
	"fmt"
	"time"
)

type (
	SubmissionStatusResponse struct {
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
)

type SubmissionStatusFunc struct {
	BaseFunc
	id       string
	Response SubmissionStatusResponse
}

func (v *VOne) SubmissionStatus(id string) *SubmissionStatusFunc {
	f := &SubmissionStatusFunc{
		id: id,
	}
	f.BaseFunc.Init(v)
	return f
}

func (f *SubmissionStatusFunc) Do() (*SubmissionStatusResponse, error) {
	if err := f.vone.Call(f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *SubmissionStatusFunc) URL() string {
	return fmt.Sprintf("/v3.0/sandbox/tasks/%s", f.id)
}

func (f *SubmissionStatusFunc) ResponseStruct() any {
	return &f.Response
}
