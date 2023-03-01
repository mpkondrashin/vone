package vone

import (
	"context"
	"sync"
)

var domains = []string{
	"api.au.xdr.trendmicro.com",
	"api.in.xdr.trendmicro.com",
	"api.xdr.trendmicro.co.jp",
	"api.sg.xdr.trendmicro.com",
	"api.xdr.trendmicro.com",
	"api.eu.xdr.trendmicro.com",
}

func DetectVisionOneDomain(ctx context.Context, token string) (result string) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	for _, d := range domains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			_, err := NewVOne(d, token).CheckConnection().Do(ctx)
			if err != nil {
				return
			}
			cancel()
			result = d
		}(d)
	}
	wg.Wait()
	return
}
