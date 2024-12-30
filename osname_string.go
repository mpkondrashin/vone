// Code generated by "stringer -type OSName -trimprefix OSName"; DO NOT EDIT.

package vone

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OSNameLinux-0]
	_ = x[OSNameWindows-1]
	_ = x[OSNameMacOS-2]
	_ = x[OSNameMacOSX-3]
}

const _OSName_name = "LinuxWindowsMacOSMacOSX"

var _OSName_index = [...]uint8{0, 5, 12, 17, 23}

func (i OSName) String() string {
	if i < 0 || i >= OSName(len(_OSName_index)-1) {
		return "OSName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OSName_name[_OSName_index[i]:_OSName_index[i+1]]
}