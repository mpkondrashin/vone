/*
Trend Micro Vision One API SDK
(c) 2026 by Mikhail Kondrashin (mkondrashin@gmail.com)

# Vision One API

get_mobile_activity.go - get mobile activity
*/
package vone

import (
	"context"
	"time"
)

type (
	GetMobileActivityResponseItem struct {
		AppLabel              string    `json:"appLabel"`
		AppIsSystem           bool      `json:"appIsSystem"`
		AppPkgName            string    `json:"appPkgName"`
		AppPublicKeySha1      string    `json:"appPublicKeySha1"`
		AppSize               string    `json:"appSize"`
		AppVerCode            string    `json:"appVerCode"`
		EndpointGUID           string    `json:"endpointGuid"`
		EndpointHostName       string    `json:"endpointHostName"`
		EndpointIP            []string  `json:"endpointIp"`
		EndpointModel          string    `json:"endpointModel"`
		EventID               int       `json:"eventId"`
		EventSubID            int       `json:"eventSubId"`
		EventTime             int64     `json:"eventTime"`
		EventTimeDT           time.Time `json:"eventTimeDT"`
		FirstSeen             string    `json:"firstSeen"`
		FilterRiskLevel       string    `json:"filterRiskLevel"`
		LogonUser             []string  `json:"logonUser"`
		ObjectAppIsSystemApp  bool      `json:"objectAppIsSystemApp"`
		ObjectAppLabel        string    `json:"objectAppLabel"`
		ObjectAppPackageName  string    `json:"objectAppPackageName"`
		ObjectAppPublicKeySha1 string    `json:"objectAppPublicKeySha1"`
		ObjectAppSha256       string    `json:"objectAppSha256"`
		ObjectAppSize         string    `json:"objectAppSize"`
		ObjectAppVerCode      string    `json:"objectAppVerCode"`
		ObjectAppInstalledTime string    `json:"objectAppInstalledTime"`
		ObjectHostName        string    `json:"objectHostName"`
		OsName                string    `json:"osName"`
		OsVer                 string    `json:"osVer"`
		Pname                 string    `json:"pname"`
		ProductCode           string    `json:"productCode"`
		Pver                  string    `json:"pver"`
		Request               string    `json:"request"`
		Tags                  []string  `json:"tags"`
		UUID                  string    `json:"uuid"`
		UserType              string    `json:"userType"`
	}
	GetMobileActivityResponse struct {
		NextLink     string                         `json:"nextLink"`
		ProgressRate int                            `json:"progressRate"`
		Items        []GetMobileActivityResponseItem `json:"items"`
	}
)

// getMobileActivityRequest - search for mobile activities
type getMobileActivityRequest struct {
	baseRequest
	response GetMobileActivityResponse
	mode     Mode
}

// Do - run request
func (f *getMobileActivityRequest) Do(ctx context.Context) (*GetMobileActivityResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	//log.Printf("Got %d items", len(f.response.Items))
	return &f.response, nil
}

// GetMobileActivity - get new search for mobile activity data request
func (v *VOne) GetMobileActivity() *getMobileActivityRequest {
	f := &getMobileActivityRequest{}
	f.baseRequest.init(v)
	return f

}

// Mode - set mode
func (f *getMobileActivityRequest) Mode(mode Mode) *getMobileActivityRequest {
	f.mode = mode
	f.setParameter("mode", mode.String())
	f.response.NextLink = ""
	return f
}

// StartDateTime - set start date time
func (f *getMobileActivityRequest) StartDateTime(t time.Time) *getMobileActivityRequest {
	f.setParameter("startDateTime", t.Format(timeFormatZ))
	return f
}

// EndDateTime - set end date time
func (f *getMobileActivityRequest) EndDateTime(t time.Time) *getMobileActivityRequest {
	f.setParameter("endDateTime", t.Format(timeFormatZ))
	return f
}

// Query - set search query
func (f *getMobileActivityRequest) Query(filter string) *getMobileActivityRequest {
	f.setHeader("TMV1-Query", filter)
	f.response.NextLink = ""
	return f
}

// Top - set limit for returned items
func (f *getMobileActivityRequest) Top(t Top) *getMobileActivityRequest {
	f.setParameter("top", t.String())
	return f
}

func (f *getMobileActivityRequest) Select(selectString string) *getMobileActivityRequest {
	f.setParameter("select", selectString)
	return f
}

func (f *getMobileActivityRequest) isDone(resp *GetMobileActivityResponse) bool {
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

func (f *getMobileActivityRequest) url() string {
	if f.response.NextLink != "" {
		return f.response.NextLink
	}
	return "/v3.0/search/mobileActivities"
}

func (f *getMobileActivityRequest) uri() string {
	return f.response.NextLink
}

func (f *getMobileActivityRequest) responseStruct() any {
	return &f.response
}

// Paginator - get paginator for mobile activity data
func (f *getMobileActivityRequest) nextLink() string {
	return f.response.NextLink
}

func (f *getMobileActivityRequest) resetPagination() {
	f.response.NextLink = ""
}

func (f *getMobileActivityRequest) Paginator() *Paginator[
	GetMobileActivityResponse,
	GetMobileActivityResponseItem,
] {
	return NewPaginator(
		f,
		func(r *GetMobileActivityResponse) []GetMobileActivityResponseItem {
			return r.Items
		},
	)
}
