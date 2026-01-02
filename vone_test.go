package vone

import (
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

	var response SandboxAnalysisResultsResponse
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
