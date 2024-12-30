package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/mpkondrashin/vone"
	"github.com/spf13/viper"
)

type commandHighRiskDevices struct {
	baseCommand
}

func newCommandHighRiskDevices() *commandHighRiskDevices {
	c := &commandHighRiskDevices{}
	c.Setup(cmdHighRiskDevices, "List high risk devices")
	c.fs.String(flagFilter, "", "id, deviceName, ip, os, riskScore, lastLogonUser; eq, and, or, not, ( ), gt, ge, le, lt")
	c.fs.String(flagOrderBy, "", "riskScore or deviceName + asc or desc")
	c.fs.Int(flagTop, 0, "Number of records displayed on a page. Possible values are 10, 50, 100 (default), 200, and 1000")
	return c
}

func (c *commandHighRiskDevices) Execute() error {
	list := c.visionOne.HighRiskDevices()
	filter := viper.GetString(flagFilter)
	if filter != "" {
		list.Filter(filter)
	}
	orderBy := viper.GetString(flagOrderBy)
	if orderBy != "" {
		list.OrderBy(orderBy)
	}
	top := viper.GetInt(flagTop)
	if top > 0 {
		t, err := vone.TopXMFromInt(top)
		if err != nil {
			return err
		}
		list.Top(t)
	}

	writer := csv.NewWriter(os.Stdout)
	dataCh := make(chan interface{})
	go func() {
		for row, err := range list.Range(context.TODO()) {
			if err != nil {
				panic(err)
			}
			dataCh <- row
		}
		close(dataCh)
	}()

	if err := gocsv.MarshalChan(dataCh, writer); err != nil {
		return fmt.Errorf("gocsv: %v", err)
	}
	return nil
}
