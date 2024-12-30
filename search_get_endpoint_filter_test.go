package vone

import (
	"testing"
)

func TestFilter(t *testing.T) {

	// Example 1: (endpointName eq 'sample-host') or (macAddress eq '00:11:22:33:44:55')
	filter1 := Or(
		EndpointNameEq("sample-host"),
		MACAddressEq("00:11:22:33:44:55"),
	)

	// Example 2: not (osName eq 'Windows')
	filter2 := Not(OSNameEq(OSNameWindows))

	// Combine filters if needed
	combinedFilter := And(filter1, filter2)

	// Print filters
	t.Log("Filter 1:", filter1.Build())
	t.Log("Filter 2:", filter2.Build())
	t.Log("Combined Filter:", combinedFilter.Build())
}
