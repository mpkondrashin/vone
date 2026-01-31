/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	sandbox_investigation_package.go - download investigation package
*/

package vone

import (
	"fmt"
)

type sandboxInvestigationPackageRequest struct {
	*sandboxDownloadResultsRequest
}

func (v *VOne) SandboxInvestigationPackage(id string) *sandboxInvestigationPackageRequest {
	return &sandboxInvestigationPackageRequest{v.SandboxDownloadResults(id)}
}

func (s *sandboxInvestigationPackageRequest) url() string {
	return fmt.Sprintf("/v3.0/sandbox/analysisResults/%s/investigationPackage", s.id)
}
