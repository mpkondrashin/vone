package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
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
	filePath := viper.GetString(flagFileName)
	log.Printf("Uploading %s", filePath)
	f, err := c.visionOne.SandboxSubmitFile().SetFileName(filePath)
	if err != nil {
		return err
	}
	response, err := f.Do()
	if err != nil {
		return nil
	}
	id := response.ID
	log.Printf("File accepted. ID: %s", id)
	return c.ProcessFile(id)
}

func (c *commandSubmit) ProcessFile(id string) error {
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
		if status.Error.Code != "" {
			return fmt.Errorf("%s. %s", status.Error.Code, status.Error.Message)
		}
		log.Printf("Status: %v", status.Status)
		switch status.Status {
		case "succeeded":
		case "running":
			return nil
		default:
			return fmt.Errorf("unknown status: %s", status.Status)
		}
		if time.Now().After(endTime) {
			return ErrTimeout
		}
		time.Sleep(5 * time.Second)
	}
}

func (c *commandSubmit) AnalysisResult(id string) (bool, error) {
	results, err := c.visionOne.SandboxAnalysisResults(id).Do()
	if err != nil {
		return false, err
	}
	log.Printf("Type: %s", results.Type)
	log.Printf("TrueFileType: %s", results.TrueFileType)
	log.Printf("RiskLevel: %s", results.RiskLevel)
	if len(results.DetectionNames) > 0 {
		log.Printf("DetectionNames: %s", strings.Join(results.DetectionNames, ", "))
	}
	if len(results.ThreatTypes) > 0 {
		log.Printf("ThreatTypes: %s", strings.Join(results.ThreatTypes, ", "))
	}
	return results.RiskLevel == "high", nil
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
	log.Println("Request Suspicious Objects")
	so, err := c.visionOne.SandboxSuspiciousObjects(id).Do()
	if err != nil {
		return err
	}
	soSHA1 := ListSHA1(so)
	for _, sha1 := range soSHA1 {
		log.Printf("Suspicious Object SHA1: %s", sha1)
	}
	soIP := ListIP(so)
	for _, ip := range soIP {
		log.Printf("Suspicious Object IP: %s", ip)
	}
	return nil
}

func (c *commandSubmit) PDFReport(id string) error {
	log.Println("Download report PDF")
	pdfFileName := id + ".pdf"
	if err := c.visionOne.SandboxDownloadResults(id).Store(pdfFileName); err != nil {
		return err
	}
	log.Printf("PDF report saved: %s", pdfFileName)
	return nil
}

func (c *commandSubmit) InvestigationPackage(id string) error {
	log.Println("Download Investigation Package")
	zipFileName := id + ".zip"
	if err := c.visionOne.SandboxInvestigationPackage(id).Store(zipFileName); err != nil {
		return err
	}
	log.Printf("Investigation Package Saved: %s", zipFileName)
	return nil
}

func (c *commandSubmit) Setup(name string) {
	c.baseCommand.Setup(name)
	c.fs.String(flagFileName, "", "Sample file path")
	c.fs.Duration(flagTimeout, 10*time.Minute, "Analisys timeout")
}
