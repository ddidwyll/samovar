package state

import (
	"samovar/lib/field"
	"samovar/lib/val"

	"errors"
	"fmt"
)

type Key string

type State map[Key]*field.Field

type FieldParams struct {
	Key  Key
	Type field.Type
	Name string
}

type Fields []FieldParams

func New(fields Fields) *State {
	newState := make(State)

	for _, f := range fields {
		newState[f.Key] = field.New(f.Name, f.Type)
	}

	return &newState
}

func (s *State) Fetch(k Key) (f *field.Field, err error) {
	if field, has := (*s)[k]; has {
		return field, err
	} else {
		err := fmt.Sprintf("field [%s] not found", k)
		return f, errors.New(err)
	}
}

func (s *State) Get(k Key) val.Val {
	if field, has := (*s)[k]; has {
		return field.Get()
	} else {
		return val.Nil{}
	}
}

func (s *State) Has(k Key) bool {
	_, has := (*s)[k]
	return has
}

func (s *State) Set(k Key, v any) error {
	if field, has := (*s)[k]; has {
		return field.Set(v)
	} else {
		err := fmt.Sprintf("field [%s] not found", k)
		return errors.New(err)
	}
}
