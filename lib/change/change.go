package change

import "samovar/lib/val"

type Request struct {
	Key       string
	Value     any
	Timestamp int64
}

type Report struct {
	Changed      bool
	Key          string
	OldValue     val.Val
	NewValue     val.Val
	OldTimestamp int64
	NewTimestamp int64
}

func (r Request) BuildReport(oldTs int64, oldV, newV val.Val) Report {
	changed := newV.ToStr() != oldV.ToStr()
	return Report{changed, r.Key, oldV, newV, oldTs, r.Timestamp}
}
