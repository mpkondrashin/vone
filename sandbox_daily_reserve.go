/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_daily_reserve.go - get dalily usage quota
*/
package vone

/*
{
   "submissionReserveCount":50,
   "submissionRemainingCount":19,
   "submissionCount":31,
   "submissionExemptionCount":5,
   "submissionCountDetail":{
      "fileCount":27,
      "fileExemptionCount":5,
      "urlCount":4,
      "urlExemptionCount":0
   }
}
*/
type SandboxDailyReserveResponse struct {
	SubmissionReserveCount   int `json:"submissionReserveCount"`
	SubmissionRemainingCount int `json:"submissionRemainingCount"`
	SubmissionCount          int `json:"submissionCount"`
	SubmissionExemptionCount int `json:"submissionExemptionCount"`
	SubmissionCountDetail    struct {
		FileCount          int `json:"fileCount"`
		FileExemptionCount int `json:"fileExemptionCount"`
		URLCount           int `json:"urlCount"`
		URLExemptionCount  int `json:"urlExemptionCount"`
	} `json:"submissionCountDetail"`
}

type SanboxDailyReserveFunc struct {
	BaseFunc
	Response SandboxDailyReserveResponse
}

func (f *SanboxDailyReserveFunc) Do() (*SandboxDailyReserveResponse, error) {
	if err := f.vone.Call(f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (v *VOne) SandboxDailyReserve() *SanboxDailyReserveFunc {
	f := &SanboxDailyReserveFunc{}
	f.BaseFunc.Init(v)
	return f
}

func (s *SanboxDailyReserveFunc) URL() string {
	return "/v3.0/sandbox/submissionUsage"
}

func (f *SanboxDailyReserveFunc) ResponseStruct() any {
	return &f.Response
}
