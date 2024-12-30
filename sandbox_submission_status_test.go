package vone

import (
	"encoding/json"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	jsonData := []byte(`
	{
		"id": "bd670550-7e57-4d10-901b-f6134a51f3bf",
		"action": "analyzeFile",
		"status": "running",
		"createdDateTime": "2024-02-04T20:07:50Z",
		"lastActionDateTime": "2024-02-04T20:07:50Z",
		"resourceLocation": "https://api.eu.xdr.trendmicro.com/v3.0/sandbox/analysisResults/bd670550-7e57-4d10-901b-f6134a51f3bf",
		"isCached": true,
		"digest": {
		  "md5": "4c7ff71bee3a293c789635837d26522d",
		  "sha1": "51ca64a676f1d2c41b4541fda88458e486a09dbb",
		  "sha256": "5265798ad1797d01c67927d110f080bfa1ecc3e58f56900f34c2072017b2ed69"
		}
	  }
	`)
	var response SandboxSubmissionStatusResponse
	if err := json.Unmarshal(jsonData, &response); err != nil {
		t.Fatal(err)
	}
	{
		expected := StatusRunning
		actual := response.Status
		if actual != expected {
			t.Errorf("Expected %v, but got %v", expected, actual)

		}
	}
	{
		expected := "4c7ff71bee3a293c789635837d26522d"
		actual := response.Digest.MD5
		if actual != expected {
			t.Errorf("Expected %v, but got %v", expected, actual)

		}
	}
}
