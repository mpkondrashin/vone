/*
	Trend Micro Vision One API SDK
	(c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	it_add_exception.go - add exceptions for Threat Intelligance
*/

package vone

import (
	"context"
	"iter"
)

type (
	ObservedAttackTechniquesEventsResponse struct {
		TotalCount int                                  `json:"totalCount"`
		Count      int                                  `json:"count"`
		Items      []ObservedAttackTechniquesEventsItem `json:"items"`
		NextLink   string                               `json:"nextLink"`
	}

	ObservedAttackTechniquesEventsItem struct {
		Source  string `json:"source"`
		UUID    string `json:"uuid"`
		Filters []struct {
			ID                 string   `json:"id"`
			Name               string   `json:"name"`
			Description        string   `json:"description"`
			MitreTacticIds     []string `json:"mitreTacticIds"`
			MitreTechniqueIds  []string `json:"mitreTechniqueIds"`
			HighlightedObjects []struct {
				Type  string `json:"type"`
				Field string `json:"field"`
				Value any    `json:"value"`
			} `json:"highlightedObjects"`
			RiskLevel string `json:"riskLevel"`
			Type      string `json:"type"`
		} `json:"filters"`
		Endpoint struct {
			EndpointName string   `json:"endpointName"`
			AgentGUID    string   `json:"agentGuid"`
			Ips          []string `json:"ips"`
		} `json:"endpoint"`
		EntityType       string   `json:"entityType"`
		EntityName       any      `json:"entityName"` // Can be list, so not string
		DetectedDateTime VOneTime `json:"detectedDateTime"`
		IngestedDateTime VOneTime `json:"ingestedDateTime"`
		Detail           struct {
			EventTime               string   `json:"eventTime"`
			Tags                    []string `json:"tags"`
			UUID                    string   `json:"uuid"`
			ProductCode             string   `json:"productCode"`
			FilterRiskLevel         string   `json:"filterRiskLevel"`
			BitwiseFilterRiskLevel  int      `json:"bitwiseFilterRiskLevel"`
			EventID                 string   `json:"eventId"`
			EventSubID              int      `json:"eventSubId"`
			EventHashID             string   `json:"eventHashId"`
			FirstSeen               string   `json:"firstSeen"`
			LastSeen                string   `json:"lastSeen"`
			EndpointGUID            string   `json:"endpointGuid"`
			EndpointHostName        string   `json:"endpointHostName"`
			EndpointIP              []string `json:"endpointIp"`
			EndpointMacAddress      []string `json:"endpointMacAddress"`
			Timezone                string   `json:"timezone"`
			Pname                   string   `json:"pname"`
			Pver                    string   `json:"pver"`
			Plang                   int      `json:"plang"`
			Pplat                   int      `json:"pplat"`
			OsName                  string   `json:"osName"`
			OsVer                   string   `json:"osVer"`
			OsDescription           string   `json:"osDescription"`
			OsType                  string   `json:"osType"`
			ProcessHashID           string   `json:"processHashId"`
			ProcessName             string   `json:"processName"`
			ProcessPid              int      `json:"processPid"`
			SessionID               int      `json:"sessionId"`
			ProcessUser             string   `json:"processUser"`
			ProcessUserDomain       string   `json:"processUserDomain"`
			ProcessLaunchTime       string   `json:"processLaunchTime"`
			ProcessCmd              string   `json:"processCmd"`
			AuthID                  string   `json:"authId"`
			IntegrityLevel          int      `json:"integrityLevel"`
			ProcessFileHashID       string   `json:"processFileHashId"`
			ProcessFilePath         string   `json:"processFilePath"`
			ProcessFileHashSha1     string   `json:"processFileHashSha1"`
			ProcessFileHashSha256   string   `json:"processFileHashSha256"`
			ProcessFileHashMd5      string   `json:"processFileHashMd5"`
			ProcessSigner           []string `json:"processSigner"`
			ProcessSignerValid      []bool   `json:"processSignerValid"`
			ProcessFileSize         string   `json:"processFileSize"`
			ProcessFileCreation     string   `json:"processFileCreation"`
			ProcessFileModifiedTime string   `json:"processFileModifiedTime"`
			ProcessTrueType         int      `json:"processTrueType"`
			ObjectHashID            string   `json:"objectHashId"`
			ObjectUser              string   `json:"objectUser"`
			ObjectUserDomain        string   `json:"objectUserDomain"`
			ObjectSessionID         string   `json:"objectSessionId"`
			ObjectFilePath          string   `json:"objectFilePath"`
			ObjectFileHashSha1      string   `json:"objectFileHashSha1"`
			ObjectFileHashSha256    string   `json:"objectFileHashSha256"`
			ObjectFileHashMd5       string   `json:"objectFileHashMd5"`
			ObjectSigner            []string `json:"objectSigner"`
			ObjectSignerValid       []bool   `json:"objectSignerValid"`
			ObjectFileSize          string   `json:"objectFileSize"`
			ObjectFileCreation      string   `json:"objectFileCreation"`
			ObjectFileModifiedTime  string   `json:"objectFileModifiedTime"`
			ObjectTrueType          int      `json:"objectTrueType"`
			ObjectName              string   `json:"objectName"`
			ObjectPid               int      `json:"objectPid"`
			ObjectLaunchTime        string   `json:"objectLaunchTime"`
			ObjectCmd               string   `json:"objectCmd"`
			ObjectAuthID            string   `json:"objectAuthId"`
			ObjectIntegrityLevel    int      `json:"objectIntegrityLevel"`
			ObjectFileHashID        string   `json:"objectFileHashId"`
			ObjectRunAsLocalAccount bool     `json:"objectRunAsLocalAccount"`
		} `json:"detail"`
	}
)

