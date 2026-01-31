/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_submit_file.go - upload file for analysis
*/

package vone

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var ErrFileNotSet = errors.New("file not set")

// SandboxSubmitFileResponse - Submit file to sandbox response JSON format.
type SandboxSubmitFileResponse struct {
	ID     string `json:"id"`
	Digest struct {
		MD5    string `json:"md5"`
		SHA1   string `json:"sha1"`
		SHA256 string `json:"sha256"`
	} `json:"digest"`
	Arguments string `json:"arguments"`
}

// SandboxSubmitFileResponseHeaders contains response headers returned by sandbox file submission.
type SandboxSubmitFileResponseHeaders struct {
	OperationLocation        string `header:"Operation-Location"`
	SubmissionReserveCount   int    `header:"TMV1-Submission-Reserve-Count"`
	SubmissionRemainingCount int    `header:"TMV1-Submission-Remaining-Count"`
	SubmissionCount          int    `header:"TMV1-Submission-Count"`
	SubmissionExemptionCount int    `header:"TMV1-Submission-Exemption-Count"`
}

// sandboxSubmitFileRequest - function to submit file to sandbox
type sandboxSubmitFileRequest struct {
	baseFunc
	request             io.Reader
	formDataContentType string
	response            SandboxSubmitFileResponse
	responseHeaders     SandboxSubmitFileResponseHeaders
}

var _ vOneFunc = &sandboxSubmitFileRequest{}

// SandboxSubmitFile - return new submit to sandbox file
func (v *VOne) SandboxSubmitFile() *sandboxSubmitFileRequest {
	f := &sandboxSubmitFileRequest{}
	f.baseFunc.init(v)
	return f
}

func (f *sandboxSubmitFileRequest) SetFilePath(filePath string) error {
	return f.SetFilePathAndName(filePath, filepath.Base(filePath))
}

func (f *sandboxSubmitFileRequest) SetFilePathAndName(filePath, fileName string) error {
	reader, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer reader.Close()
	return f.SetReader(reader, fileName)
}

func (f *sandboxSubmitFileRequest) SetReader(reader io.Reader, fileName string) error {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()
		if c, ok := reader.(io.Closer); ok {
			defer c.Close()
		}
		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err := io.Copy(part, reader); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	f.request = pr
	f.formDataContentType = writer.FormDataContentType()
	return nil
}

func (s *sandboxSubmitFileRequest) SetDocumentPassword(documentPassword string) *sandboxSubmitFileRequest {
	s.setParameter("documentPassword", documentPassword)
	return s
}

func (s *sandboxSubmitFileRequest) SetArchivePassword(archivePassword string) *sandboxSubmitFileRequest {
	s.setParameter("archivePassword", archivePassword)
	return s
}

func (s *sandboxSubmitFileRequest) SetArguments(arguments string) *sandboxSubmitFileRequest {
	s.setParameter("arguments", arguments)
	return s
}

func (f *sandboxSubmitFileRequest) Do(ctx context.Context) (*SandboxSubmitFileResponse, *SandboxSubmitFileResponseHeaders, error) {
	if f.request == nil {
		return nil, nil, fmt.Errorf("submit file: %w", ErrFileNotSet)
	}
	if f.vone.mockup != nil {
		return f.vone.mockup.SubmitFile(f)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, nil, err
	}
	return &f.response, &f.responseHeaders, nil
}

func (f *sandboxSubmitFileRequest) method() string {
	return methodPost
}

func (s *sandboxSubmitFileRequest) url() string {
	return "/v3.0/sandbox/files/analyze"
}

func (f *sandboxSubmitFileRequest) requestBody() io.Reader {
	return f.request
}

func (f *sandboxSubmitFileRequest) contentType() string {
	return f.formDataContentType
}

func (f *sandboxSubmitFileRequest) responseStruct() any {
	return &f.response
}

func (f *sandboxSubmitFileRequest) responseHeader() any {
	return &f.responseHeaders
}
