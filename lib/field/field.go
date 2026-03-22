package field

import (
  "samovar/lib/value"
)

type Type value.Type

type Field struct {
  Type Type
  rawValue value.Raw
}

func New(v string, t Type) Field {
  return Field{t, v}
}
