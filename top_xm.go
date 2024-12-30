/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	top.go - limit of API calls response type
*/

package vone

import "fmt"

type TopXM int

const (
	TopXM10 TopXM = iota
	TopXM50
	TopXM100
	TopXM200
	TopXM1000
)

func (t TopXM) String() string {
	return [...]string{"10", "50", "100", "200", "1000"}[t]
}

func (t TopXM) Int() int {
	return [...]int{10, 50, 100, 200, 1000}[t]
}

var IntToTopXM = map[int]TopXM{
	10:   TopXM10,
	50:   TopXM50,
	100:  TopXM100,
	200:  TopXM200,
	1000: TopXM1000,
}

func TopXMFromInt(v int) (TopXM, error) {
	t, ok := IntToTopXM[v]
	if !ok {
		return 0, fmt.Errorf("%d: top invalid value", v)
	}
	return t, nil
}

func TopXMFromString(v string) (TopXM, error) {
	switch v {
	case "10":
		return TopXM10, nil
	case "50":
		return TopXM50, nil
	case "100":
		return TopXM100, nil
	case "200":
		return TopXM200, nil
	case "1000":
		return TopXM1000, nil
	default:
		return 0, fmt.Errorf("%s: top invalid value", v)
	}
}
