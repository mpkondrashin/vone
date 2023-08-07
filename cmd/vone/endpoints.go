package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mpkondrashin/vone"
	"github.com/spf13/viper"
)

type commandGetEndpointData struct {
	baseCommand
}

func newCommandGetEndpointData() *commandGetEndpointData {
	c := &commandGetEndpointData{}
	c.Setup(cmdGetEndpointData)
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
		log.Fatalf("%s parameter can not be empty", flagQuery)
	}
	search.Query(query)
	top := vone.Top50
	topAmount := viper.GetInt(flagTop)
	if topAmount != 0 {
		var err error
		top, err = vone.TopFromInt(topAmount)
		if err != nil {
			log.Fatal(err)
		}
		search.Top(top)
	}
	first := true
	fmt.Println("[")
	for {
		response, err := search.Do(context.TODO())
		if err != nil {
			return err
		}
		for _, item := range response.Items {
			s, err := json.MarshalIndent(item, "    ", "    ")
			if err != nil {
				log.Fatal(err)
			}
			if first {
				fmt.Printf("    %s", string(s))
				first = false
			} else {
				fmt.Printf(",\n    %s", string(s))
			}
		}
		if topAmount != 0 {
			// If used did not provided "top" parameter limiting amount of data we will stop here
			break
		}
		if response.NextLink == "" {
			break
		}
	}
	fmt.Println("\n]")
	return nil
}
