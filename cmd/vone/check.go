package main

import (
	"context"
	"log"
)

type commandCheck struct {
	baseCommand
}

func newCommandCheck() *commandCheck {
	c := &commandCheck{}
	c.Setup(cmdCheck, "Check Vision One token and connectivity")
	return c
}

func (c *commandCheck) Execute() error {
	check, err := c.visionOne.CheckConnection().Do(context.TODO())
	if err != nil {
		return err
	}
	log.Printf("Status: %s", check.Status)
	return nil
}
