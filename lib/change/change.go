package change

import (
	"samovar/lib/val"

	"fmt"
	"slices"
	"time"
)

func NewRequest(key string, val any, from string, ts int64) Request {
	return Request{key, val, []string{from}, ts}
}

type Request struct {
	Key       string
	Value     any
	From      []string
	Timestamp int64
}

type Report struct {
	IsNew        bool
	Changed      bool
	Key          string
	From         []string
	OldValue     val.Val
	NewValue     val.Val
	OldTimestamp int64
	NewTimestamp int64
	FieldName    string
	FieldUnit    string
}

func (r Request) BuildReport(oldTs int64, oldV, newV val.Val, name, unit string) Report {
	isNew := oldTs == 0 && oldV.IsNil()
	changed := newV.String() != oldV.String()
	newFrom := slices.Clone(r.From)
	return Report{isNew, changed, r.Key, newFrom, oldV, newV, oldTs, r.Timestamp, name, unit}
}

func (r Report) NewRequest() (req Request) {
	req.Key = r.Key
	req.Value = r.NewValue
	req.Timestamp = r.NewTimestamp
	req.From = slices.Clone(r.From)
	return req
}

func (r Report) FormatTime() string {
	return time.UnixMicro(r.NewTimestamp).Format(time.TimeOnly)
}

func (r Report) FormatValue() string {
	if r.OldValue.IsNil() {
		return r.NewValue.String()
	} else {
		return fmt.Sprintf("%s -> %s", r.OldValue, r.NewValue)
	}
}

func (r Report) FormatField() string {
	if r.FieldUnit == "" {
		return r.FieldName
	} else {
		return fmt.Sprintf("%s (%s)", r.FieldName, r.FieldUnit)
	}
}

func (r Report) AddFrom(from string) Report {
	r.From = append(r.From, from)
	return r
}

func (r Report) LastFrom() string {
	length := len(r.From)

	if length == 0 {
		return ""
	} else {
		return r.From[length-1]
	}
}
