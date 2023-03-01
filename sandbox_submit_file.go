/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_submit_file.go - upload file for analysis
*/

package vone

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SandboxSubmitFileResponse - Submit file or url to sanbox response JSON format
type SandboxSubmitFileResponse struct {
	ID     string `json:"id"`
	Digest struct {
		MD5    string `json:"md5"`
		SHA1   string `json:"sha1"`
		SHA256 string `json:"sha256"`
	} `json:"digest"`
	Arguments string `json:"arguments"`
}

// SandboxSubmitFileResponse - struct used to return value of HTTP headers from file or url submit to sanbox
type SandboxSubmitFileResponseHeaders struct {
	OperationLocation        string `header:"Operation-Location"`
	SubmissionReserveCount   int    `header:"TMV1-Submission-Reserve-Count"`
	SubmissionRemainingCount int    `header:"TMV1-Submission-Remaining-Count"`
	SubmissionCount          int    `header:"TMV1-Submission-Count"`
	SubmissionExemptionCount int    `header:"TMV1-Submission-Exemption-Count"`
}

// SandboxSubmitFileToSandboxFunc - function to submit file to sandbox
type SandboxSubmitFileToSandboxFunc struct {
	BaseFunc
	//filePath            string
	Request             io.Reader
	formDataContentType string
	Response            SandboxSubmitFileResponse
	ResponseHeaders     SandboxSubmitFileResponseHeaders
}

var _ Func = &SandboxSubmitFileToSandboxFunc{}

// SandboxSubmitFile - return new submit to sandbox file
func (v *VOne) SandboxSubmitFile() *SandboxSubmitFileToSandboxFunc {
	f := &SandboxSubmitFileToSandboxFunc{}
	f.BaseFunc.Init(v)
	return f
}

func (f *SandboxSubmitFileToSandboxFunc) SetFilePath(filePath string) (*SandboxSubmitFileToSandboxFunc, error) {
	return f.SetFilePathAndName(filePath, filepath.Base(filePath))
}

func (f *SandboxSubmitFileToSandboxFunc) SetFilePathAndName(filePath, fileName string) (*SandboxSubmitFileToSandboxFunc, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		return f, err
	}
	defer reader.Close()
	return f.SetReader(reader, fileName)
}

func (f *SandboxSubmitFileToSandboxFunc) SetReader(reader io.Reader, fileName string) (*SandboxSubmitFileToSandboxFunc, error) {
	var data bytes.Buffer
	writer := multipart.NewWriter(&data)
	w, err := writer.CreateFormFile("file", fileName)
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

func (s *SandboxSubmitFileToSandboxFunc) SetDocumentPassword(documentPassword string) *SandboxSubmitFileToSandboxFunc {
	s.SetParameter("documentPassword", documentPassword)
	return s
}

func (s *SandboxSubmitFileToSandboxFunc) SetArchivePassword(archivePassword string) *SandboxSubmitFileToSandboxFunc {
	s.SetParameter("archivePassword", archivePassword)
	return s
}

func (s *SandboxSubmitFileToSandboxFunc) SetArguments(arguments string) *SandboxSubmitFileToSandboxFunc {
	s.SetParameter("arguments", arguments)
	return s
}

func (f *SandboxSubmitFileToSandboxFunc) Do(ctx context.Context) (*SandboxSubmitFileResponse, *SandboxSubmitFileResponseHeaders, error) {
	if err := f.vone.Call(ctx, f); err != nil {
		return nil, nil, err
	}
	return &f.Response, &f.ResponseHeaders, nil
}

func (f *SandboxSubmitFileToSandboxFunc) Method() string {
	return POST
}

func (s *SandboxSubmitFileToSandboxFunc) URL() string {
	return "/v3.0/sandbox/files/analyze"
}

func (f *SandboxSubmitFileToSandboxFunc) RequestBody() io.Reader {
	return f.Request
}

func (f *SandboxSubmitFileToSandboxFunc) ContentType() string {
	return f.formDataContentType
}

func (f *SandboxSubmitFileToSandboxFunc) ResponseStruct() any {
	return &f.Response
}

func (f *SandboxSubmitFileToSandboxFunc) ResponseHeader() any {
	return &f.ResponseHeaders
}
