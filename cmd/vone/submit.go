package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mpkondrashin/vone"
	"github.com/spf13/viper"
)

type commandSubmit struct {
	baseCommand
}

func newCommandSubmit() *commandSubmit {
	c := &commandSubmit{}
	c.Setup(cmdSubmit)
	return c
}

func (c *commandSubmit) Execute() error {
	var wg sync.WaitGroup

	filePath := viper.GetString(flagFileName)
	if filePath != "" {
		go c.SubmitFileGoRoutine(filePath, &wg)
	}

	fileMask := viper.GetString(flagMask)
	if fileMask != "" {
		matches, err := filepath.Glob(fileMask)
		if err != nil {
			log.Printf("%s: %v", fileMask, err)
		}
		for _, m := range matches {
			wg.Add(1)
			go c.SubmitFileGoRoutine(m, &wg)
		}
		// np tears stats
		// 328
	}

	url := viper.GetString(flagURL)
	if url != "" {
		//log.Println("XXX url", url)
		wg.Add(1)
		go c.SubmitURLGoRoutine(url, &wg)
	}
	c.ProcessURLsFile(&wg)
	wg.Wait()
	return nil
}

func (c *commandSubmit) ProcessURLsFile(wg *sync.WaitGroup) {
	urlfile := viper.GetString(flagURLsFile)
	if urlfile == "" {
		return
	}
	data, err := os.ReadFile(urlfile)
	if err != nil {
		log.Println(err)
		return
	}
	urls := strings.Split(string(data), "\n")
	for _, url := range urls {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}
		wg.Add(1)
		c.SubmitURLGoRoutine(url, wg)
	}
}

func (c *commandSubmit) SubmitFileGoRoutine(filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := c.SubmitFile(filePath)
	if err != nil {
		log.Println(err)
	}
}

func (c *commandSubmit) SubmitURLGoRoutine(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := c.SubmitURL(url)
	if err != nil {
		log.Println(err)
	}
}

func (c *commandSubmit) SubmitFile(filePath string) error {
	log.Printf("Uploading %s", filePath)
	f, err := c.visionOne.SandboxSubmitFile().SetFilePath(filePath)
	if err != nil {
		return err
	}
	response, headers, err := f.Do()
	if err != nil {
		return err
	}
	id := response.ID
	log.Printf("File accepted. ID = %s", id)
	c.LogQuota(id, headers)
	return c.ProcessObject(id)
}

func (c *commandSubmit) SubmitURL(url string) error {
	log.Printf("Uploading URL %s", url)
	submit := c.visionOne.SandboxSubmitURLs().AddURL(url)
	response, headers, err := submit.Do()
	if err != nil {
		return err
	}
	if len(response) != 1 {
		return fmt.Errorf("wrong response length: %d", len(response))
	}
	resp := response[0]
	if vone.GetHTTPCodeRange(resp.Status) != vone.HTTPCodeSuccess {
		return fmt.Errorf("%s: %s", resp.Body.Error.Code, resp.Body.Error.Code)
	}
	id := response[0].Body.ID
	log.Printf("URL accepted. ID = %s", id)
	c.LogQuota(id, headers)
	return c.ProcessObject(id)
}

func (c *commandSubmit) LogQuota(id string, headers *vone.SandboxSubmitFileResponseHeaders) {
	log.Printf("%s Daily quota: %d", id, headers.SubmissionReserveCount)
	note := ""
	if headers.SubmissionRemainingCount == 0 {
		note = " (consider to add Credits to sandbox feature)"
	}
	log.Printf("%s Quota left: %d%s", id, headers.SubmissionRemainingCount, note)
	log.Printf("%s Today submissions: %d", id, headers.SubmissionCount)
	log.Printf("%s Today submissions of unsupported files: %d (not accounted in quota)", id, headers.SubmissionExemptionCount)
}

func (c *commandSubmit) ProcessObject(id string) error {
	if err := c.WaitForResult(id); err != nil {
		return err
	}
	malicious, err := c.AnalysisResult(id)
	if err != nil {
		return err
	}
	if !malicious {
		return nil
	}
	if err := c.SuspiciousObjects(id); err != nil {
		return err
	}
	if err := c.PDFReport(id); err != nil {
		return err
	}
	if err := c.InvestigationPackage(id); err != nil {
		return err
	}
	return nil
}

