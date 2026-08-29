package change

import (
	"samovar/lib/val"

	"slices"
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

func (r Request) LastFrom() string {
	length := len(r.From)

	if length == 0 {
		return ""
	} else {
		return r.From[length-1]
	}
}

func (r Request) BuildReport(oldTs int64, oldV, newV val.Val, name, unit string) Report {
	isNew := oldTs == 0 && oldV.IsNil()
	changed := newV.String() != oldV.String()
	newFrom := slices.Clone(r.From)
	return Report{isNew, changed, r.Key, newFrom, oldV, newV, oldTs, r.Timestamp, name, unit}
}
