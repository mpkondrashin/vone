package vone

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type SubmitFileToSandboxResponse struct {
	ID     string `json:"id"`
	Digest struct {
		MD5    string `json:"md5"`
		SHA1   string `json:"sha1"`
		SHA256 string `json:"sha256"`
	} `json:"digest"`
	Arguments string `json:"arguments"`
}

type SubmitFileToSandboxFunc struct {
	BaseFunc
	filePath            string
	Request             io.Reader
	formDataContentType string
	Response            SubmitFileToSandboxResponse
}

var _ Func = &SubmitFileToSandboxFunc{}

func (v *VOne) SubmitFileToSandbox(filePath string) (*SubmitFileToSandboxResponse, error) {
	f, err := NewSubmitFileToSandboxFunc(filePath)
	if err != nil {
		return nil, fmt.Errorf("NewSubmitFileToSandboxFunc: %w", err)
	}
	if err := v.Call(f); err != nil {
		return nil, fmt.Errorf("Call: %w", err)
	}
	return &f.Response, nil
}

func NewSubmitFileToSandboxFunc(filePath string) (*SubmitFileToSandboxFunc, error) {
	f := &SubmitFileToSandboxFunc{
		filePath: filePath,
	}
	f.BaseFunc.Init()
	var data bytes.Buffer
	writer := multipart.NewWriter(&data)
	w, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(w, reader); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	f.Request = &data
	f.formDataContentType = writer.FormDataContentType()
	return f, nil
}

func (s *SubmitFileToSandboxFunc) SetDocumentPassword(documentPassword string) *SubmitFileToSandboxFunc {
	s.SetParameter("documentPassword", documentPassword)
	return s
}

func (s *SubmitFileToSandboxFunc) SetArchivePassword(archivePassword string) *SubmitFileToSandboxFunc {
	s.SetParameter("archivePassword", archivePassword)
	return s
}

func (s *SubmitFileToSandboxFunc) SetArguments(arguments string) *SubmitFileToSandboxFunc {
	s.SetParameter("arguments", arguments)
	return s
}

func (f *SubmitFileToSandboxFunc) Method() string {
	return "POST"
}

func (s *SubmitFileToSandboxFunc) URL() string {
	return "/v3.0/sandbox/files/analyze"
}

func (f *SubmitFileToSandboxFunc) RequestBody() io.Reader {
	return f.Request
}

func (f *SubmitFileToSandboxFunc) ContentType() string {
	return f.formDataContentType
}

func (f *SubmitFileToSandboxFunc) ResponseStruct() any {
	return &f.Response
}

/*
func (s *SubmitFileToSandboxFunc) Do() (*SubmitFileToSandboxData, error) {
	url := "/v3.0/sandbox/files/analyze"
	var data bytes.Buffer
	writer := multipart.NewWriter(&data)
	w, err := writer.CreateFormFile("file", "file")
	if err != nil {
		return nil, err
	}
	reader, err := os.Open(s.filePath)
	/// err!!!
	if _, err := io.Copy(w, reader); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	body, err := s.Perform("POST", url, &data, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(body)
	var respData SubmitFileToSandboxData
	if err := decoder.Decode(&respData); err != nil && err != io.EOF {
		return nil, fmt.Errorf("response error: %w", err)
	}
	return &respData, nil
}
*/
/*
func (v *VOne) SubmitFileToSandboxZ(filePath, fileName, documentPassword, archivePassword, arguments string) (*SubmitFileToSandboxData, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return v.SubmitFileToSandboxAsReader(f, fileName, documentPassword, archivePassword, arguments)
}
*/
