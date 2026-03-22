package value

import (
  "math"
  "fmt"
  "strconv"
  "strings"
  "errors"
)

type Type rune

type Value struct {
  Type Type
  str string
  num int64
}

func FromFloat(f float64) Value {
  num := int64(math.Round(f * 100))
  str := fmt.Sprintf("%d", num / 100)
  return Value{'f', str, num}
}

func FromInt(i int64) Value {
  return Value{'i', string(i), i}
}

func Parse(s string, t Type) (v Value, err error) {
  switch t {
  case 'i':
    if i, err := strconv.ParseInt(s, 10, 64); err == nil {
      return FromInt(i), nil
    } else {
      return v, err
    }
  case 'f':
    if f, err := strconv.ParseFloat(s, 64); err == nil {
      return FromFloat(f), nil
    } else {
      return v, err
    }
  case 's':
    if s = strings.TrimSpace(s); s == "" {
      return v, errors.New("empty value")
    } else {
      return Value{'s', s, 0}, nil
    }
  }
  return v, errors.New("unexpected value type")
}

func (v Value) ToInt() int64 {
  if v.Type != 'i' && v.Type != 'f'{
    panic("not integer")
  }

  return v.num
}

func (v Value) ToFloat() float64 {
  if v.Type != 'f' && v.Type != 'i' {
    panic("not float")
  }

  return float64(v.num / 100)
}
