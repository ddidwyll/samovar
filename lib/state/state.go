package state

import (
  "samovar/lib/field"
)

type Key string

type State map[Key]field.Field

func NewState() *State {
  s := make(State)
  return &s
}
