/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	get_endpoint_list.go - get endpoint data
*/

package vone

import (
	"context"
	"iter"
	"strings"
)

type (
	// EndpointListResponse -  response for get endpoint list query
	EndpointListResponse struct {
		Items      []EndpointListItem `json:"items"`
		Count      int                `json:"count"`
		TotalCount int                `json:"totalCount"`
		NextLink   string             `json:"nextLink"`
	}
	EndpointListItem struct {
		EndpointName          string        `json:"endpointName"`
		AgentGUID             string        `json:"agentGuid"`
		DisplayName           string        `json:"displayName"`
		OsName                string        `json:"osName"`
		OsVersion             string        `json:"osVersion"`
		OsKernelVersion       string        `json:"osKernelVersion"`
		OsArchitecture        string        `json:"osArchitecture"`
		LastUsedIP            string        `json:"lastUsedIp"`
		ServiceGatewayOrProxy string        `json:"serviceGatewayOrProxy"`
		CPUArchitecture       string        `json:"cpuArchitecture"`
		LastLoggedOnUser      string        `json:"lastLoggedOnUser"`
		IsolationStatus       string        `json:"isolationStatus"`
		IPAddresses           StringsSlice  `json:"ipAddresses"`
		SerialNumber          string        `json:"serialNumber"`
		EppAgent              EppAgentData  `json:"eppAgent"`
		EdrSensor             EdrSensorData `json:"edrSensor"`
	}

	PatternData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	Patterns     []PatternData
	EppAgentData struct {
		EndpointGroup         string        `json:"endpointGroup"`
		ProtectionManager     string        `json:"protectionManager"`
		PolicyName            string        `json:"policyName"`
		Status                string        `json:"status"`
		LastConnectedDateTime VisionOneTime `json:"lastConnectedDateTime"`
		Version               string        `json:"version"`
		LastScannedDateTime   VisionOneTime `json:"lastScannedDateTime"`
		ComponentVersion      string        `json:"componentVersion"`
		ComponentUpdatePolicy string        `json:"componentUpdatePolicy"`
		ComponentUpdateStatus string        `json:"componentUpdateStatus"`
		InstalledComponentIds StringsSlice  `json:"installedComponentIds"`
		Patterns              Patterns      `json:"patterns"`
	}
	EdrSensorData struct {
		EndpointGroup               string        `json:"endpointGroup"`
		Connectivity                string        `json:"connectivity"`
		Version                     string        `json:"version"`
		LastConnectedDateTime       VisionOneTime `json:"lastConnectedDateTime"`
		Status                      string        `json:"status"`
		AdvancedRiskTelemetryStatus string        `json:"advancedRiskTelemetryStatus"`
	}
)

// Convert Patterns the internal date as CSV string
func (p Patterns) MarshalCSV() (string, error) {
	var sb strings.Builder
	for _, pattern := range p {
		sb.WriteString(pattern.ID)
		sb.WriteString(",")
		sb.WriteString(pattern.Name)
		sb.WriteString("|")
	}
	return sb.String(), nil
}

// SearchEndPointDataFunc - search for endpoints
type GetEndPointListFunc struct {
	baseFunc
	Response EndpointListResponse
	top      int
}

// Filter - filter endpoints
func (f *GetEndPointListFunc) Filter(filter string) *GetEndPointListFunc {
	f.setHeader("TMV1-Filter", filter)
	//	if f.query != query {
	//		f.query = query
	f.Response.NextLink = ""
	//	}
	return f
}

// OrderBy - sort endpoints
func (f *GetEndPointListFunc) OrderBy(orderBy string) *GetEndPointListFunc {
	f.setParameter("orderBy", orderBy)
	return f
}

// Top - set limit for returned amount of items
func (f *GetEndPointListFunc) Top(t Top) *GetEndPointListFunc {
	f.setParameter("top", t.String())
	f.top = t.Int()
	return f
}

// Do - run request
func (f *GetEndPointListFunc) Do(ctx context.Context) (*EndpointListResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

// Iterate - get all endpoints matching query one by one. If callback returns
// non nil error, iteration is aborted and this error is returned
func (f *GetEndPointListFunc) Iterate(ctx context.Context,
	callback func(item *EndpointListItem) error) error {
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

// Range - iterator for all endpoints matching query (go 1.23 and later)
func (f *GetEndPointListFunc) Range(ctx context.Context) iter.Seq2[*EndpointListItem, error] {
	return func(yield func(*EndpointListItem, error) bool) {
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
func (v *VOne) EndPointList() *GetEndPointListFunc {
	f := &GetEndPointListFunc{
		top: 100,
	}
	f.baseFunc.init(v)
	return f

}

func (s *GetEndPointListFunc) uri() string {
	if s.Response.NextLink != "" {
		return s.Response.NextLink
	}
	return ""
}

func (s *GetEndPointListFunc) url() string {
	return "/v3.0/endpointSecurity/endpoints"
}

func (f *GetEndPointListFunc) responseStruct() any {
	return &f.Response
}
