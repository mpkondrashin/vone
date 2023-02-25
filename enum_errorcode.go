// Code generated by enum (github.com/mpkondrashin/enum). DO NOT EDIT

package vone

import (
    "encoding/json"
    "fmt"
    "strconv"
)

type ErrorCode int

const (
    ErrorCodeAccessDenied ErrorCode = iota
    ErrorCodeBadRequest ErrorCode = iota
    ErrorCodeConditionNotMet ErrorCode = iota
    ErrorCodeInternalServerError ErrorCode = iota
    ErrorCodeNotFound ErrorCode = iota
    ErrorCodeParameterNotAccepted ErrorCode = iota
    ErrorCodeRequestEntityTooLarge ErrorCode = iota
    ErrorCodeTooManyRequests ErrorCode = iota
    ErrorCodeUnsupported ErrorCode = iota
)

func (v ErrorCode)String() string {
    s, ok := map[ErrorCode]string {
        ErrorCodeAccessDenied: "AccessDenied",
        ErrorCodeBadRequest: "BadRequest",
        ErrorCodeConditionNotMet: "ConditionNotMet",
        ErrorCodeInternalServerError: "InternalServerError",
        ErrorCodeNotFound: "NotFound",
        ErrorCodeParameterNotAccepted: "ParameterNotAccepted",
        ErrorCodeRequestEntityTooLarge: "RequestEntityTooLarge",
        ErrorCodeTooManyRequests: "TooManyRequests",
        ErrorCodeUnsupported: "Unsupported",
    }[v]
    if ok {
        return s
    }
    return "ErrorCode(" + strconv.FormatInt(int64(v), 10) + ")"
}

func (s *ErrorCode) UnmarshalJSON(data []byte) error {
    var v string
    if err := json.Unmarshal(data, &v); err != nil {
        return err
    }
    result, ok := map[string]ErrorCode{
        "AccessDenied": ErrorCodeAccessDenied,
        "BadRequest": ErrorCodeBadRequest,
        "ConditionNotMet": ErrorCodeConditionNotMet,
        "InternalServerError": ErrorCodeInternalServerError,
        "NotFound": ErrorCodeNotFound,
        "ParameterNotAccepted": ErrorCodeParameterNotAccepted,
        "RequestEntityTooLarge": ErrorCodeRequestEntityTooLarge,
        "TooManyRequests": ErrorCodeTooManyRequests,
        "Unsupported": ErrorCodeUnsupported,
    }[v]
    if !ok {
        return fmt.Errorf("%w: %s", ErrEnumUnknown, v)
    }
    *s = result
    return nil
}
