package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/mpkondrashin/vone"
	"github.com/spf13/viper"
)

type commandGetEndpointData struct {
	baseCommand
}

func newCommandGetEndpointData() *commandGetEndpointData {
	c := &commandGetEndpointData{}
	c.Setup(cmdGetEndpointData, "get endpoints data")
	c.fs.String(flagQuery, "", "Endpoints filter. Parameters:"+
		" agentGuid, loginAccount, endpointName, macAddress, ip,"+
		" osName (Linux, Windows, macOS, macOSX), osVersion,"+
		" productCode (sao, sds, xes), installedProductCodes."+
		" Operators: eq, and, or, not.")
	c.fs.Int(flagTop, 0, "Response limit. Possible values are 50 (default), 100, and 200. If omited, all data is downloaded")
	return c
}

func (c *commandGetEndpointData) Execute() error {
	search := c.visionOne.SearchEndPointData()
	query := viper.GetString(flagQuery)
	if query == "" {
		log.Fatalf("--%s parameter can not be empty", flagQuery)
	}
	search.QueryString(query)
	topAmount := viper.GetInt(flagTop)
	if topAmount != 0 {
		top, err := vone.TopFromInt(topAmount)
		if err != nil {
			return err
		}
		search.Top(top)
		response, err := search.Do(context.TODO())
		if err != nil {
			return err
		}
		s, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			return err
		}
		fmt.Println(string(s))
		return nil
	}
	writer := csv.NewWriter(os.Stdout)
	dataCh := make(chan interface{})
	go func() {
		for row, err := range search.Range(context.TODO()) {
			if err != nil {
				panic(err)
			}
			dataCh <- row
		}
		close(dataCh)
	}()
	err := gocsv.MarshalChan(dataCh, writer)
	if err != nil {
		return fmt.Errorf("gocsv: %v", err)
	}
	return nil
}
