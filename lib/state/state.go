package state

import (
	"samovar/lib/change"
	"samovar/lib/field"
	"samovar/lib/val"

	"encoding/json"
	"errors"
	"fmt"
)

type State map[string]*field.Field
type KeyVals map[string]val.Val

type Entry struct {
	Key   string  `json:"key"`
	Value val.Val `json:"value"`
}
type Entries []Entry

type FieldParams struct {
	Key  string
	Type field.Type
	Name string
	Unit string
}

type Fields []FieldParams

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

func (s *State) KeyVals() KeyVals {
	kv := make(KeyVals, len(*s))
	for key, field := range *s {
		kv[key] = field.Get()
	}
	return kv
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

func (s *State) Change(req change.Request) (rep change.Report, err error) {
	if field, has := (*s)[req.Key]; has {
		return field.Change(req)
	} else {
		err := fmt.Sprintf("field [%s] not found", req.Key)
		return rep, errors.New(err)
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
