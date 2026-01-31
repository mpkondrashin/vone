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

type (
	HighRiskDevicesResponse struct {
		Items      []HighRiskDevicesItem `json:"items"`
		Count      int                   `json:"count"`
		TotalCount int                   `json:"totalCount"`
		NextLink   string                `json:"nextLink"`
	}
	HighRiskDevicesItem struct {
		ID            string       `json:"id"`
		DeviceName    string       `json:"deviceName"`
		Os            string       `json:"os"`
		RiskScore     int          `json:"riskScore"`
		LastLogonUser string       `json:"lastLogonUser"`
		IP            StringsSlice `json:"ip"`
	}
)

// SearchEndPointDataFunc - search for endpoints
type HighRiskDevicesRequest struct {
	baseRequest
	response HighRiskDevicesResponse
	top      int
}

// Filter - filter endpoints
func (f *HighRiskDevicesRequest) Filter(filter string) *HighRiskDevicesRequest {
	f.setHeader("TMV1-Filter", filter)
	//	if f.query != query {
	//		f.query = query
	f.response.NextLink = ""
	//	}
	return f
}

// OrderBy - sort endpoints
func (f *HighRiskDevicesRequest) OrderBy(orderBy string) *HighRiskDevicesRequest {
	f.setParameter("orderBy", orderBy)
	return f
}

// Top - set limit for returned amount of items
func (f *HighRiskDevicesRequest) Top(t TopXM) *HighRiskDevicesRequest {
	f.setParameter("top", t.String())
	f.top = t.Int()
	return f
}

// Do - run request
func (f *HighRiskDevicesRequest) Do(ctx context.Context) (*HighRiskDevicesResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.response, nil
}

// Iterate - get all endpoints matching query one by one. If callback returns
// non nil error, iteration is aborted and this error is returned
func (f *HighRiskDevicesRequest) Iterate(ctx context.Context,
	callback func(item *HighRiskDevicesItem) error) error {
	for {
		response, err := f.Do(ctx)
		//fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXX:", err, response.Count, response.TotalCount) ////, /response, err)
		if err != nil {
			return err
		}
		//return errors.New("Quit")
		for n := range response.Items {
			//	fmt.Println("XXX call callback")
			if err := callback(&response.Items[n]); err != nil {
				return err
			}
		}
		if response.NextLink == "" {
			break
		}
		if response.Count != f.top {
			break
		}
	}
	return nil
}

// Range - iterator for all devices matching query (go 1.23 and later)
func (f *HighRiskDevicesRequest) Range(ctx context.Context) iter.Seq2[*HighRiskDevicesItem, error] {
	return func(yield func(*HighRiskDevicesItem, error) bool) {
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
			if response.Count != f.top {
				break
			}
		}
	}
}

// SearchEndPointData - get new search for endpoint data function
func (v *VOne) HighRiskDevices() *HighRiskDevicesRequest {
	f := &HighRiskDevicesRequest{
		top: 100,
	}
	f.baseRequest.init(v)
	return f

}

func (s *HighRiskDevicesRequest) uri() string {
	if s.response.NextLink != "" {
		return s.response.NextLink
	}
	return ""
}

func (s *HighRiskDevicesRequest) url() string {
	return "/v3.0/asrm/highRiskDevices"
}

func (f *HighRiskDevicesRequest) responseStruct() any {
	return &f.response
}
