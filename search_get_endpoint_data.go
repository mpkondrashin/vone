/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	search_get_endpoint_data.go - get endpoint data
*/

package vone

import (
	"context"
	"iter"
)

// SearchEndPointDataResponseItem - get endpoint data response for each endpoint
type SearchEndPointDataResponseItem struct {
	AgentGUID    string `json:"agentGuid"`
	LoginAccount struct {
		Value           []string      `json:"value"`
		UpdatedDateTime VisionOneTime `json:"updatedDateTime"`
	} `json:"loginAccount"`
	EndpointName struct {
		Value           string        `json:"value"`
		UpdatedDateTime VisionOneTime `json:"updatedDateTime"`
	} `json:"endpointName"`
	MacAddress struct {
		Value           []string      `json:"value"`
		UpdatedDateTime VisionOneTime `json:"updatedDateTime"`
	} `json:"macAddress"`
	IP struct {
		Value           []string      `json:"value"`
		UpdatedDateTime VisionOneTime `json:"updatedDateTime"`
	} `json:"ip"`
	OsName                string   `json:"osName"`
	OsVersion             string   `json:"osVersion"`
	OsDescription         string   `json:"osDescription"`
	ProductCode           string   `json:"productCode"`
	InstalledProductCodes []string `json:"installedProductCodes"`
}

// SearchEndPointDataResponse - get endpoint data response
type SearchEndPointDataResponse struct {
	Items    []SearchEndPointDataResponseItem `json:"items"`
	NextLink string                           `json:"nextLink"`
}

// SearchEndPointDataFunc - search for endpoints
type SearchEndPointDataFunc struct {
	baseFunc
	Response SearchEndPointDataResponse
	//query    string
}

// Do - run request
func (f *SearchEndPointDataFunc) Do(ctx context.Context) (*SearchEndPointDataResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

// Iterate - get all endpoints matching query one by one. If callback returns
// non nil error, iteration is aborted and this error is returned
func (f *SearchEndPointDataFunc) Iterate(ctx context.Context,
	callback func(item *SearchEndPointDataResponseItem) error) error {
	for {
		response, err := f.Do(ctx)
		if err != nil {
			return err
		}
		for n := range response.Items {
			if err := callback(&response.Items[n]); err != nil {
				return err
			}
		}
		if response.NextLink == "" {
			break
		}
	}
	return nil
}

// Range - iterator for all endpoints matching query (go 1.23 and later)
func (f *SearchEndPointDataFunc) Range(ctx context.Context) iter.Seq2[*SearchEndPointDataResponseItem, error] {
	return func(yield func(*SearchEndPointDataResponseItem, error) bool) {
		for {
			response, err := f.Do(ctx)
			if err != nil {
				yield(nil, err)
				return
			}
			for n := range response.Items {
				if !yield(&response.Items[n], nil) {
					return
				}
			}
			if response.NextLink == "" {
				break
			}
		}
	}
}

// SearchEndPointData - get new search for endpoint data function
func (v *VOne) SearchEndPointData() *SearchEndPointDataFunc {
	f := &SearchEndPointDataFunc{}
	f.baseFunc.init(v)
	return f

}

// Query - set search query
func (f *SearchEndPointDataFunc) Query(filter Filter) *SearchEndPointDataFunc {
	f.setHeader("TMV1-Query", filter.Build())
	//	if f.query != query {
	//		f.query = query
	f.Response.NextLink = ""
	//	}
	return f
}

// Query - set search query
func (f *SearchEndPointDataFunc) QueryString(query string) *SearchEndPointDataFunc {
	f.setHeader("TMV1-Query", query)
	//	if f.query != query {
	//		f.query = query
	f.Response.NextLink = ""
	//	}
	return f
}

// top - set limit for returned amount of items
func (f *SearchEndPointDataFunc) Top(t Top) *SearchEndPointDataFunc {
	f.setParameter("top", t.String())
	return f
}

func (s *SearchEndPointDataFunc) url() string {
	if s.Response.NextLink != "" {
		return s.Response.NextLink
	}
	return "/v3.0/eiqs/endpoints"
}

//func (s *SearchEndPointDataFunc) populate(req *http.Request) {
//	s.baseFunc.populate(req)
//	req.Header.Set("TMV1-Query", s.query)
//}

func (f *SearchEndPointDataFunc) responseStruct() any {
	return &f.Response
}
