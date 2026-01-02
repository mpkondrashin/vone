package vone

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SubmissionState int

const (
	state SubmissionState = iota
)

type (
	submission struct {
		state            SubmissionState
		analysisResult   SandboxAnalysisResultsResponse
		submissionStatus SandboxSubmissionStatusResponse
	}

	SandboxMockup struct {
		submissions map[string]*submission
		logger      *log.Logger
	}
)

func NewSandboxMockup() *SandboxMockup {
	return &SandboxMockup{
		submissions: make(map[string]*submission),
		logger:      log.New(io.Discard, "", 0),
	}
}

func (s *SandboxMockup) EnableFileLogging(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	s.logger = log.New(
		f,
		"",
		log.LstdFlags,
	)
	return nil
}

func (sm *SandboxMockup) SubmitFile(f *SandboxSubmitFileToSandboxFunc) (*SandboxSubmitFileResponse, *SandboxSubmitFileResponseHeaders, error) {
	sm.logger.Println("SubmitFile")
	data, err := io.ReadAll(f.Request)
	if err != nil {
		return nil, nil, err
	}
	md5Hash := md5.Sum(data)
	sha1Hash := sha1.Sum(data)
	sha256Hash := sha256.Sum256(data)
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
	sm.logger.Printf("SubmitFile (%s): response created", id)
	headers := SandboxSubmitFileResponseHeaders{
		OperationLocation:        "",
		SubmissionReserveCount:   100,
		SubmissionRemainingCount: 100,
		SubmissionCount:          1,
		SubmissionExemptionCount: 1,
	}
	sm.logger.Printf("SubmitFile (%s): headers created", id)

	sub := &submission{}
	sm.submissions[id] = sub
	sm.logger.Printf("SubmitFile (%s): submission tracked", id)

	err = json.Unmarshal(data, &sub.submissionStatus)
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

	err = json.Unmarshal(data, &sub.analysisResult)
	if err != nil {
		sm.logger.Printf("SubmitFile (%s): unmarshal analysisResult failed: %v", id, err)
		return nil, nil, fmt.Errorf("unmarshal: %w", err)
	}
	sm.logger.Printf("SubmitFile (%s): success", id)
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

func (sm *SandboxMockup) SubmissionStatus(f *SandboxSubmissionStatusFunc) (*SandboxSubmissionStatusResponse, error) {
	sm.logger.Printf("SubmissionStatus (%s)", f.id)
	submission, ok := sm.submissions[f.id]
	if !ok {
		return nil, fmt.Errorf("submission %s not found", f.id)
	}
	result := submission.submissionStatus
	submission.submissionStatus.Status = StatusSucceeded
	return &result, nil
}

func (sm *SandboxMockup) AnalysisResults(f *SandboxAnalysisResultsFunc) (*SandboxAnalysisResultsResponse, error) {
	sm.logger.Printf("AnalysisResults (%s)", f.id)
	submission, ok := sm.submissions[f.id]
	sm.logger.Printf("AnalysisResults (%s): Submission found: %v", f.id, ok)
	if !ok {
		return nil, fmt.Errorf("submission %s not found", f.id)
	}
	if submission.submissionStatus.Status != StatusSucceeded {
		return nil, fmt.Errorf("submission %s status is not succeed", f.id)
	}
	sm.logger.Printf("AnalysisResults (%s): returning risk level %v", f.id, submission.analysisResult.RiskLevel)
	return &submission.analysisResult, nil
}

func (sm *SandboxMockup) ListSubmissions(f *SandboxSubmissionsFunc) (*SandboxSubmissionsResponse, error) {
	sm.logger.Printf("ListSubmissions")
	response := SandboxSubmissionsResponse{
		Items:    nil,
		NextLink: "",
	}
	for id, s := range sm.submissions {
		if strings.Contains(f.parameters["filter"], id) {
			response.Items = append(response.Items, s.submissionStatus)
			s.submissionStatus.Status += 1
		}
	}
	sm.logger.Printf("ListSubmissions: %d items", len(response.Items))
	return &response, nil
}

func (sm *SandboxMockup) ListAnalysisResults(f *SandboxListAnalysisResultsFunc) (*SandboxListAnalysisResultResponse, error) {
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
	}
	sm.logger.Printf("ListAnalysisResults: %d items", len(results.Items))
	return &results, nil
}
