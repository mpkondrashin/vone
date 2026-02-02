package vone

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SandboxMockup interface {
	SubmitFile(f *sandboxSubmitFileRequest) (*SandboxSubmitFileResponse, *SandboxSubmitFileResponseHeaders, error)
	SubmissionStatus(f *sandboxSubmissionStatusRequest) (*SandboxSubmissionStatusResponse, error)
	AnalysisResults(f *sandboxAnalysisResultsRequest) (*SandboxAnalysisResultsResponseItem, error)
	ListSubmissions(f *sandboxSubmissionsRequest) (*SandboxSubmissionsResponse, error)
	ListAnalysisResults(f *sandboxListAnalysisResultsRequest) (*SandboxListAnalysisResultResponse, error)
}

type SubmissionState int

const (
	state SubmissionState = iota
)

type (
	submission struct {
		state            SubmissionState
		analysisResult   SandboxAnalysisResultsResponseItem
		submissionStatus SandboxSubmissionStatusResponse
	}

	SandboxMockupRAM struct {
		submissions map[string]*submission
		logger      *log.Logger
	}
)

func NewSandboxMockup() *SandboxMockupRAM {
	return &SandboxMockupRAM{
		submissions: make(map[string]*submission),
		logger:      log.New(io.Discard, "", 0),
	}
}

func (s *SandboxMockupRAM) EnableFileLogging(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("mockup logging: %v", err)
	}
	s.logger = log.New(
		f,
		"",
		log.LstdFlags,
	)
	return nil
}

func (sm *SandboxMockupRAM) extractFirstPart(body []byte, contentType string) ([]byte, error) {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}
	boundary := params["boundary"]
	if boundary == "" {
		return nil, errors.New("no boundary")
	}
	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	part, err := reader.NextPart()
	if err == io.EOF {
		return nil, errors.New("file content not found")
	}
	if err != nil {
		return nil, err
	}
	return io.ReadAll(part)
}

func (sm *SandboxMockupRAM) SubmitFile(f *sandboxSubmitFileRequest) (*SandboxSubmitFileResponse, *SandboxSubmitFileResponseHeaders, error) {
	sm.logger.Println("SubmitFile")
	data, err := io.ReadAll(f.request)
	if err != nil {
		return nil, nil, fmt.Errorf("submit file: %v", err)
	}
	re := regexp.MustCompile(`\s+`)
	strippedData := re.ReplaceAllString(string(data), "")
	sm.logger.Printf("Got data \"%s\"", strippedData)

	jsonData, err := sm.extractFirstPart(data, f.formDataContentType)
	if err != nil {
		return nil, nil, fmt.Errorf("submit file: %v", err)
	}
	strippedJsonData := re.ReplaceAllString(string(jsonData), "")
	sm.logger.Printf("Got JSON data \"%s\"", strippedJsonData)

	md5Hash := md5.Sum(jsonData)
	sha1Hash := sha1.Sum(jsonData)
	sha256Hash := sha256.Sum256(jsonData)
	digest := Digest{
		MD5:    fmt.Sprintf("%x", md5Hash),
		SHA1:   fmt.Sprintf("%x", sha1Hash),
		SHA256: fmt.Sprintf("%x", sha256Hash),
	}
	id := uuid.New().String()
	sm.logger.Printf("SubmitFile (%s)", id)
	response := SandboxSubmitFileResponse{
		ID:        id,
		Digest:    digest,
		Arguments: "",
	}
	headers := SandboxSubmitFileResponseHeaders{
		OperationLocation:        "",
		SubmissionReserveCount:   100,
		SubmissionRemainingCount: 100,
		SubmissionCount:          1,
		SubmissionExemptionCount: 1,
	}

	sub := &submission{}
	sm.submissions[id] = sub

	err = json.Unmarshal(jsonData, &sub.submissionStatus)
	//if err != nil {
	//	sm.logger.Printf("SubmitFile (%s): unmarshal submissionStatus failed: %v", id, err)
	//	return nil, nil, fmt.Errorf("unmarshal: %w", err)
	//}
	sm.logger.Printf("SubmitFile (%s): submissionStatus unmarshal err: %v", id, err)
	if err == nil && sub.submissionStatus.Status == StatusFailed {
		sm.logger.Printf("SubmitFile (%s): failed", id)
		sub.submissionStatus.ID = id
		sub.submissionStatus.Action = ActionAnalyzeFile
		sub.submissionStatus.CreatedDateTime = VisionOneTime(time.Now())
		sub.submissionStatus.LastActionDateTime = VisionOneTime(time.Now())
		sub.submissionStatus.IsCached = false
		sub.submissionStatus.Digest = digest
		return &response, &headers, nil
	}

	err = json.Unmarshal(jsonData, &sub.analysisResult)
	if err != nil {
		sm.logger.Printf("SubmitFile (%s): unmarshal analysisResult failed: %v", id, err)
		return nil, nil, fmt.Errorf("submit file: %v", err)
	}
	sm.logger.Printf("SubmitFile (%s): unmarshaled analysisResult", id)
	sub.analysisResult.ID = id
	sub.analysisResult.Type = "file"
	sub.analysisResult.Digest = digest

	sub.submissionStatus = SandboxSubmissionStatusResponse{
		ID:                 id,
		Action:             ActionAnalyzeFile,
		Status:             StatusRunning,
		Error:              Error{},
		CreatedDateTime:    VisionOneTime(time.Now()),
		LastActionDateTime: VisionOneTime(time.Now()),
		ResourceLocation:   "",
		IsCached:           false,
		Digest:             digest,
		Arguments:          "",
	}
	sm.logger.Printf("SubmitFile (%s): created submission with status %v", id, sub.submissionStatus.Status)
	sm.logger.Printf("SubmitFile (%s): created analysis result with risk level %v", id, sub.analysisResult.RiskLevel)
	return &response, &headers, nil
}

