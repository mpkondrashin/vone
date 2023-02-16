package vone

import (
	"bytes"
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
	//filePath            string
	Request             io.Reader
	formDataContentType string
	Response            SubmitFileToSandboxResponse
}

var _ Func = &SubmitFileToSandboxFunc{}

/*
func (v *VOne) SubmitFileToSandbox(filePath string) (*SubmitFileToSandboxResponse, error) {
	f, err := SandboxSubmitFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("SandboxSubmitFile: %w", err)
	}
	if err := v.Call(f); err != nil {
		return nil, fmt.Errorf("Call: %w", err)
	}
	return &f.Response, nil
}
*/
func (v *VOne) SandboxSubmitFile() *SubmitFileToSandboxFunc {
	f := &SubmitFileToSandboxFunc{}
	f.BaseFunc.Init(v)
	return f
}

func (f *SubmitFileToSandboxFunc) SetFileName(fileName string) (*SubmitFileToSandboxFunc, error) {
	//	f.filePath = fileName
	reader, err := os.Open(fileName)
	if err != nil {
		return f, err
	}
	defer reader.Close()
	return f, f.SetReader(reader, filepath.Base(fileName))
}

func (f *SubmitFileToSandboxFunc) SetReader(reader io.Reader, fileName string) error {
	var data bytes.Buffer
	writer := multipart.NewWriter(&data)
	w, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, reader); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	f.Request = &data
	f.formDataContentType = writer.FormDataContentType()
	return nil
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

func (f *SubmitFileToSandboxFunc) Do() (*SubmitFileToSandboxResponse, error) {
	if err := f.vone.Call(f); err != nil {
		return nil, err
	}
	return &f.Response, nil
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