type GetOATEventsFunc struct {
	baseFunc
	Response ObservedAttackTechniquesEventsResponse
	top      int
}

var _ vOneFunc = &GetOATEventsFunc{}

func (v *VOne) GetOATEvents() *GetOATEventsFunc {
	f := &GetOATEventsFunc{}
	f.baseFunc.init(v)
	return f
}

func (f *GetOATEventsFunc) DetectedStart(vOneTime VOneTime) *GetOATEventsFunc {
	f.setParameter("detectedStartDateTime", vOneTime.String()+"Z")
	return f
}

func (f *GetOATEventsFunc) DetectedEnd(vOneTime VOneTime) *GetOATEventsFunc {
	f.setParameter("detectedEndDateTime", vOneTime.String()+"Z")
	return f
}

func (f *GetOATEventsFunc) IngestedStart(vOneTime VOneTime) *GetOATEventsFunc {
	f.setParameter("ingestedStartDateTime", vOneTime.String()+"Z")
	return f
}

func (f *GetOATEventsFunc) IngestedEnd(vOneTime VOneTime) *GetOATEventsFunc {
	f.setParameter("ingestedEndDateTime", vOneTime.String()+"Z")
	return f
}

// Filter - filter events
func (f *GetOATEventsFunc) Filter(filter string) *GetOATEventsFunc {
	f.setHeader("TMV1-Filter", filter)
	f.Response.NextLink = ""
	return f
}

// Top - set limit for returned amount of items
func (f *GetOATEventsFunc) Top(t Top) *GetOATEventsFunc {
	f.setParameter("top", t.String())
	f.top = t.Int()
	return f
}

// Iterate - get all events matching query one by one. If callback returns
// non nil error, iteration is aborted and this error is returned
func (f *GetOATEventsFunc) Iterate(ctx context.Context,
	callback func(item *ObservedAttackTechniquesEventsItem) error) error {
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
		if response.Count != f.top {
			break
		}
	}
	return nil
}

// Range - iterator for all endpoints matching query (go 1.23 and later)
func (f *GetOATEventsFunc) Range(ctx context.Context) iter.Seq2[*ObservedAttackTechniquesEventsItem, error] {
	return func(yield func(*ObservedAttackTechniquesEventsItem, error) bool) {
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

// Do - execute the API call
func (f *GetOATEventsFunc) Do(ctx context.Context) (*ObservedAttackTechniquesEventsResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *GetOATEventsFunc) method() string {
	return methodGet
}

func (s *GetOATEventsFunc) url() string {
	return "/v3.0/oat/detections"
}

func (f *GetOATEventsFunc) responseStruct() any {
	return &f.Response
}
