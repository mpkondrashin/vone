package vone

import (
	"context"
	"sync"
)

type RegionalDomain struct {
	Region string
	Domain string
}

var RegionalDomains = []RegionalDomain{
	{"Australia", "api.au.xdr.trendmicro.com"},
	{"India", "api.in.xdr.trendmicro.com"},
	{"Japan", "api.xdr.trendmicro.co.jp"},
	{"Singapore", "api.sg.xdr.trendmicro.com"},
	{"United States", "api.xdr.trendmicro.com"},
	{"Europe", "api.eu.xdr.trendmicro.com"},
}

// DetectVisionOneDomain return correct domain for given token or empty string
func DetectVisionOneDomain(ctx context.Context, token string) (result string) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	for _, rd := range RegionalDomains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			_, err := NewVOne(d, token).CheckConnection().Do(ctx)
			if err != nil {
				return
			}
			cancel()
			result = d
		}(rd.Domain)
	}
	wg.Wait()
	return
}
