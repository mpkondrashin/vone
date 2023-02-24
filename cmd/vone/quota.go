package main

import "log"

type commandQuota struct {
	baseCommand
}

func newCommandQuota() *commandQuota {
	c := &commandQuota{}
	c.Setup(cmdQuota)
	return c
}

func (c *commandQuota) Execute() error {
	quota, err := c.visionOne.SandboxDailyReserve().Do()
	if err != nil {
		return err
	}
	log.Printf("Daily quota: %d", quota.SubmissionReserveCount)
	log.Printf("Quota left: %d", quota.SubmissionRemainingCount)
	log.Printf("Submitted objects: %d (files: %d, urls: %d)", quota.SubmissionCount, quota.SubmissionCountDetail.FileCount, quota.SubmissionCountDetail.URLCount)
	log.Printf("Submitted unsupported (not accounted in quota) objects: %d (files: %d, urls: %d) ", quota.SubmissionExemptionCount, quota.SubmissionCountDetail.FileExemptionCount, quota.SubmissionCountDetail.URLExemptionCount)
	return nil
}

/*
	log.Printf("%s Daily quota: %d", id, headers.SubmissionReserveCount)
	note := ""
	if headers.SubmissionRemainingCount == 0 {
		note = " (consider to add Credits to sandbox feature)"
	}
	log.Printf("%s Quota left: %d%s", id, headers.SubmissionRemainingCount, note)
	log.Printf("%s Today submissions: %d", id, headers.SubmissionCount)
	log.Printf("%s Today submissions of unsupported files: %d (not accounted in quota)", id, headers.SubmissionExemptionCount)
}
*/
