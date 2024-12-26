package main

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/dnlo/struct2csv"
	"github.com/mpkondrashin/vone"
	"github.com/spf13/viper"
)

type commandListEndpoints struct {
	baseCommand
}

func newCommandListEndpoints() *commandListEndpoints {
	c := &commandListEndpoints{}
	c.Setup(cmdListGetEndpoints, "List endpoints")
	c.fs.String(flagFilter, "", "https://automation.trendmicro.com/xdr/api-v3/#tag/Endpoint-Security/paths/~1v3.0~1endpointSecurity~1endpoints/get")
	c.fs.String(flagOrderBy, "", "agentGuid, edrSensorLastConnectedDateTime, eppAgentLastConnectedDateTime, eppAgentLastScannedDateTime, +asc or desc")
	c.fs.Int(flagTop, 0, "Number of records displayed on a page. Possible values are 50, 100 (default), 200, 500, and 1000")
	return c
}

func (c *commandListEndpoints) Execute() error {
	list := c.visionOne.EndPointList()
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
		t, err := vone.TopFromInt(top)
		if err != nil {
			return err
		}
		list.Top(t)
	}

	w := csv.NewWriter(os.Stdout)

	data := vone.EndpointListItem{}
	enc := struct2csv.New()
	// get the column names first
	colhdrs, err := enc.GetColNames(data)
	if err != nil {
		return err
	}
	if err := w.Write(colhdrs); err != nil {
		return err
	}

	for item, err := range list.Range(context.TODO()) {
		if err != nil {
			return err
		}
		row, err := enc.GetRow(*item)
		if err != nil {
			return err
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}
