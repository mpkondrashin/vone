/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_submit_file.go - upload file for analysis
*/

package vone

import (
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

// sandboxSubmitFileFunc - function to submit file to sandbox
type sandboxSubmitFileFunc struct {
	baseFunc
	//filePath            string
	request             io.Reader
	formDataContentType string
	response            SandboxSubmitFileResponse
	responseHeaders     SandboxSubmitFileResponseHeaders
}

var _ vOneFunc = &sandboxSubmitFileFunc{}

// SandboxSubmitFile - return new submit to sandbox file
func (v *VOne) SandboxSubmitFile() *sandboxSubmitFileFunc {
	f := &sandboxSubmitFileFunc{}
	f.baseFunc.init(v)
	return f
}

func (f *sandboxSubmitFileFunc) SetFilePath(filePath string) error {
	return f.SetFilePathAndName(filePath, filepath.Base(filePath))
}

func (f *sandboxSubmitFileFunc) SetFilePathAndName(filePath, fileName string) error {
	reader, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer reader.Close()
	return f.SetReader(reader, fileName)
}

func (f *sandboxSubmitFileFunc) SetReader(reader io.Reader, fileName string) error {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()

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

func (s *sandboxSubmitFileFunc) SetDocumentPassword(documentPassword string) *sandboxSubmitFileFunc {
	s.setParameter("documentPassword", documentPassword)
	return s
}

func (s *sandboxSubmitFileFunc) SetArchivePassword(archivePassword string) *sandboxSubmitFileFunc {
	s.setParameter("archivePassword", archivePassword)
	return s
}

func (s *sandboxSubmitFileFunc) SetArguments(arguments string) *sandboxSubmitFileFunc {
	s.setParameter("arguments", arguments)
	return s
}

func (f *sandboxSubmitFileFunc) Do(ctx context.Context) (*SandboxSubmitFileResponse, *SandboxSubmitFileResponseHeaders, error) {
	if f.vone.mockup != nil {
		return f.vone.mockup.SubmitFile(f)
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, nil, err
	}
	return &f.response, &f.responseHeaders, nil
}

func (f *sandboxSubmitFileFunc) method() string {
	return methodPost
}

func (s *sandboxSubmitFileFunc) url() string {
	return "/v3.0/sandbox/files/analyze"
}

func (f *sandboxSubmitFileFunc) requestBody() io.Reader {
	return f.request
}

func (f *sandboxSubmitFileFunc) contentType() string {
	return f.formDataContentType
}

func (f *sandboxSubmitFileFunc) responseStruct() any {
	return &f.response
}

func (f *sandboxSubmitFileFunc) responseHeader() any {
	return &f.responseHeaders
}
