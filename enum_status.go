// Code generated by enum (github.com/mpkondrashin/enum). DO NOT EDIT

package vone

import (
    "encoding/json"
	"errors"
    "fmt"
    "strconv"
	"strings"
)

type Status int

const (
    StatusSucceeded Status = iota
    StatusRunning Status = iota
    StatusFailed Status = iota
)

// String - return string representation for Status value
func (v Status)String() string {
    s, ok := map[Status]string {
        StatusSucceeded: "succeeded",
        StatusRunning: "running",
        StatusFailed: "failed",
    }[v]
    if ok {
        return s
    }
    return "Status(" + strconv.FormatInt(int64(v), 10) + ")"
}

// ErrUnknownStatus - will be returned wrapped when parsing string
// containing unrecognized value.
var ErrUnknownStatus = errors.New("unknown Status")

var mapStatusFromString = map[string]Status{
    "succeeded": StatusSucceeded,
    "running": StatusRunning,
    "failed": StatusFailed,
}

// UnmarshalJSON implements the Unmarshaler interface of the json package for Status.
func (s *Status) UnmarshalJSON(data []byte) error {
    var v string
    if err := json.Unmarshal(data, &v); err != nil {
        return err
    }
    result, ok := mapStatusFromString[strings.ToLower(v)]
    if !ok {
        return fmt.Errorf("%w: %s", ErrUnknownStatus, v)
    }
    *s = result
    return nil
}

// UnmarshalYAML implements the Unmarshaler interface of the yaml.v3 package for Status.
func (s *Status) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	if err := unmarshal(&v); err != nil {
		return err
	}
	result, ok := mapStatusFromString[strings.ToLower(v)]		
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownStatus, v)
	}
	*s = result
	return nil
}
