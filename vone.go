/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	VOne - Web API SDK

	vone.go - main SDK declarations
*/

package vone

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

const (
	methodGet  = "GET"
	methodPost = "POST"

	applicationJSON = "application/json"

	timeFormat = "2006-1-02T15:04:05Z"
)

type Error struct {
	ErrorData struct {
		Message    string    `json:"message"`
		Code       ErrorCode `json:"code"`
		Innererror struct {
			Service string `json:"service"`
			Code    string `json:"code"`
		} `json:"innererror"`
	} `json:"error"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v: %s", e.ErrorData.Code, e.ErrorData.Message)
}

type VOne struct {
	Domain string
	Token  string
}

func NewVOne(domain string, token string) *VOne {
	return &VOne{
		Domain: domain,
		Token:  token,
	}
}

func (v *VOne) call(ctx context.Context, f vOneFunc) error {
	uri := f.uri()
	if uri == "" {
		uri = "https://" + v.Domain + f.url()
	}
	return v.callURL(ctx, f, uri)
}

func (v *VOne) callURL(ctx context.Context, f vOneFunc, uri string) error {
	req, err := http.NewRequestWithContext(ctx, f.method(), uri, f.requestBody())
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+v.Token)
	if f.requestBody() != nil {
		req.Header.Set("Content-Type", f.contentType())
	}
	f.populate(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request: %w", err)
	}
	if GetHTTPCodeRange(resp.StatusCode) != HTTPCodeSuccessRange {
		var data bytes.Buffer
		if _, err := io.Copy(&data, resp.Body); err != nil {
			return fmt.Errorf("download error: %w", err)
		}
		vOneErr := new(Error)
		if err := json.Unmarshal(data.Bytes(), vOneErr); err != nil {
			return fmt.Errorf("parse error: %w", err)
		}
		return fmt.Errorf("request error: %w", vOneErr)
	}

	if err := v.PopulateResponseStruct(f.responseHeader(), resp.Header); err != nil {
		return err
	}

	if f.responseStruct() == nil {
		f.responseBody(resp.Body)
		return nil
	}

	return v.DecodeBody(f, resp.Body)
}

func (v *VOne) DecodeBody(f vOneFunc, body io.ReadCloser) error {
	defer body.Close()
	err := json.NewDecoder(body).Decode(f.responseStruct())
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("response parse error: %w", err)
	}
	return nil
}

var ErrUnsupportedType = errors.New("unsupported type")

func (v *VOne) PopulateResponseStruct(structPtr any, header http.Header) error {
	if structPtr == nil {
		return nil
	}
	structPtrValue := reflect.ValueOf(structPtr)
	structValue := reflect.Indirect(structPtrValue)
	structValueType := structValue.Type()
	for i := 0; i < structValueType.NumField(); i++ {
		fieldValue := structValue.Field(i)
		headerName := structValueType.Field(i).Tag.Get("header")
		headerValue := header.Get(headerName)
		kind := fieldValue.Kind()
		switch kind {
		case reflect.String:
			fieldValue.SetString(headerValue)
		case reflect.Int:
			x, err := strconv.Atoi(headerValue)
			if err != nil {
				return err
			}
			fieldValue.SetInt(int64(x))
		default:
			return fmt.Errorf("%s: %w", kind, ErrUnsupportedType)
		}
	}
	return nil
}

type HTTPCodeRange int

const (
	HTTPCodeInformationalRange HTTPCodeRange = iota + 1
	HTTPCodeSuccessRange
	HTTPCodeRedirectRange
	HTTPCodeClientErrorRange
	HTTPCodeServerErrorRange
)

func GetHTTPCodeRange(code int) HTTPCodeRange {
	return HTTPCodeRange(code / 100)
}
