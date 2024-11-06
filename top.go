/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	top.go - limit of API calls response type
*/

package vone

import "fmt"

type Top int

const (
	Top50 Top = iota
	Top100
	Top200
	Top500
	Top1000
)

func (t Top) String() string {
	return [...]string{"50", "100", "200", "500", "1000"}[t]
}

func (t Top) Int() int {
	return [...]int{50, 100, 200, 500, 1000}[t]
}

var IntToTop = map[int]Top{
	50:   Top50,
	100:  Top100,
	200:  Top200,
	500:  Top500,
	1000: Top1000,
}

func TopFromInt(v int) (Top, error) {
	t, ok := IntToTop[v]
	if !ok {
		return 0, fmt.Errorf("%d: top invalid value", v)
	}
	return t, nil
}

func TopFromString(v string) (Top, error) {
	switch v {
	case "50":
		return Top50, nil
	case "100":
		return Top100, nil
	case "200":
		return Top200, nil
	case "500":
		return Top500, nil
	case "1000":
		return Top1000, nil
	default:
		return 0, fmt.Errorf("%s: top invalid value", v)
	}
}
