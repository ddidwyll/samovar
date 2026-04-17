package state

import (
	"samovar/lib/change"
	"samovar/lib/field"
	"samovar/lib/val"

	"errors"
	"fmt"
)

type State map[string]*field.Field

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
		keys := [2]string{params.Key, params.Name}

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
