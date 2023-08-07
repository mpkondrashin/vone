package vone

import "fmt"

type Top int

const (
	Top50 Top = iota
	Top100
	Top200
	TopDefault = Top50
)

func (t Top) String() string {
	return [...]string{"50", "100", "200"}[t]
}

var IntToTop = map[int]Top{
	50:  Top50,
	100: Top100,
	200: Top200,
}

func TopFromInt(v int) (Top, error) {
	t, ok := IntToTop[v]
	if !ok {
		return 0, fmt.Errorf("%d: top invalid value", v)
	}
	return t, nil
}
