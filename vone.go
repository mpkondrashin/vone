/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	VOne - Web API SDK

	vone.go - main SDK declarations
*/

package vone

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

type (
	Innererror struct {
		Service string `json:"service"`
		Code    string `json:"code"`
	}

	Error struct {
		Message    string     `json:"message"`
		Code       ErrorCode  `json:"code"`
		Innererror Innererror `json:"innererror"`
	}

	ErrorData struct {
		Error Error `json:"error"`
	}
)

func ErrorFromReader(reader io.Reader) (Error, error) {
	dec := json.NewDecoder(reader)
	var data ErrorData
	err := dec.Decode(&data)
	if err != nil {
		return Error{}, err
	}
	return data.Error, nil
}

func (e Error) Error() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d (%v): %s",
		e.Code,
		e.Code,
		e.Message)
	if e.Innererror.Code != "" {
		fmt.Fprintf(&sb, " (%s: %s)",
			e.Innererror.Service,
			e.Innererror.Code,
		)
	}
	return sb.String()
}

type HTTPError struct {
	Status int
	Err    error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http %d: %v", e.Status, e.Err)
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

type VisionOneTime time.Time

const (
	timeFormat  = `2006-01-02T15:04:05`
	timeFormatZ = timeFormat + "Z"
)

var _ json.Unmarshaler = (*VisionOneTime)(nil)

// Implement Marshaler interface
func (vot *VisionOneTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*vot = VisionOneTime(time.Time{})
		return nil
	}

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "" {
		*vot = VisionOneTime(time.Time{})
		return nil
	}

	t, err := time.Parse(timeFormat, s)
	if err != nil {
		t, err = time.Parse(timeFormatZ, s)
		if err != nil {
			return err
		}
	}

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
	Domain      string
	Token       string
	client      *http.Client
	rateLimiter RateLimiter
	mockup      *SandboxMockup
}

//transportModifier func(*http.Transport)

// NewVOne - create VOne struct
func NewVOne(domain string, token string) *VOne {
	return &VOne{
		Domain: domain,
		Token:  token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (v *VOne) SetRateLimiter(rateLimiter RateLimiter) *VOne {
	v.rateLimiter = rateLimiter
	return v
}

func (v *VOne) AddTransportModifier(transportModifier func(*http.Transport)) {
	tr, ok := v.client.Transport.(*http.Transport)
	if !ok || tr == nil {
		tr = http.DefaultTransport.(*http.Transport).Clone()
	}
	transportModifier(tr)
	v.client.Transport = tr
}

func (v *VOne) SetMockup(mockup *SandboxMockup) *VOne {
	v.mockup = mockup
	return v
}

const HTTPResponseTooManyRequests = 429

var VOneRateLimitSurpassedError RateLimitSurpassed = func(err error) bool {
	var vOneErr *Error
	if !errors.As(err, &vOneErr) {
		return false
	}
	return vOneErr.Code == HTTPResponseTooManyRequests
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
			return ErrStop
		}
		err := v.callWithoutLimiter(ctx, f)
		err = v.rateLimiter.CheckError(err)
		if err == ErrOnceMore {
			continue
		}
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

	/*
		fmt.Printf("Method: %v\n", req.Method)
		fmt.Printf("URL: %v\n", req.URL)
		fmt.Printf("Proto: %v\n", req.Proto)
		fmt.Printf("ContentLength: %v\n", req.ContentLength)
		fmt.Printf("TransferEncoding: %v\n", req.TransferEncoding)
		fmt.Printf("Host: %v\n", req.Host)
		fmt.Printf("Form: %v\n", req.Form)
		fmt.Printf("PostForm: %v\n", req.PostForm)
		fmt.Printf("MultipartForm: %v\n", req.MultipartForm)
		fmt.Printf("Trailer: %v\n", req.Trailer)
		fmt.Printf("RemoteAddr: %v\n", req.RemoteAddr)
		fmt.Printf("ContentLength: %v\n", req.ContentLength)
	*/
	resp, err := v.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()
	if GetHTTPCodeRange(resp.StatusCode) != HTTPCodeSuccessRange {
		vOneErr, err := ErrorFromReader(resp.Body)
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}
		return &HTTPError{
			Status: resp.StatusCode,
			Err:    vOneErr,
		}
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
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read body : %w", err)
	}
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
		if !fieldValue.CanSet() {
			continue
		}
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
