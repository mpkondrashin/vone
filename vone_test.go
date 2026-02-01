package vone

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestVisionOneTime_String(t *testing.T) {
	var v VisionOneTime

	data := []byte(`"2024-12-01T10:00:00Z"`)

	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if time.Time(v).IsZero() {
		t.Fatalf("time should not be zero")
	}
}

func TestUnmarshalSandboxAnalysisResultsResponse(t *testing.T) {
	jsonData := `{
  "riskLevel": "noRisk",
  "trueFileType": "exe"
}`

	var response SandboxAnalysisResultsResponseItem
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.RiskLevel != RiskLevelNoRisk {
		t.Fatalf("expected riskLevel noRisk, got %s", response.RiskLevel)
	}

	if response.TrueFileType != "exe" {
		t.Fatalf("expected trueFileType exe, got %s", response.TrueFileType)
	}
}

func TestSandboxSubmissionStatusResponse_NullDates(t *testing.T) {
	data := []byte(`
{
  "id": "123",
  "action": "analyzeFile",
  "status": "running",
  "error": {},
  "createdDateTime": null,
  "lastActionDateTime": null,
  "resourceLocation": "/tmp/file.exe",
  "isCached": false,
  "digest": {},
  "arguments": ""
}
`)

	var resp SandboxSubmissionStatusResponse

	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseError(t *testing.T) {
	data := []byte(`{"error":{"code":"AccessDenied","message":"No permission to access this resource"}}`)

	vOneErr, err := ErrorFromReader(bytes.NewReader(data))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vOneErr.Code != ErrorCodeAccessDenied {
		t.Errorf("expected code AccessDenied, got %s", vOneErr.Code)
	}

	if vOneErr.Message != "No permission to access this resource" {
		t.Errorf("expected message No permission to access this resource, got %s", vOneErr.Message)
	}
}
