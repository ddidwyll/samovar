package field

import (
	"samovar/lib/val"

	"errors"
)

type Type rune

type Field struct {
	Name string
	t    Type
	val  val.Val
	timestamp int64
}

func New(n string, t Type) *Field {
	return &Field{n, t, val.Nil{}}
}

func (f *Field) Get() val.Val { return f.val }

func (f *Field) Set(v any) error {
	switch f.t {
	case 'i':
		return f.castInt(v)
	case 'f':
		return f.castFlt(v)
	case 's':
		return f.castStr(v)
	default:
		return errors.New("invalid field type")
	}
}

func (f *Field) Update(fun func(val.Val) any) error {
	return f.Set(fun(f.val))
}

func (f *Field) IsNil() bool { return f.val.IsNil() }

func (f *Field) Type() Type { return f.t }

func (f *Field) castInt(v any) (err error) {
	switch c := v.(type) {
	case int:
		f.val = val.IntAsInt(int64(c))
	case int64:
		f.val = val.IntAsInt(c)
	case float64:
		f.val = val.FltAsInt(c)
	case string:
		f.val, err = val.StrAsInt(c)
	default:
		err = errors.New("failed to cast field")
	}

	return err
}

func (f *Field) castFlt(v any) (err error) {
	switch c := v.(type) {
	case int:
		f.val = val.IntAsFlt(int64(c))
	case int64:
		f.val = val.IntAsFlt(c)
	case float64:
		f.val = val.FltAsFlt(c)
	case string:
		f.val, err = val.StrAsFlt(c)
	default:
		err = errors.New("failed to cast field")
	}

	return err
}

func (f *Field) castStr(v any) (err error) {
	switch c := v.(type) {
	case int:
		f.val = val.IntAsStr(int64(c))
	case int64:
		f.val = val.IntAsStr(c)
	case float64:
		f.val = val.FltAsStr(c)
	case string:
		f.val, err = val.StrAsStr(c)
	default:
		err = errors.New("failed to cast field")
	}

	return err
}
