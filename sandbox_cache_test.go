package vone

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func prepare(t *testing.T) {

}

func TestAddTwice(t *testing.T) {
	folder := "testing/cache"
	err := os.RemoveAll(folder)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(folder, 0775)
	if err != nil {
		t.Fatal(err)
	}

	data1 := &SandboxAnalysisResultsResponse{
		ID:   "abc",
		Type: "file",
		Digest: Digest{
			MD5:    "abcd",
			SHA1:   "efgh",
			SHA256: "jklm",
		},
		Arguments:                  "",
		AnalysisCompletionDateTime: VisionOneTime(time.Now()),
		RiskLevel:                  RiskLevelLow,
		DetectionNames:             []string{"Virus"},
		ThreatTypes:                []string{"Malware"},
		TrueFileType:               "PE-EXE",
	}
	dbPath := filepath.Join(folder, "cache.sqlite3")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(fmt.Errorf("%s: %w", dbPath, err))
	}
	cache, err := NewCache(db, dbPath)
	err = cache.Add(context.Background(), data1)
	if err != nil {
		t.Fatal(err)
	}
	data2 := &SandboxAnalysisResultsResponse{
		ID:   "abc",
		Type: "file",
		Digest: Digest{
			MD5:    "abcd",
			SHA1:   "efgh2",
			SHA256: "jklm3",
		},
		Arguments:                  "",
		AnalysisCompletionDateTime: VisionOneTime(time.Now()),
		RiskLevel:                  RiskLevelLow,
		DetectionNames:             []string{"Virus"},
		ThreatTypes:                []string{"Malware"},
		TrueFileType:               "PE-EXE",
	}
	err = cache.Add(context.Background(), data2)
	if err != nil {
		t.Fatal(err)
	}
	count, err := cache.Count(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expected := 1
	actual := count
	if expected != actual {
		t.Errorf("Expected count %d, but got %d", expected, actual)
	}
}

func TestAddTwo(t *testing.T) {
	folder := "testing/cache"
	err := os.RemoveAll(folder)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(folder, 0775)
	if err != nil {
		t.Fatal(err)
	}

	data1 := &SandboxAnalysisResultsResponse{
		ID:   "abc",
		Type: "file",
		Digest: Digest{
			MD5:    "abcd",
			SHA1:   "efgh",
			SHA256: "jklm",
		},
		Arguments:                  "",
		AnalysisCompletionDateTime: VisionOneTime(time.Now()),
		RiskLevel:                  RiskLevelLow,
		DetectionNames:             []string{"Virus"},
		ThreatTypes:                []string{"Malware"},
		TrueFileType:               "PE-EXE",
	}
	dbPath := filepath.Join(folder, "cache.sqlite3")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(fmt.Errorf("%s: %w", dbPath, err))
	}
	cache, err := NewCache(db, dbPath)
	err = cache.Add(context.Background(), data1)
	if err != nil {
		t.Fatal(err)
	}
	data2 := &SandboxAnalysisResultsResponse{
		ID:   "abc",
		Type: "file",
		Digest: Digest{
			MD5:    "abcd1",
			SHA1:   "efgh2",
			SHA256: "jklm3",
		},
		Arguments:                  "",
		AnalysisCompletionDateTime: VisionOneTime(time.Now()),
		RiskLevel:                  RiskLevelLow,
		DetectionNames:             []string{"Virus"},
		ThreatTypes:                []string{"Malware"},
		TrueFileType:               "PE-EXE",
	}
	err = cache.Add(context.Background(), data2)
	if err != nil {
		t.Fatal(err)
	}
	count, err := cache.Count(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expected := 2
	actual := count
	if expected != actual {
		t.Errorf("Expected count %d, but got %d", expected, actual)
	}
}

func TestAddAndQuery(t *testing.T) {
	folder := "testing/cache"
	err := os.RemoveAll(folder)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(folder, 0775)
	if err != nil {
		t.Fatal(err)
	}
	sha1 := "efgh"
	data1 := &SandboxAnalysisResultsResponse{
		ID:   "abc",
		Type: "file",
		Digest: Digest{
			MD5:    "abcd",
			SHA1:   sha1,
			SHA256: "jklm",
		},
		Arguments:                  "",
		AnalysisCompletionDateTime: VisionOneTime(time.Now()),
		RiskLevel:                  RiskLevelLow,
		DetectionNames:             []string{"Virus"},
		ThreatTypes:                []string{"Malware"},
		TrueFileType:               "PE-EXE",
	}
	dbPath := filepath.Join(folder, "cache.sqlite3")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(fmt.Errorf("%s: %w", dbPath, err))
	}
	cache, err := NewCache(db, dbPath)
	err = cache.Add(context.Background(), data1)
	if err != nil {
		t.Fatal(err)
	}
	data2, updated, err := cache.Query(context.Background(), sha1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Updated: %v", updated)
	expected := "PE-EXE"
	actual := data2.TrueFileType
	if expected != actual {
		t.Errorf("Expected TrueFileType %s, but got %s", expected, actual)
	}
}

func TestAddAndIterate(t *testing.T) {
	folder := "testing/cache"
	err := os.RemoveAll(folder)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(folder, 0775)
	if err != nil {
		t.Fatal(err)
	}
	sha1 := "efgh"
	data1 := &SandboxAnalysisResultsResponse{
		ID:   "abc",
		Type: "file",
		Digest: Digest{
			MD5:    "abcd",
			SHA1:   sha1,
			SHA256: "jklm",
		},
		Arguments:                  "",
		AnalysisCompletionDateTime: VisionOneTime(time.Now()),
		RiskLevel:                  RiskLevelLow,
		DetectionNames:             []string{"Virus"},
		ThreatTypes:                []string{"Malware"},
		TrueFileType:               "PE-EXE",
	}
	dbPath := filepath.Join(folder, "cache.sqlite3")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(fmt.Errorf("%s: %w", dbPath, err))
	}
	cache, err := NewCache(db, dbPath)
	err = cache.Add(context.Background(), data1)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	err = cache.IterateCache(context.Background(), func(data *SandboxAnalysisResultsResponse, updated time.Time) error {
		count += 1
		if !strings.EqualFold(data.Digest.SHA1, sha1) {
			t.Errorf("Expected count %s, but got %s", sha1, data.Digest.SHA1)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := 1
	actual := count
	if expected != actual {
		t.Errorf("Expected count %d, but got %d", expected, actual)
	}
}
