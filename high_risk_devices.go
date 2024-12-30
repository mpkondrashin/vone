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
		ID            string   `json:"id"`
		DeviceName    string   `json:"deviceName"`
		Os            string   `json:"os"`
		RiskScore     int      `json:"riskScore"`
		LastLogonUser string   `json:"lastLogonUser"`
		IP            []string `json:"ip"`
	}
)

// SearchEndPointDataFunc - search for endpoints
type HighRiskDevicesFunc struct {
	baseFunc
	Response HighRiskDevicesResponse
	top      int
}

// Filter - filter endpoints
func (f *HighRiskDevicesFunc) Filter(filter string) *HighRiskDevicesFunc {
	f.setHeader("TMV1-Filter", filter)
	//	if f.query != query {
	//		f.query = query
	f.Response.NextLink = ""
	//	}
	return f
}

// OrderBy - sort endpoints
func (f *HighRiskDevicesFunc) OrderBy(orderBy string) *HighRiskDevicesFunc {
	f.setParameter("orderBy", orderBy)
	return f
}

// Top - set limit for returned amount of items
func (f *HighRiskDevicesFunc) Top(t TopXM) *HighRiskDevicesFunc {
	f.setParameter("top", t.String())
	f.top = t.Int()
	return f
}

// Do - run request
func (f *HighRiskDevicesFunc) Do(ctx context.Context) (*HighRiskDevicesResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

// Iterate - get all endpoints matching query one by one. If callback returns
// non nil error, iteration is aborted and this error is returned
func (f *HighRiskDevicesFunc) Iterate(ctx context.Context,
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
func (f *HighRiskDevicesFunc) Range(ctx context.Context) iter.Seq2[*HighRiskDevicesItem, error] {
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
func (v *VOne) HighRiskDevices() *HighRiskDevicesFunc {
	f := &HighRiskDevicesFunc{
		top: 100,
	}
	f.baseFunc.init(v)
	return f

}

func (s *HighRiskDevicesFunc) uri() string {
	if s.Response.NextLink != "" {
		return s.Response.NextLink
	}
	return ""
}

func (s *HighRiskDevicesFunc) url() string {
	return "/v3.0/asrm/highRiskDevices"
}

func (f *HighRiskDevicesFunc) responseStruct() any {
	return &f.Response
}
