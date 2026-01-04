package vone

import (
	"bytes"
	"context"
	"os"
	"testing"
)

func TestSandboxMockup(t *testing.T) {
	v1 := NewVOne("https://api.trendmicro.com", "token")
	mockup := NewSandboxMockup()
	mockupFilePath := "mockup.log"
	if err := os.RemoveAll(mockupFilePath); err != nil {
		t.Fatalf("RemoveAll Error: %v", err)
	}
	if err := mockup.EnableFileLogging(mockupFilePath); err != nil {
		t.Fatalf("EnableFileLogging Error: %v", err)
	}
	v1.SetMockup(mockup)
	t.Logf("Test Submit File with failed status")
	submitFile := v1.SandboxSubmitFile()
	submitFile.Request = bytes.NewReader([]byte(`{
  "status": "failed",
  "error": {
    "code": "InternalServerError",
    "message": "error_message"
  }
}`))
	response, headers, err := submitFile.Do(context.Background())
	if err != nil {
		t.Fatalf("Submit File Error: %v", err)
	}
	id := response.ID
	t.Logf("response: ID=%v", id)
	t.Logf("headers: SubmissionCount=%v, SubmissionRemainingCount=%v", headers.SubmissionCount, headers.SubmissionRemainingCount)
	t.Logf("Get Submission Status")
	submisionStatus := v1.SandboxSubmissionStatus(id)
	status, err := submisionStatus.Do(context.Background())
	if err != nil {
		t.Fatalf("Submission Status Error: %v", err)
	}
	t.Logf("status: %v", status.Status)
	t.Logf("Test Submit Malicious File")
	submitFile = v1.SandboxSubmitFile()
	submitFile.Request = bytes.NewReader([]byte(`{
  "riskLevel": "high",
  "detectionNames": [
    "VAN_DROPPER.UMXX"
  ],
  "threatTypes": [
    "Dropper"
  ],
  "trueFileType": "exe"
}`))
	response, headers, err = submitFile.Do(context.Background())
	if err != nil {
		t.Fatalf("Submit Malicious File Error: %v", err)
	}
	id = response.ID
	t.Logf("response: ID assigned %v", id)
	t.Logf("headers: SubmissionCount=%v, SubmissionRemainingCount=%v", headers.SubmissionCount, headers.SubmissionRemainingCount)
	t.Logf("Submission: %v", v1.mockup.submissions[id].submissionStatus.Status)
	t.Logf("Get Analysis Result For Not ready result")
	analysisResult := v1.SandboxAnalysisResults(id)
	result, err := analysisResult.Do(context.Background())
	if err != nil {
		t.Fatalf("Analysis Results Error: %v", err)
	}
	if result != nil {
		t.Fatalf("No ready result is returned!: %v", result)
	}
	t.Logf("Get First Submission Status for Malicious File")
	submisionStatus = v1.SandboxSubmissionStatus(id)
	status, err = submisionStatus.Do(context.Background())
	if err != nil {
		t.Fatalf("Submission Malicious FileStatus Error: %v", err)
	}
	t.Logf("status: %v", status.Status)
	t.Logf("Get Second Submission Status for Malicious File")
	submisionStatus = v1.SandboxSubmissionStatus(id)
	status, err = submisionStatus.Do(context.Background())
	if err != nil {
		t.Fatalf("Submission Malicious FileStatus Error: %v", err)
	}
	t.Logf("status: %v", status.Status)
	if err != nil {
		t.Fatalf("Submission Malicious FileStatus Error: %v", err)
	}
	t.Logf("Get Analysis Result For Malicious File")
	analysisResult = v1.SandboxAnalysisResults(id)
	result, err = analysisResult.Do(context.Background())
	if err != nil {
		t.Fatalf("Analysis Results Error: %v", err)
	}
	t.Logf("result: RiskLevel=%v, DetectionNames=%v, ThreatTypes=%v", result.RiskLevel, result.DetectionNames, result.ThreatTypes)
	t.Logf("Test Complete")
}
