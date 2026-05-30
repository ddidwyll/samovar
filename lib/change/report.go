package change

import (
	"samovar/lib/val"

	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

type Report struct {
	IsNew        bool     `json:"-"`
	Changed      bool     `json:"-"`
	Key          string   `json:"key"`
	From         []string `json:"-"`
	OldValue     val.Val  `json:"-"`
	NewValue     val.Val  `json:"value"`
	OldTimestamp int64    `json:"-"`
	NewTimestamp int64    `json:"-"`
	FieldName    string   `json:"-"`
	FieldUnit    string   `json:"-"`
}

func (r Report) ReRequest() (req Request) {
	req.Key = r.Key
	req.Value = r.NewValue
	req.Timestamp = r.NewTimestamp
	req.From = slices.Clone(r.From)
	return req
}

func (r Report) NewRequest(from, key string, val any) (req Request) {
	req.Key = key
	req.Value = val
	req.Timestamp = r.NewTimestamp
	req.From = append(r.From, from)
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

func (r Report) MakeLogString(prefix string) string {
	field := r.FormatField()
	field = strings.ReplaceAll(field, "%", "%%")

	if r.IsNew {
		return fmt.Sprintf(
			"%s [%s\t]:\t%s",
			prefix,
			field,
			r.NewValue,
		)
	} else {
		return fmt.Sprintf(
			"%s [%s\t]:\t%s -> %s\t\t%s\t%s",
			prefix,
			field,
			r.OldValue,
			r.NewValue,
			r.FormatTime(),
			r.LastFrom(),
		)
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

func (r Report) ToJson() []byte {
	json, _ := json.Marshal(r)
	return json
}
