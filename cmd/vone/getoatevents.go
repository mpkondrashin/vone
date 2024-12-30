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

type commandGetOATEvents struct {
	baseCommand
}

func newCommandGetOATEvents() *commandGetOATEvents {
	c := &commandGetOATEvents{}
	c.Setup(cmdGetOATEvents, "Get Observed Attack Techniques events")
	c.fs.String(flagFilter, "", "Events filter. Parameters:"+
		" uuid, riskLevel (undefined, info, low, medium, high, critical), filterName"+
		" filterMitreTacticId, filterMitreTechniqueId, endpointName, agentGuid, endpointIp"+
		"productCode, containerName. Operations: eq, and, or, not, ()")

	c.fs.String(flagDetectedStart, "", "The start of the event detection data retrieval time range in ISO 8601 format.")
	c.fs.String(flagDetectedEnd, "", "The end of the event detection data retrieval time range in ISO 8601 format. Default: The time you make the request.")

	c.fs.String(flagIngestedStart, "", "The beginning of the data ingestion time range in ISO 8601 format.")
	c.fs.String(flagIngestedEnd, "", "The end of the data ingestion time range in ISO 8601 format.")

	c.fs.Int(flagTop, 0, "Response limit. Possible values are 50 (default), 100, and 200. If omited, all data is downloaded")
	return c
}

func VOneTimeFromString(s string) (vone.VOneTime, error) {
	var v vone.VOneTime
	err := (&v).UnmarshalJSON([]byte(s))
	if err != nil {
		return vone.VOneTime{}, err
	}
	return v, nil
}

func (c *commandGetOATEvents) Execute() error {
	events := c.visionOne.GetOATEvents()
	filter := viper.GetString(flagFilter)
	if filter != "" {
		log.Println("Filter:", filter)
		events.Filter(filter)

	}

	detectedStart := viper.GetString(flagDetectedStart)
	if detectedStart != "" {
		t, err := VOneTimeFromString(detectedStart)
		if err != nil {
			return err
		}
		events.DetectedStart(t)
	}

	detectedEnd := viper.GetString(flagDetectedEnd)
	if detectedEnd != "" {
		t, err := VOneTimeFromString(detectedEnd)
		if err != nil {
			return err
		}
		events.DetectedEnd(t)
	}

	ingestedStart := viper.GetString(flagIngestedStart)
	if ingestedStart != "" {
		t, err := VOneTimeFromString(ingestedStart)
		if err != nil {
			return err
		}
		events.IngestedStart(t)
	}

	ingestedEnd := viper.GetString(flagIngestedEnd)
	if ingestedEnd != "" {
		t, err := VOneTimeFromString(ingestedEnd)
		if err != nil {
			return err
		}
		events.IngestedEnd(t)
	}

	topAmount := viper.GetInt(flagTop)
	if topAmount != 0 {
		top, err := vone.TopFromInt(topAmount)
		if err != nil {
			return err
		}
		events.Top(top)
		response, err := events.Do(context.TODO())
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
		for row, err := range events.Range(context.TODO()) {
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
