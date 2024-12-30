package vone

import (
	"fmt"
	"strings"
)

type Field int

const (
	FieldAgentGuid Field = iota
	FieldLoginAccount
	FieldEndpointName
	FieldMACAddress
	FieldIP
	FieldProtectionManager
	FieldPolicyName
	FieldComponentUpdatePolicy
	FieldComponentUpdateStatus
	FieldComponentVersion
	FieldOSName
	FieldOSVersion
	FieldProductCode
	FieldInstalledProductCodes
)

//go:generate stringer -type Field -trimprefix Field

type ProductCode int

const (
	ProductCodeSAO ProductCode = iota
	ProductCodeSDS
	ProductCodeXES
)

//go:generate stringer -type ProductCode -trimprefix ProductCode

type OSName int

const (
	OSNameLinux OSName = iota
	OSNameWindows
	OSNameMacOS
	OSNameMacOSX
)

//go:generate stringer -type OSName -trimprefix OSName

type ComponentVersion int

const (
	ComponentVersionOutdatedVersion ComponentVersion = iota
	ComponentVersionLatestVersion
	ComponentVersionControlledLatestVersion
	ComponentVersionUnknownVersions
)

//go:generate stringer -type ComponentVersion -trimprefix ComponentVersion

type ComponentUpdateStatus int

const (
	ComponentUpdateStatusPause ComponentUpdateStatus = iota
	ComponentUpdateStatusOnSchedule
	ComponentUpdateStatusNotSupported
)

//go:generate stringer -type ComponentUpdateStatus -trimprefix ComponentUpdateStatus

// Filter represents a condition or a group of conditions
type Filter struct {
	Field    Field
	Operator string
	Value    string
	Children []Filter
	Logic    string
	Negate   bool
}

// Build generates the filter string based on the structure
func (f Filter) Build() string {
	if len(f.Children) > 0 {
		var parts []string
		for _, child := range f.Children {
			parts = append(parts, child.Build())
		}
		combined := strings.Join(parts, fmt.Sprintf(" %s ", strings.ToLower(f.Logic)))
		if f.Negate {
			return fmt.Sprintf("not (%s)", combined)
		}
		return fmt.Sprintf("(%s)", combined)
	}
	condition := fmt.Sprintf("%v %s '%s'", f.Field, strings.ToLower(f.Operator), f.Value)
	if f.Negate {
		return fmt.Sprintf("not (%s)", condition)
	}
	return condition
}

// Helper functions for creating filters
//func NewCondition(field Field, operator, value string) Filter {
//	return Filter{Field: field, Operator: operator, Value: value}
//}

func AgentGuidEq(value string) Filter {
	return Filter{
		Field:    FieldLoginAccount,
		Operator: "eq",
		Value:    value,
	}
}

func EndpointNameEq(value string) Filter {
	return Filter{
		Field:    FieldEndpointName,
		Operator: "eq",
		Value:    value,
	}
}

func LoginAccountEq(value string) Filter {
	return Filter{
		Field:    FieldLoginAccount,
		Operator: "eq",
		Value:    value,
	}
}

func MACAddressEq(value string) Filter {
	return Filter{
		Field:    FieldMACAddress,
		Operator: "eq",
		Value:    value,
	}
}

func IPEq(value string) Filter {
	return Filter{
		Field:    FieldIP,
		Operator: "eq",
		Value:    value,
	}
}

func ProtectionManagerEq(value string) Filter {
	return Filter{
		Field:    FieldProtectionManager,
		Operator: "eq",
		Value:    value,
	}
}

func PolicyNameEq(value string) Filter {
	return Filter{
		Field:    FieldPolicyName,
		Operator: "eq",
		Value:    value,
	}
}

func ComponentUpdatePolicyEq(value string) Filter {
	return Filter{
		Field:    FieldComponentUpdatePolicy,
		Operator: "eq",
		Value:    value,
	}
}

func ComponentUpdateStatusEq(value string) Filter {
	return Filter{
		Field:    FieldComponentUpdateStatus,
		Operator: "eq",
		Value:    value,
	}
}

func ComponentVersionEq(value string) Filter {
	return Filter{
		Field:    FieldComponentVersion,
		Operator: "eq",
		Value:    value,
	}
}

func OSNameEq(value OSName) Filter {
	return Filter{
		Field:    FieldOSName,
		Operator: "eq",
		Value:    value.String(),
	}
}

func OSVersionEq(value OSName) Filter {
	return Filter{
		Field:    FieldOSVersion,
		Operator: "eq",
		Value:    value.String(),
	}
}

func ProductCodeEq(value ProductCode) Filter {
	return Filter{
		Field:    FieldProductCode,
		Operator: "eq",
		Value:    value.String(),
	}
}

func InstalledProductCodesEq(value ProductCode) Filter {
	return Filter{
		Field:    FieldInstalledProductCodes,
		Operator: "eq",
		Value:    value.String(),
	}
}

func And(filters ...Filter) Filter {
	return Filter{Children: filters, Logic: "and"}
}

func Or(filters ...Filter) Filter {
	return Filter{Children: filters, Logic: "or"}
}

func Not(filter Filter) Filter {
	filter.Negate = true
	return filter
}
