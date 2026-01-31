package vone

import (
	"context"
	"errors"
	"net/http"
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
func DetectVisionOneDomain(
	ctx context.Context,
	token string,
	modifier func(*http.Transport),
) (string, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type result struct {
		domain string
		err    error
	}

	results := make(chan result, len(RegionalDomains))

	for _, rd := range RegionalDomains {
		domain := rd.Domain

		go func() {
			vOne := NewVOne(domain, token)
			if modifier != nil {
				vOne.AddTransportModifier(modifier)
			}

			_, err := vOne.CheckConnection().Do(ctx)
			if err != nil {
				// если это не API-ошибка — возвращаем её
				var apiErr *Error
				if !errors.As(err, &apiErr) {
					results <- result{err: err}
				}
				return
			}

			// успех
			results <- result{domain: domain}
		}()
	}

	var firstErr error

	for range RegionalDomains {
		select {
		case r := <-results:
			if r.domain != "" {
				cancel() // останавливаем остальные
				return r.domain, nil
			}
			if r.err != nil && firstErr == nil {
				firstErr = r.err
			}
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	if firstErr != nil {
		return "", firstErr
	}

	return "", errors.New("vision one domain not detected")
}
