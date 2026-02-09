/*
	Trend Micro Vision One API SDK
	(c) 2025 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	main.go - various SDK tests
*/

package vone

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func GetVOne(t *testing.T) *VOne {
	url := os.Getenv("TEST_VONE_URL")
	if url == "" {
		t.Skip("No url, so skip")
	}
	token := os.Getenv("TEST_VONE_TOKEN")
	if token == "" {
		t.Skip("No token, so skip")
	}
	return NewVOne(url, token)
}

func TestSandbox_DailyReserve(t *testing.T) {
	v1 := GetVOne(t)
	t.Log("*** SandboxGetDailyReserve ***")
	reserve, err := v1.SandboxDailyReserve().Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Submission Remaining Count: %d", reserve.SubmissionRemainingCount)
	t.Logf("Result: %v", reserve)
}

func TestSandbox_SubmitURLs(t *testing.T) {
	v1 := GetVOne(t)
	t.Log("*** Submit URLs To Sandbox ***")
	urls := []string{"test0001", "test0002"}
	resp, _, err := v1.SandboxSubmitURLs().AddURLs(urls).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	for _, each := range resp {
		if each.Status == 207 {
			t.Logf("%v: %v", each.Body.URL, each.Body.ID)
		} else {
			t.Logf("%v: %v", each.Body.URL, each.Body.Error.Code)
		}
	}
}

func TestSandbox_SubmitFile(t *testing.T) {
	v1 := GetVOne(t)
	t.Log("*** Submit File To Sandbox ***")
	filePath := "sample.exe"
	submit := v1.SandboxSubmitFile()
	err := submit.SetFilePath(context.Background(), filePath)
	if err != nil {
		t.Fatal(err)
	}
	resp, _, err := submit.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Accepted: %v", filePath)
	t.Logf("ID: %s", resp.ID)
	t.Logf("MD5: %s", resp.Digest.MD5)
	t.Logf("SHA1: %s", resp.Digest.SHA1)
	t.Logf("SHA256: %s", resp.Digest.SHA256)

	virus := strings.NewReader(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`)
	fileName := "eicar.com"
	submit = v1.SandboxSubmitFile()
	if err := submit.SetReader(context.Background(), virus, fileName); err != nil {
		t.Fatal(err)
	}
	resp, _, err = submit.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Accepted: %v", fileName)
	t.Logf("ID: %s", resp.ID)
	t.Logf("MD5: %s", resp.Digest.MD5)
	t.Logf("SHA1: %s", resp.Digest.SHA1)
	t.Logf("SHA256: %s", resp.Digest.SHA256)
}

func TestSandbox_SubmitFile_Malware(t *testing.T) {
	v1 := GetVOne(t)
	t.Logf("*** Submit File To Sandbox ***")
	filePath := "sample.exe"
	submit := v1.SandboxSubmitFile()
	err := submit.SetFilePath(context.Background(), filePath)
	if err != nil {
		t.Fatal(err)
	}
	resp, _, err := submit.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Accepted: %v", filePath)
	t.Logf("ID: %s", resp.ID)
	t.Logf("MD5: %s", resp.Digest.MD5)
	t.Logf("SHA1: %s", resp.Digest.SHA1)
	t.Logf("SHA256: %s", resp.Digest.SHA256)

	virus := strings.NewReader(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`)
	fileName := "eicar.com"
	submit = v1.SandboxSubmitFile()
	if err := submit.SetReader(context.Background(), virus, fileName); err != nil {
		panic(err)
	}
	resp, _, err = submit.Do(context.Background())
	if err != nil {
		panic(err)
	}
	t.Logf("Accepted: %v", fileName)
	t.Logf("ID: %s", resp.ID)
	t.Logf("MD5: %s", resp.Digest.MD5)
	t.Logf("SHA1: %s", resp.Digest.SHA1)
	t.Logf("SHA256: %s", resp.Digest.SHA256)
}