var ErrTimeout = errors.New("Timeout")

func (c *commandSubmit) WaitForResult(id string) error {
	timeout := viper.GetDuration(flagTimeout)
	endTime := time.Now().Add(timeout)
	for {
		status, err := c.visionOne.SandboxSubmissionStatus(id).Do()
		if err != nil {
			return fmt.Errorf("WaitForResult(%s): %w", id, err)
		}
		log.Printf("%s Status: %v", id, status.Status)
		switch status.Status {
		case vone.StatusSucceeded:
			return nil
		case vone.StatusRunning:
			if time.Now().After(endTime) {
				return ErrTimeout
			}
			time.Sleep(5 * time.Second)
		case vone.StatusFailed:
			return fmt.Errorf("%s: %s", status.Error.Code, status.Error.Message)
		default:
			return fmt.Errorf("unknown status: %s", status.Status)
		}
	}
}

func (c *commandSubmit) AnalysisResult(id string) (bool, error) {
	results, err := c.visionOne.SandboxAnalysisResults(id).Do()
	if err != nil {
		return false, err
	}
	log.Printf("%s Type: %s", id, results.Type)
	log.Printf("%s TrueFileType: %s", id, results.TrueFileType)
	log.Printf("%s RiskLevel: %s", id, results.RiskLevel)
	if len(results.DetectionNames) > 0 {
		log.Printf("%s DetectionNames: %s", id, strings.Join(results.DetectionNames, ", "))
	}
	if len(results.ThreatTypes) > 0 {
		log.Printf("%s ThreatTypes: %s", id, strings.Join(results.ThreatTypes, ", "))
	}
	return results.RiskLevel != vone.RiskLevelNoRisk, nil //"high", nil
}

func ListSHA1(so *vone.SandboxSuspiciousObjectsResponse) (result []string) {
	m := make(map[string]struct{})
	for _, item := range so.Items {
		if item.RootSHA1 == "" {
			continue
		}
		m[item.RootSHA1] = struct{}{}
	}
	for sha1 := range m {
		result = append(result, sha1)
	}
	return
}

func ListIP(so *vone.SandboxSuspiciousObjectsResponse) (result []string) {
	m := make(map[string]struct{})
	for _, item := range so.Items {
		if item.IP == "" {
			continue
		}
		m[item.IP] = struct{}{}
	}
	for sha1 := range m {
		result = append(result, sha1)
	}
	return
}

func (c *commandSubmit) SuspiciousObjects(id string) error {
	log.Printf("%s Request Suspicious Objects", id)
	so, err := c.visionOne.SandboxSuspiciousObjects(id).Do()
	if err != nil {
		return err
	}
	soSHA1 := ListSHA1(so)
	for _, sha1 := range soSHA1 {
		log.Printf("%s Suspicious Object SHA1: %s", id, sha1)
	}
	soIP := ListIP(so)
	for _, ip := range soIP {
		log.Printf("%s Suspicious Object IP: %s", id, ip)
	}
	return nil
}

func (c *commandSubmit) PDFReport(id string) error {
	log.Printf("%s Download report PDF", id)
	pdfFileName := id + ".pdf"
	if err := c.visionOne.SandboxDownloadResults(id).Store(pdfFileName); err != nil {
		return err
	}
	log.Printf("%s PDF report saved: %s", id, pdfFileName)
	return nil
}

func (c *commandSubmit) InvestigationPackage(id string) error {
	log.Printf("%s Download Investigation Package", id)
	zipFileName := id + ".zip"
	if err := c.visionOne.SandboxInvestigationPackage(id).Store(zipFileName); err != nil {
		return err
	}
	log.Printf("%s Investigation Package Saved: %s", id, zipFileName)
	return nil
}

func (c *commandSubmit) Setup(name string) {
	c.baseCommand.Setup(name)
	c.fs.String(flagFileName, "", "Sample file path")
	c.fs.String(flagMask, "", "Sample files mask")
	c.fs.String(flagURL, "", "Sample URL")
	c.fs.String(flagURLsFile, "", "File with URLs")
	c.fs.Duration(flagTimeout, 10*time.Minute, "Analisys timeout")
}
