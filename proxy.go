package vone

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/launchdarkly/go-ntlm-proxy-auth"
)

func AddTransportModifier(pModifier *func(*http.Transport), modifier func(*http.Transport)) {
	if *pModifier == nil {
		*pModifier = modifier
		return
	}
	original := *pModifier
	*pModifier = func(t *http.Transport) {
		original(t)
		modifier(t)
	}
}

type AuthType int

const (
	AuthTypeNone AuthType = iota
	AuthTypeBasic
	AuthTypeNTLM
)

var AuthTypeString = [...]string{
	"None",
	"Basic",
	"NTLM",
}

func (r AuthType) String() string {
	return AuthTypeString[r]
}

// ErrUnknownAuthType - unknown authentication type error
var ErrUnknownAuthType = errors.New("unknown auth type")

// UnmarshalJSON implements the Unmarshaler interface of the json package for AuthType.
func (a *AuthType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	for i, s := range AuthTypeString {
		if strings.EqualFold(s, v) {
			*a = AuthType(i)
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrUnknownAuthType, v)
}

// MarshalJSON implements the Marshaler interface of the json package for AuthType.
func (a AuthType) MarshalJSON() ([]byte, error) {
	if a < 0 || a >= AuthTypeNTLM {
		return nil, ErrUnknownAuthType
	}
	return []byte(fmt.Sprintf("\"%s\"", a.String())), nil
}

type Proxy struct {
	Type      AuthType
	URL       *url.URL
	Username  string
	Password  string
	Domain    string
	Timeout   time.Duration
	KeepAlive time.Duration
}

func NewProxy(URL *url.URL) *Proxy {
	return &Proxy{
		Type: AuthTypeNone,
		URL:  URL,
	}
}

func (p *Proxy) BasicAuth(Username string, Password string) *Proxy {
	p.Type = AuthTypeBasic
	p.Username = Username
	p.Password = Password
	return p
}

func (p *Proxy) NTLMAuth(Username string, Password string, Domain string) *Proxy {
	p.Type = AuthTypeNTLM
	p.Username = Username
	p.Password = Password
	p.Domain = Domain
	return p
}

func (p *Proxy) GetModifier() func(*http.Transport) {
	switch p.Type {
	default:
		fallthrough
	case AuthTypeNone:
		return p.TransportNoAuth
	case AuthTypeBasic:
		return p.TransportBasic
	case AuthTypeNTLM:
		return p.TransportNTLM
	}
}

func (p *Proxy) ChangeTransport(t *http.Transport) {
	switch p.Type {
	default:
		fallthrough
	case AuthTypeNone:
		p.TransportNoAuth(t)
	case AuthTypeBasic:
		p.TransportBasic(t)
	case AuthTypeNTLM:
		p.TransportNTLM(t)
	}
}

func (p *Proxy) TransportNoAuth(t *http.Transport) {
	t.Proxy = http.ProxyURL(p.URL)
}

func (p *Proxy) TransportNTLM(t *http.Transport) {
	dialer := &net.Dialer{
		Timeout:   p.Timeout,
		KeepAlive: p.KeepAlive,
	}
	ntlmDialContext := ntlm.NewNTLMProxyDialContext(dialer, *p.URL, p.Username, p.Password, p.Domain, nil)
	t.Proxy = nil
	t.DialContext = ntlmDialContext

}

func (p *Proxy) TransportBasic(t *http.Transport) {
	u := *p.URL
	u.User = url.UserPassword(p.Username, p.Password)
	t.Proxy = http.ProxyURL(&u)
}
