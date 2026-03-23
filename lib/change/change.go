package change

import (
	"samovar/lib/val"

	"time"
)

type Request struct {
	Key       string
	Value     any
	Reason    string
	Timestamp int64
}

func NewRequest(key string, val any, reason string, ts int64) Request {
	return Request{key, val, reason, ts}
}

type Report struct {
	IsNew        bool
	Changed      bool
	Key          string
	Reason       string
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
	return Report{isNew, changed, r.Key, r.Reason, oldV, newV, oldTs, r.Timestamp, name, unit}
}

func (r Report) FormatTS() string {
	return time.UnixMicro(r.NewTimestamp).Format(time.TimeOnly)
}
