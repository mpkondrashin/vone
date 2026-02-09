/*
	Trend Micro Vision One API SDK
	(c) 2026 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Workbench API capabilities

	workbench_list.go - list workbench alerts
*/

package vone

import (
	"context"
	"io"
	"time"
)

type (
	// EntityValue - entity value structure (can be string or object)
	EntityValue struct {
		GUID string   `json:"guid,omitempty"`
		Name string   `json:"name,omitempty"`
		IPs  []string `json:"ips,omitempty"`
	}

	// Entity - impact scope entity
	Entity struct {
		EntityType                    string      `json:"entityType"`
		EntityValue                   interface{} `json:"entityValue"` // can be string or EntityValue object
		EntityID                      string      `json:"entityId"`
		ManagementScopeGroupID        string      `json:"managementScopeGroupId,omitempty"`
		ManagementScopeInstanceID     string      `json:"managementScopeInstanceId,omitempty"`
		ManagementScopePartitionKey   string      `json:"managementScopePartitionKey,omitempty"`
		RelatedEntities               []string    `json:"relatedEntities"`
		RelatedIndicatorIDs           []int       `json:"relatedIndicatorIds"`
		Provenance                    []string    `json:"provenance"`
	}

	// ImpactScope - impact scope of the alert
	ImpactScope struct {
		DesktopCount        int      `json:"desktopCount"`
		ServerCount         int      `json:"serverCount"`
		AccountCount        int      `json:"accountCount"`
		EmailAddressCount   int      `json:"emailAddressCount"`
		ContainerCount      int      `json:"containerCount"`
		CloudIdentityCount  int      `json:"cloudIdentityCount"`
		Entities            []Entity `json:"entities"`
	}

	// MatchedEvent - matched event in a filter
	MatchedEvent struct {
		UUID            string        `json:"uuid"`
		MatchedDateTime VisionOneTime `json:"matchedDateTime"`
		Type            string        `json:"type"`
	}

	// MatchedFilter - matched filter in a rule
	MatchedFilter struct {
		ID                string         `json:"id"`
		Name              string         `json:"name"`
		MatchedDateTime   VisionOneTime  `json:"matchedDateTime"`
		MitreTechniqueIDs []string       `json:"mitreTechniqueIds"`
		MatchedEvents     []MatchedEvent `json:"matchedEvents"`
	}

	// MatchedRule - matched rule
	MatchedRule struct {
		ID             string          `json:"id"`
		Name           string          `json:"name"`
		MatchedFilters []MatchedFilter `json:"matchedFilters"`
	}

	// Indicator - alert indicator
	Indicator struct {
		ID              int      `json:"id"`
		Type            string   `json:"type"`
		Field           string   `json:"field"`
		Value           string   `json:"value"`
		RelatedEntities []string `json:"relatedEntities"`
		FilterIDs       []string `json:"filterIds"`
		Provenance      []string `json:"provenance"`
	}

	// WorkbenchAlert - structure of a single workbench alert
	WorkbenchAlert struct {
		SchemaVersion              string        `json:"schemaVersion"`
		ID                         string        `json:"id"`
		InvestigationStatus        string        `json:"investigationStatus"`
		Status                     string        `json:"status"`
		InvestigationResult        string        `json:"investigationResult"`
		WorkbenchLink              string        `json:"workbenchLink"`
		AlertProvider              string        `json:"alertProvider"`
		ModelID                    string        `json:"modelId"`
		Model                      string        `json:"model"`
		ModelType                  string        `json:"modelType"`
		Score                      int           `json:"score"`
		Severity                   string        `json:"severity"`
		FirstInvestigatedDateTime  VisionOneTime `json:"firstInvestigatedDateTime,omitempty"`
		CreatedDateTime            VisionOneTime `json:"createdDateTime"`
		UpdatedDateTime            VisionOneTime `json:"updatedDateTime"`
		IncidentID                 string        `json:"incidentId,omitempty"`
		CaseID                     string        `json:"caseId,omitempty"`
		OwnerIDs                   []string      `json:"ownerIds,omitempty"`
		ImpactScope                ImpactScope   `json:"impactScope"`
		Description                string        `json:"description"`
		MatchedRules               []MatchedRule `json:"matchedRules"`
		Indicators                 []Indicator   `json:"indicators"`
	}

	// WorkbenchAlertsResponse - response structure for workbench alerts list
	WorkbenchAlertsResponse struct {
		TotalCount int              `json:"totalCount"`
		Count      int              `json:"count"`
		Items      []WorkbenchAlert `json:"items"`
		NextLink   string           `json:"nextLink"`
	}
)

type workbenchListRequest struct {
	baseRequest
	response WorkbenchAlertsResponse
}

// WorkbenchListAlerts - create a new request to list workbench alerts
func (v *VOne) WorkbenchListAlerts() *workbenchListRequest {
	f := &workbenchListRequest{}
	f.baseRequest.init(v)
	return f
}

// StartDateTime - set start date/time filter
func (f *workbenchListRequest) StartDateTime(t time.Time) *workbenchListRequest {
	f.setParameter("startDateTime", t.Format(timeFormatZ))
	return f
}

// EndDateTime - set end date/time filter
func (f *workbenchListRequest) EndDateTime(t time.Time) *workbenchListRequest {
	f.setParameter("endDateTime", t.Format(timeFormatZ))
	return f
}

// DateTimeTarget - set which datetime field to filter on
func (f *workbenchListRequest) DateTimeTarget(target string) *workbenchListRequest {
	f.setParameter("dateTimeTarget", target)
	return f
}

// OrderBy - set ordering for results
func (f *workbenchListRequest) OrderBy(orderBy string) *workbenchListRequest {
	f.setParameter("orderBy", orderBy)
	return f
}

// Filter - set TMV1-Filter header for filtering results
func (f *workbenchListRequest) Filter(filter string) *workbenchListRequest {
	f.setHeader("TMV1-Filter", filter)
	return f
}

// Do - execute the request and return workbench alerts
func (f *workbenchListRequest) Do(ctx context.Context) (*WorkbenchAlertsResponse, error) {
	if f.vone.mockup != nil {
		// Add mockup support if needed
		return nil, nil
	}
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.response, nil
}

func (*workbenchListRequest) url() string {
	return "/v3.0/workbench/alerts"
}

func (f *workbenchListRequest) uri() string {
	return f.response.NextLink
}

func (f *workbenchListRequest) responseStruct() any {
	return &f.response
}

func (f *workbenchListRequest) nextLink() string {
	return f.response.NextLink
}

func (f *workbenchListRequest) resetPagination() {
	f.response.NextLink = ""
}

// Next - get next page of results
func (f *workbenchListRequest) Next(ctx context.Context) (*WorkbenchAlertsResponse, error) {
	if f.response.NextLink == "" {
		return nil, io.EOF
	}
	return f.Do(ctx)
}

// Paginator - create a paginator for iterating through all results
func (f *workbenchListRequest) Paginator() *Paginator[
	WorkbenchAlertsResponse,
	WorkbenchAlert,
] {
	return NewPaginator(
		f,
		func(r *WorkbenchAlertsResponse) []WorkbenchAlert {
			return r.Items
		},
	)
}
