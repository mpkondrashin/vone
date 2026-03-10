/*
Trend Micro Vision One API SDK
(c) 2026 by Mikhail Kondrashin (mkondrashin@gmail.com)

# Vision One API

get_endpoint_activity.go - get endpoint activity
*/
package vone

import (
	"context"
	"time"
)

type (
	GetEndpointActivityResponseItem struct {
		Dpt                     int       `json:"dpt"`
		Dst                     string    `json:"dst"`
		EndpointGUID            string    `json:"endpointGuid"`
		EndpointHostName        string    `json:"endpointHostName"`
		EndpointIP              []string  `json:"endpointIp"`
		EventID                 string    `json:"eventId"`
		EventSubID              int       `json:"eventSubId"`
		ObjectIntegrityLevel    int       `json:"objectIntegrityLevel"`
		ObjectTrueType          int       `json:"objectTrueType"`
		ObjectSubTrueType       int       `json:"objectSubTrueType"`
		WinEventID              int       `json:"winEventId"`
		EventTime               int64     `json:"eventTime"`
		EventTimeDT             time.Time `json:"eventTimeDT"`
		HostName                string    `json:"hostName"`
		LogonUser               []string  `json:"logonUser"`
		ObjectCmd               string    `json:"objectCmd"`
		ObjectFileHashSha1      string    `json:"objectFileHashSha1"`
		ObjectFilePath          string    `json:"objectFilePath"`
		ObjectHostName          string    `json:"objectHostName"`
		ObjectIP                string    `json:"objectIp"`
		ObjectIps               []string  `json:"objectIps"`
		ObjectPort              int       `json:"objectPort"`
		ObjectRegistryData      string    `json:"objectRegistryData"`
		ObjectRegistryKeyHandle string    `json:"objectRegistryKeyHandle"`
		ObjectRegistryValue     string    `json:"objectRegistryValue"`
		ObjectSigner            []string  `json:"objectSigner"`
		ObjectSignerValid       []bool    `json:"objectSignerValid"`
		ObjectUser              string    `json:"objectUser"`
		Os                      string    `json:"os"`
		ParentCmd               string    `json:"parentCmd"`
		ParentFileHashSha1      string    `json:"parentFileHashSha1"`
		ParentFilePath          string    `json:"parentFilePath"`
		ProcessCmd              string    `json:"processCmd"`
		ProcessFileHashSha1     string    `json:"processFileHashSha1"`
		ProcessFilePath         string    `json:"processFilePath"`
		Request                 string    `json:"request"`
		SearchDL                string    `json:"searchDL"`
		Spt                     int       `json:"spt"`
		Src                     string    `json:"src"`
		SrcFileHashSha1         string    `json:"srcFileHashSha1"`
		SrcFilePath             string    `json:"srcFilePath"`
		Tags                    []string  `json:"tags"`
		UUID                    string    `json:"uuid"`
	}
	GetEndpointActivityResponse struct {
		NextLink     string                            `json:"nextLink"`
		ProgressRate int                               `json:"progressRate"`
		Items        []GetEndpointActivityResponseItem `json:"items"`
	}
)

// searchEndPointDataRequest - search for endpoints
type getEndpointActivityRequest struct {
	baseRequest
	response GetEndpointActivityResponse
	mode     Mode
}

// Do - run request
func (f *getEndpointActivityRequest) Do(ctx context.Context) (*GetEndpointActivityResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	//log.Printf("Got %d items", len(f.response.Items))
	return &f.response, nil
}

// GetEndpointActivity - get new search for endpoint activity data request
func (v *VOne) GetEndpointActivity() *getEndpointActivityRequest {
	f := &getEndpointActivityRequest{}
	f.baseRequest.init(v)
	return f

}

// Mode - set mode
func (f *getEndpointActivityRequest) Mode(mode Mode) *getEndpointActivityRequest {
	f.mode = mode
	f.setParameter("mode", mode.String())
	f.response.NextLink = ""
	return f
}

// StartDateTime - set start date time
func (f *getEndpointActivityRequest) StartDateTime(t time.Time) *getEndpointActivityRequest {
	f.setParameter("startDateTime", t.Format(timeFormatZ))
	return f
}

// EndDateTime - set end date time
func (f *getEndpointActivityRequest) EndDateTime(t time.Time) *getEndpointActivityRequest {
	f.setParameter("endDateTime", t.Format(timeFormatZ))
	return f
}

// Query - set search query
func (f *getEndpointActivityRequest) Query(filter string) *getEndpointActivityRequest {
	f.setHeader("TMV1-Query", filter)
	f.response.NextLink = ""
	return f
}

// Top - set limit for returned items
func (f *getEndpointActivityRequest) Top(t Top) *getEndpointActivityRequest {
	f.setParameter("top", t.String())
	return f
}

func (f *getEndpointActivityRequest) Select(selectString string) *getEndpointActivityRequest {
	f.setParameter("select", selectString)
	return f
}

func (f *getEndpointActivityRequest) isDone(resp *GetEndpointActivityResponse) bool {
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

func (f *getEndpointActivityRequest) url() string {
	if f.response.NextLink != "" {
		return f.response.NextLink
	}
	return "/v3.0/search/endpointActivities"
}

func (f *getEndpointActivityRequest) uri() string {
	return f.response.NextLink
}

func (f *getEndpointActivityRequest) responseStruct() any {
	return &f.response
}

// Paginator - get paginator for endpoint activity data
func (f *getEndpointActivityRequest) nextLink() string {
	return f.response.NextLink
}

func (f *getEndpointActivityRequest) resetPagination() {
	f.response.NextLink = ""
}

func (f *getEndpointActivityRequest) Paginator() *Paginator[
	GetEndpointActivityResponse,
	GetEndpointActivityResponseItem,
] {
	return NewPaginator(
		f,
		func(r *GetEndpointActivityResponse) []GetEndpointActivityResponseItem {
			return r.Items
		},
	)
}
