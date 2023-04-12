package main

import (
	"time"

	"github.com/spf13/viper"
)

type commandPDF struct {
	commandSubmit
}

func newCommandPDF() *commandPDF {
	c := &commandPDF{}
	c.Setup(cmdPDF)
	return c
}

func (c *commandPDF) Execute() error {
	id := viper.GetString(flagID)
	return c.ProcessObject(id)
}

func (c *commandPDF) ProcessObject(id string) error {
	if err := c.WaitForResult(id); err != nil {
		return err
	}
	return c.PDFReport(id)
}

/*
func (c *commandPDF) WaitForResult(id string) error {
	timeout := viper.GetDuration(flagTimeout)
	endTime := time.Now().Add(timeout)
	for {
		status, err := c.visionOne.SandboxSubmissionStatus(id).Do(context.TODO())
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

func (c *commandPDF) PDFReport(id string) error {
	log.Printf("%s Download report PDF", id)
	pdfFileName := id + ".pdf"
	if err := c.visionOne.SandboxDownloadResults(id).Store(context.TODO(), pdfFileName); err != nil {
		return err
	}
	log.Printf("%s PDF report saved: %s", id, pdfFileName)
	return nil
}
*/
func (c *commandPDF) Setup(name string) {
	c.baseCommand.Setup(name)
	c.fs.String(flagID, "", "Sample file path")
	c.fs.Duration(flagTimeout, 10*time.Minute, "Analysis timeout")
}