func (sm *SandboxMockupRAM) SubmissionStatus(f *sandboxSubmissionStatusRequest) (*SandboxSubmissionStatusResponse, error) {
	sm.logger.Printf("SubmissionStatus (%s)", f.id)
	submission, ok := sm.submissions[f.id]
	if !ok {
		return nil, fmt.Errorf("submission %s not found", f.id)
	}
	result := submission.submissionStatus
	submission.submissionStatus.Status = StatusSucceeded
	return &result, nil
}

func (sm *SandboxMockupRAM) AnalysisResults(f *sandboxAnalysisResultsRequest) (*SandboxAnalysisResultsResponseItem, error) {
	sm.logger.Printf("AnalysisResults (%s)", f.id)
	submission, ok := sm.submissions[f.id]
	sm.logger.Printf("AnalysisResults (%s): Submission found: %v", f.id, ok)
	if !ok {
		return nil, fmt.Errorf("submission %s not found", f.id)
	}
	if submission.submissionStatus.Status != StatusSucceeded {
		// Question is, how actual VOneSandbox API behaves?...
		return nil, nil
	}
	sm.logger.Printf("AnalysisResults (%s): returning risk level %v", f.id, submission.analysisResult.RiskLevel)
	return &submission.analysisResult, nil
}

func (sm *SandboxMockupRAM) ListSubmissions(f *sandboxSubmissionsRequest) (*SandboxSubmissionsResponse, error) {
	sm.logger.Printf("ListSubmissions")
	response := SandboxSubmissionsResponse{
		Items:    nil,
		NextLink: "",
	}
	for id, s := range sm.submissions {
		if strings.Contains(f.parameters["filter"], id) {
			response.Items = append(response.Items, s.submissionStatus)
			sm.logger.Printf("ListSubmissions: id=%s, status=%v", id, s.submissionStatus.Status)
			s.submissionStatus.Status = StatusSucceeded
		}
	}
	sm.logger.Printf("ListSubmissions: %d items", len(response.Items))
	return &response, nil
}

func (sm *SandboxMockupRAM) ListAnalysisResults(f *sandboxListAnalysisResultsRequest) (*SandboxListAnalysisResultResponse, error) {
	sm.logger.Printf("ListAnalysisResults")
	results := SandboxListAnalysisResultResponse{
		Items:    nil,
		NextLink: "",
	}
	for id, s := range sm.submissions {
		if !strings.Contains(f.parameters["filter"], id) {
			continue
		}
		if s.submissionStatus.Status != StatusSucceeded {
			continue
		}
		results.Items = append(results.Items, s.analysisResult)
		sm.logger.Printf("ListAnalysisResults: id=%s, riskLevel=%v", id, s.analysisResult.RiskLevel)
	}
	sm.logger.Printf("ListAnalysisResults: %d items", len(results.Items))
	return &results, nil
}
