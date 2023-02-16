package vone

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

var _ Func = &SanboxDailyReserveFunc{}

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
