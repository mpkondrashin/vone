package vone

import (
	"encoding/json"
	"strings"
	"time"
)

type VOneTime time.Time

const layout = `2006-01-02T15:04:05`

// UnmarshalJSON for VOneTime with a custom layout
func (ct *VOneTime) UnmarshalJSON(data []byte) error {
	//fmt.Printf("PRZ: |%s|\n", string(data))
	s := strings.Trim(string(data), `"`)
	//fmt.Printf("S: |%s|, %d\n", s, len(s))
	if len(s) == 0 {
		*ct = VOneTime{}
		return nil
	}
	parsedTime, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	*ct = VOneTime(parsedTime)
	return nil
}

// Implement Marshaler interface for VOneTime
func (j VOneTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

// To format the output when printing
func (ct VOneTime) String() string {
	return time.Time(ct).Format(layout)
}
