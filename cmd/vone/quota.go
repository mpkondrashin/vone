package main

import (
	"context"
	"log"
)

type commandQuota struct {
	baseCommand
}

func newCommandQuota() *commandQuota {
	c := &commandQuota{}
	c.Setup(cmdQuota, "Get allocated daily files/URLs submission limit")
	return c
}

func (c *commandQuota) Execute() error {
	quota, err := c.visionOne.SandboxDailyReserve().Do(context.TODO())
	if err != nil {
		return err
	}
	log.Printf("Daily quota: %d", quota.SubmissionReserveCount)
	log.Printf("Quota left: %d", quota.SubmissionRemainingCount)
	log.Printf("Submitted objects: %d (files: %d, urls: %d)", quota.SubmissionCount, quota.SubmissionCountDetail.FileCount, quota.SubmissionCountDetail.URLCount)
	log.Printf("Submitted unsupported (not accounted in quota) objects: %d (files: %d, urls: %d) ", quota.SubmissionExemptionCount, quota.SubmissionCountDetail.FileExemptionCount, quota.SubmissionCountDetail.URLExemptionCount)
	return nil
}
