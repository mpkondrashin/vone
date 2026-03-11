/*
Trend Micro Vision One API SDK
(c) 2026 by Mikhail Kondrashin (mkondrashin@gmail.com)

# Vision One API

get_network_activity.go - get network activity
*/
package vone

import (
	"context"
	"time"
)

type (
	GetNetworkActivityResponseItem struct {
		Application         string `json:"application"`
		Act                 string `json:"act"`
		ClientIP            string `json:"clientIp"`
		CompanyName         string `json:"companyName"`
		Dpt                 int    `json:"dpt"`
		DetectionType       string `json:"detectionType"`
		DeviceGUID          string `json:"deviceGuid"`
		Dst                 string `json:"dst"`
		Duration            int    `json:"duration"`
		EndpointGUID        string `json:"endpointGuid"`
		EndpointHostName    string `json:"endpointHostName"`
		EventName           string `json:"eventName"`
		EventSubName        string `json:"eventSubName"`
		EventTime           int64  `json:"eventTime"`
		EventTimeDT         time.Time `json:"eventTimeDT"`
		FileHash            string `json:"fileHash"`
		FileHashSha256      string `json:"fileHashSha256"`
		FileName            string `json:"fileName"`
		FileSize            string `json:"fileSize"`
		FileType            string `json:"fileType"`
		MalName             string `json:"malName"`
		MimeType            string `json:"mimeType"`
		ObjectId            string `json:"objectId"`
		OsName              string `json:"osName"`
		Pname               string `json:"pname"`
		PolicyUuid          string `json:"policyUuid"`
		Pver                string `json:"pver"`
		PrincipalName       string `json:"principalName"`
		Profile             string `json:"profile"`
		Request             string `json:"request"`
		RequestBase         string `json:"requestBase"`
		RequestMimeType     string `json:"requestMimeType"`
		RequestMethod       string `json:"requestMethod"`
		RuleName            string `json:"ruleName"`
		RuleType            string `json:"ruleType"`
		RuleUuid            string `json:"ruleUuid"`
		Score               int    `json:"score"`
		Sender              string `json:"sender"`
		ServerProtocol      string `json:"serverProtocol"`
		ServerTls           string `json:"serverTls"`
		Spt                 int    `json:"spt"`
		Src                 string `json:"src"`
		Start               int64  `json:"start"`
		Suid                string `json:"suid"`
		TenantGuid          string `json:"tenantGuid"`
		UserAgent           string `json:"userAgent"`
		UserDepartment      string `json:"userDepartment"`
		UserDomain          string `json:"userDomain"`
		Rt                  int64  `json:"rt"`
	}
	GetNetworkActivityResponse struct {
		NextLink     string                          `json:"nextLink"`
		ProgressRate int                             `json:"progressRate"`
		Items        []GetNetworkActivityResponseItem `json:"items"`
	}
)

// getNetworkActivityRequest - search for network activities
type getNetworkActivityRequest struct {
	baseRequest
	response GetNetworkActivityResponse
	mode     Mode
}

// Do - run request
func (f *getNetworkActivityRequest) Do(ctx context.Context) (*GetNetworkActivityResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	//log.Printf("Got %d items", len(f.response.Items))
	return &f.response, nil
}

// GetNetworkActivity - get new search for network activity data request
func (v *VOne) GetNetworkActivity() *getNetworkActivityRequest {
	f := &getNetworkActivityRequest{}
	f.baseRequest.init(v)
	return f

}

// Mode - set mode
func (f *getNetworkActivityRequest) Mode(mode Mode) *getNetworkActivityRequest {
	f.mode = mode
	f.setParameter("mode", mode.String())
	f.response.NextLink = ""
	return f
}

// StartDateTime - set start date time
func (f *getNetworkActivityRequest) StartDateTime(t time.Time) *getNetworkActivityRequest {
	f.setParameter("startDateTime", t.Format(timeFormatZ))
	return f
}

// EndDateTime - set end date time
func (f *getNetworkActivityRequest) EndDateTime(t time.Time) *getNetworkActivityRequest {
	f.setParameter("endDateTime", t.Format(timeFormatZ))
	return f
}

// Query - set search query
func (f *getNetworkActivityRequest) Query(filter string) *getNetworkActivityRequest {
	f.setHeader("TMV1-Query", filter)
	f.response.NextLink = ""
	return f
}

// Top - set limit for returned items
func (f *getNetworkActivityRequest) Top(t Top) *getNetworkActivityRequest {
	f.setParameter("top", t.String())
	return f
}

func (f *getNetworkActivityRequest) Select(selectString string) *getNetworkActivityRequest {
	f.setParameter("select", selectString)
	return f
}

func (f *getNetworkActivityRequest) isDone(resp *GetNetworkActivityResponse) bool {
	switch f.mode {
	case ModePerformance:
		return resp.ProgressRate >= 100 && resp.NextLink == ""

	case ModeCountOnly:
		return true

	case ModeDefault:
		fallthrough
	default:
		return resp.NextLink == ""
	}
}

func (f *getNetworkActivityRequest) url() string {
	if f.response.NextLink != "" {
		return f.response.NextLink
	}
	return "/v3.0/search/networkActivities"
}

func (f *getNetworkActivityRequest) uri() string {
	return f.response.NextLink
}

func (f *getNetworkActivityRequest) responseStruct() any {
	return &f.response
}

// Paginator - get paginator for network activity data
func (f *getNetworkActivityRequest) nextLink() string {
	return f.response.NextLink
}

func (f *getNetworkActivityRequest) resetPagination() {
	f.response.NextLink = ""
}

func (f *getNetworkActivityRequest) Paginator() *Paginator[
	GetNetworkActivityResponse,
	GetNetworkActivityResponseItem,
] {
	return NewPaginator(
		f,
		func(r *GetNetworkActivityResponse) []GetNetworkActivityResponseItem {
			return r.Items
		},
	)
}
