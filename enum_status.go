// Code generated by enum (github.com/mpkondrashin/enum). DO NOT EDIT

package vone

import (
    "encoding/json"
    "fmt"
    "strconv"
)

type Status int

const (
    StatusSucceeded Status = iota
    StatusRunning Status = iota
    StatusFailed Status = iota
)

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

func (s *Status) UnmarshalJSON(data []byte) error {
    var v string
    if err := json.Unmarshal(data, &v); err != nil {
        return err
    }
    result, ok := map[string]Status{
        "succeeded": StatusSucceeded,
        "running": StatusRunning,
        "failed": StatusFailed,
    }[v]
    if !ok {
        return fmt.Errorf("%w: %s", ErrEnumUnknown, v)
    }
    *s = result
    return nil
}
