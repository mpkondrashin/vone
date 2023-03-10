// Code generated by enum (github.com/mpkondrashin/enum). DO NOT EDIT

package vone

import (
    "encoding/json"
    "fmt"
    "strconv"
)

type Action int

const (
    ActionAnalyzeFile Action = iota
    ActionAnalyzeUrl Action = iota
)

func (v Action)String() string {
    s, ok := map[Action]string {
        ActionAnalyzeFile: "analyzeFile",
        ActionAnalyzeUrl: "analyzeUrl",
    }[v]
    if ok {
        return s
    }
    return "Action(" + strconv.FormatInt(int64(v), 10) + ")"
}

func (s *Action) UnmarshalJSON(data []byte) error {
    var v string
    if err := json.Unmarshal(data, &v); err != nil {
        return err
    }
    result, ok := map[string]Action{
        "analyzeFile": ActionAnalyzeFile,
        "analyzeUrl": ActionAnalyzeUrl,
    }[v]
    if !ok {
        return fmt.Errorf("%w: %s", ErrEnumUnknown, v)
    }
    *s = result
    return nil
}
