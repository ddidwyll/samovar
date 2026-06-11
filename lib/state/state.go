package state

import (
	"samovar/lib/change"
	"samovar/lib/field"
	"samovar/lib/val"

	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type State map[string]*field.Field
type KeyVals map[string]val.Val

type Entry struct {
	Key   string  `json:"key"`
	Value val.Val `json:"value"`
}
type Entries []Entry

type FieldParams struct {
	Key  string     `json:"key"`
	Type field.Type `json:"-"`
	Name string     `json:"name"`
	Unit string     `json:"unit"`
}

type Fields []FieldParams

type request = change.Request
type report = change.Report

func New(fields Fields) *State {
	newState := make(State)

	for _, params := range fields {
		field := field.New(params.Name, params.Type)
		field.Unit = params.Unit
		// STATE BY NAME
		// keys := [...]string{params.Key, params.Name}
		keys := [...]string{params.Key}

		for _, key := range keys {
			if key == "" {
				panic("empty state key")
			}

			_, exists := newState[key]

			if exists && params.Key != params.Name {
				panic("duplicate state key")
			}

			newState[key] = field
		}
	}

	return &newState
}

func (s *State) AllKeyVals() KeyVals {
	kv := make(KeyVals, len(*s))
	for key, field := range *s {
		kv[key] = field.Get()
	}
	return kv
}

func (s *State) KeyVals(keys ...string) (KeyVals, error) {
	kv := make(KeyVals, len(keys))
	for _, key := range keys {
		if f, err := s.FetchField(key); err != nil {
			return nil, err
		} else {
			kv[key] = f.Get()
		}
	}
	return kv, nil
}

func (s *State) Entries() Entries {
	entries := make(Entries, 0, len(*s))
	for key, field := range *s {
		entry := Entry{key, field.Get()}
		entries = append(entries, entry)
	}
	return entries
}

func (s *State) HasField(key string) bool {
	_, has := (*s)[key]
	return has
}

func (s *State) FetchField(k string) (f *field.Field, err error) {
	if field, has := (*s)[k]; has {
		return field, err
	} else {
		err := fmt.Sprintf("field [%s] not found", k)
		return f, errors.New(err)
	}
}

func (s *State) Get(k string) val.Val {
	if field, has := (*s)[k]; has {
		return field.Get()
	} else {
		return val.Nil{}
	}
}

func (s *State) Has(k string) bool {
	_, has := (*s)[k]
	return has
}

func (s *State) Set(k string, v any) error {
	if field, has := (*s)[k]; has {
		return field.Set(v)
	} else {
		err := fmt.Sprintf("field [%s] not found", k)
		return errors.New(err)
	}
}

func (s *State) Change(req request) (rep report, err error) {
	if field, has := (*s)[req.Key]; has {
		return field.Change(req)
	} else {
		err := fmt.Sprintf("field [%s] not found", req.Key)
		return rep, errors.New(err)
	}
}

func (s *State) BulkChangeFromMap(kv map[string]any, from, lastFrom string) ([]report, error) {
	ts := time.Now().UnixMicro()
	reqs := make([]request, 0, len(kv))
	for k, v := range kv {
		req := change.NewRequest(k, v, from, ts)
		reqs = append(reqs, req)
	}
	return s.BulkChange(reqs, lastFrom)
}

func (s *State) BulkChange(reqs []request, lastFrom string) ([]report, error) {
	reports := make([]report, 0, len(reqs))

	for _, req := range reqs {
		if rep, err := s.Change(req); err != nil {
			return reports, err
		} else {
			if rep.Changed {
				reports = append(reports, rep.AddFrom(lastFrom))
			}
		}
	}

	return reports, nil
}

func (s *State) HandleReq(req any) (any, error) {
	switch k := req.(type) {
	case []string:
		return s.KeyVals(k...)
	case string:
		return s.FetchField(k)
	case nil:
		return s.Entries(), nil
	default:
		return nil, errors.New("Unexpected state request")
	}
}

func (kv *KeyVals) ToJson() []byte {
	json, _ := json.Marshal(kv)
	return json
}

func (entries *Entries) ToJson() []byte {
	json, _ := json.Marshal(entries)
	return json
}

func (fields *Fields) ToJson() []byte {
	json, _ := json.Marshal(fields)
	return json
}