func TestSandbox_SubmitFile_FromFolder(t *testing.T) {
	v1 := GetVOne(t)
	t.Log("*** Submit File To Sandbox ***")
	dir, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, each := range dir {
		if each.IsDir() {
			continue
		}
		if !strings.HasSuffix(each.Name(), ".exe") {
			continue
		}
		filePath := each.Name()
		submit := v1.SandboxSubmitFile()
		err := submit.SetFilePath(context.Background(), filePath)
		if err != nil {
			panic(err)
		}
		resp, _, err := submit.Do(context.Background())
		if err != nil {
			panic(err)
		}
		t.Logf("Accepted: %v", filePath)
		t.Logf("ID: %s", resp.ID)
		t.Logf("MD5: %s", resp.Digest.MD5)
		t.Logf("SHA1: %s", resp.Digest.SHA1)
		t.Logf("SHA256: %s", resp.Digest.SHA256)
	}
}
func TestSandbox_SubmitFile_Unsupported(t *testing.T) {
	v1 := GetVOne(t)
	filePath := "main.go"
	submit := v1.SandboxSubmitFile()
	err := submit.SetFilePath(context.Background(), filePath)
	if err != nil {
		t.Fatal(err)
	}
	resp, _, err := submit.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Accepted: %v", filePath)
	t.Logf("ID: %s", resp.ID)
	t.Logf("MD5: %s", resp.Digest.MD5)
	t.Logf("SHA1: %s", resp.Digest.SHA1)
	t.Logf("SHA256: %s", resp.Digest.SHA256)
}
func TestSandbox_SubmitFile_ListSubmissions(t *testing.T) {
	v1 := GetVOne(t)
	t.Log("*** List Submissions ***")
	listSubmissions := v1.SandboxListSubmissions().OrderBy(CreatedDateTime, Asc)
	resp, err := listSubmissions.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range resp.Items {
		t.Logf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
	}
	t.Logf("Next: %s", resp.NextLink)
	if resp.NextLink != "" {
		t.Log("*** List Submissions Next ***")
		resp2, err := listSubmissions.Next(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		for _, item := range resp2.Items {
			t.Logf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
		}
	}
}

func TestSandbox_SubmitFile_RangeSubmissions(t *testing.T) {
	v1 := GetVOne(t)
	t.Log("*** Iterate List Submissions ***")
	ls := v1.SandboxListSubmissions().OrderBy(CreatedDateTime, Asc)
	for item, err := range ls.Paginator().Range(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
		status := v1.SandboxSubmissionStatus(item.ID)
		result, err := status.Do(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("ID: %v, ResourceLocation: %s", item.ID, result.ResourceLocation)
	}
}

func TestSandbox_SubmitFile_ListSubmissions_AnalysisResults(t *testing.T) {
	v1 := GetVOne(t)
	t.Log("*** Iterate List Submissions & Get Analysis Results ***")
	ls := v1.SandboxListSubmissions().OrderBy(CreatedDateTime, Asc)
	for item, err := range ls.Paginator().Range(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		results := v1.SandboxAnalysisResults(item.ID)
		result, err := results.Do(context.Background())
		if err != nil {
			var perr *Error
			if !errors.As(err, &perr) {
				t.Fatal(err)
			}
			if perr.Code != ErrorCodeNotFound {
				t.Fatal(err)
			}
			t.Logf("ID: %v, NotFound", item.ID)
			break
		}
		t.Logf("ID: %v, RiskLevel: %s", item.ID, result.RiskLevel)
	}
}

func TestSandbox_SubmitFile_ListSubmissions_DownloadResults(t *testing.T) {
	v1 := GetVOne(t)
	t.Logf("*** Iterate List Submissions & Download Result ***")
	ls := v1.SandboxListSubmissions()
	for item, err := range ls.Paginator().Range(context.Background()) {
		if err != nil {
			t.Log(err)
		}
		result, err := v1.SandboxAnalysisResults(item.ID).Do(context.Background())
		if err != nil {
			var perr *Error
			if !errors.As(err, &perr) {
				t.Fatal(err)
			}
			if perr.Code != ErrorCodeNotFound {
				t.Fatal(err)
			}
			t.Logf("ID: %v, NotFound", item.ID)
			break
		}
		t.Logf("ID: %v, RiskLevel: %s", item.ID, result.RiskLevel)
		err = v1.SandboxDownloadResults(item.ID).Store(context.Background(), item.ID+".pdf")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSandbox_SubmitFile_ListSubmissions_SuspiciousObjects(t *testing.T) {
	v1 := GetVOne(t)
	t.Log("*** Iterate List Submissions & Suspicious Objects ***")
	ls := v1.SandboxListSubmissions()
	for item, err := range ls.Paginator().Range(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		result, err := v1.SandboxSuspiciousObjects(item.ID).Do(context.Background())
		if err != nil {
			var perr *Error
			if !errors.As(err, &perr) {
				t.Fatal(err)
			}
			if perr.Code != ErrorCodeNotFound {
				t.Fatal(err)
			}
			t.Logf("ID: %v, NotFound", item.ID)
			break
		}
		for _, each := range result.Items {
			t.Logf("ID: %v, RiskLevel: %s, SHA1: %s, IP: %s", item.ID, each.RiskLevel, each.RootSHA1, each.IP)
		}
	}
}

func TestSandbox_SubmitFile_IterateResults(t *testing.T) {
	v1 := GetVOne(t)
	t.Log("*** Submit Files And Iterate Analysis Results ***")
	filePath := "main.go"
	submit := v1.SandboxSubmitFile()
	err := submit.SetFilePath(context.Background(), filePath)
	if err != nil {
		t.Fatal(err)
	}
	resp1, _, err := submit.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Accepted: %v", filePath)
	log.Printf("ID: %s", resp1.ID)
	log.Printf("MD5: %s", resp1.Digest.MD5)
	log.Printf("SHA1: %s", resp1.Digest.SHA1)
	log.Printf("SHA256: %s", resp1.Digest.SHA256)

	virus := strings.NewReader(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`)
	fileName := "eicar.com"
	submit = v1.SandboxSubmitFile()
	if err := submit.SetReader(context.Background(), virus, fileName); err != nil {
		t.Fatal(err)
	}
	resp2, _, err := submit.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Accepted: %v", fileName)
	log.Printf("ID: %s", resp2.ID)
	log.Printf("MD5: %s", resp2.Digest.MD5)
	log.Printf("SHA1: %s", resp2.Digest.SHA1)
	log.Printf("SHA256: %s", resp2.Digest.SHA256)
	filter := fmt.Sprintf("(id eq '%s') or (id eq '%s')", resp1.ID, resp2.ID)
	listSubmissions := func() {
		ls := v1.SandboxListSubmissions()
		ls.Filter(filter)
		for item, err := range ls.Paginator().Range(context.Background()) {
			if err != nil {
				t.Fatal(err)
			}
			t.Logf(`Submission ID %v
Action: %v
Status: %v
Error Code: %v
Error Message: %v
IsCached: %v
SHA1: %v`,
				item.ID,
				item.Action,
				item.Status,
				item.Error.Code,
				item.Error.Message,
				item.IsCached,
				item.Digest.SHA1,
			)

		}
	}
	listSubmissions()
	log.Println("Sleep 10 seconds")
	time.Sleep(10 * time.Second)
	listSubmissions()
	listResults := v1.SandboxListAnalysisResults()
	listResults.StartDateTime(time.Now().Add(-time.Hour * 10000))
	listResults.EndDateTime(time.Now().Add(time.Hour * 10000))
	listResults.Filter(filter)
	for report, err := range listResults.Paginator().Range(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		log.Printf(`Report for ID: %v
DetectionNames: %v
SHA1: %v
RiskLevel: %v
ThreatTypes: %v
TrueFileType: %v`,
			report.ID,
			report.DetectionNames,
			report.Digest.SHA1,
			report.RiskLevel,
			report.ThreatTypes,
			report.TrueFileType)
	}
}
