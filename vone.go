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
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	methodGet  = "GET"
	methodPost = "POST"

	applicationJSON = "application/json"

	//timeFormat = "2006-01-02T15:04:05Z"
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
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d (%v): %s",
		e.ErrorData.Code,
		e.ErrorData.Code,
		e.ErrorData.Message)
	if e.ErrorData.Innererror.Code != "" {
		fmt.Fprintf(&sb, " (%s: %s)",
			e.ErrorData.Innererror.Service,
			e.ErrorData.Innererror.Code,
		)
	}
	return sb.String()
}

type VisionOneTime time.Time

const (
	timeFormat  = `2006-01-02T15:04:05`
	timeFormatZ = layout + "Z"
)

var _ json.Unmarshaler = (*VisionOneTime)(nil)

// Implement Marshaler interface
func (vot *VisionOneTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if len(s) == 0 {
		*vot = VisionOneTime(time.Time{})
		return nil
	}
	//log.Println("UnmarshalJSON s ", s)
	t, err := time.Parse(timeFormat, s)
	if err != nil {
		t, err = time.Parse(timeFormatZ, s)
		if err != nil {
			return err
		}
	}
	//log.Println("UnmarshalJSON  t ", t)
	*vot = VisionOneTime(t)
	return nil
}

// Implement Unmarshaler interface
func (vot VisionOneTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + vot.String() + "\""), nil
}

func (vot VisionOneTime) String() string {
	return time.Time(vot).Format(timeFormatZ)
}

// Convert the internal date as CSV string
func (vot VisionOneTime) MarshalCSV() (string, error) {
	return vot.String(), nil
}

// VOne - Vision One API struct
type VOne struct {
	Domain            string
	Token             string
	transportModifier func(*http.Transport)
	rateLimiter       RateLimiter
}

// NewVOne - create VOne struct
func NewVOne(domain string, token string) *VOne {
	return &VOne{
		Domain: domain,
		Token:  token,
	}
}

func (v *VOne) SetRateLimiter(rateLimiter RateLimiter) *VOne {
	v.rateLimiter = rateLimiter
	return v
}

func (v *VOne) AddTransportModifier(transportModifier func(*http.Transport)) {
	AddTransportModifier(&v.transportModifier, transportModifier)
}

const HTTPResponseTooManyRequests = 429

var VOneRateLimitSurpassedError RateLimitSurpassed = func(err error) bool {
	var vOneErr *Error
	if !errors.As(err, &vOneErr) {
		return false
	}
	return vOneErr.ErrorData.Code == HTTPResponseTooManyRequests
}

func (v *VOne) callWithoutLimiter(ctx context.Context, f vOneFunc) error {
	uri := f.uri()
	if uri == "" {
		uri = "https://" + v.Domain + f.url()
	}
	return v.callURL(ctx, f, uri)
}

func (v *VOne) callWithLimiter(ctx context.Context, f vOneFunc) error {
	for {
		if v.rateLimiter.ShouldAbort() {
			//	log.Println("return ErrStop")
			return ErrStop
		}
		//log.Println("no should abort")
		err := v.callWithoutLimiter(ctx, f)
		//log.Println("err", err)
		err = v.rateLimiter.CheckError(err)
		//log.Println("err after CheckError", err)
		if err == ErrOnceMore {
			//log.Println("ErrOnceMore")
			continue
		}
		//log.Println("return", err)
		return err
	}
}

/*
	func AddLimiter(f func(context.Context, vOneFunc) error, limiter RateLimiter) func(context.Context, vOneFunc) error {
		return func(context.Context, vOneFunc) error {
			for {
				if err := limiter.Sleep(); err != nil {
					return err
				}
				err := f(ctx, f)
				if err.Error() == "419" {

				}
			}
		}
	}
*/
func (v *VOne) call(ctx context.Context, f vOneFunc) error {
	if v.rateLimiter != nil {
		return v.callWithLimiter(ctx, f)
	}
	return v.callWithoutLimiter(ctx, f)
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
	if v.transportModifier != nil {
		transport := &http.Transport{}
		v.transportModifier(transport)
		client.Transport = transport
	}
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
		//vOneErr.ErrorData.Message += strconv.Itoa(resp.StatusCode)
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

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read body : %w", err)
	}
	w, _ := os.Create("body.json")
	w.Write(bodyBytes)
	w.Close()
	r := f.responseStruct()
	err = json.Unmarshal(bodyBytes, r)
	if err != nil {
		return fmt.Errorf("response parse error: %w [%s]", err, string(bodyBytes))
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
