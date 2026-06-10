package field

import (
	"samovar/lib/change"
	"samovar/lib/val"

	"errors"
	"fmt"
	"time"
)

type Type rune

type Field struct {
	Name string
	Unit string
	t    Type
	val  val.Val
	ts   int64
}

func New(name string, t Type) *Field {
	return &Field{name, "", t, val.Nil{}, 0}
}

func (f *Field) String() string {
	if f.Unit == "" {
		return f.Name
	} else {
		return fmt.Sprintf("%s (%s)", f.Name, f.Unit)
	}
}

func (f *Field) Get() val.Val { return f.val }

func (f *Field) Cast(a any) (v val.Val, err error) {
	switch f.t {
	case 'i':
		v, err = f.castInt(a)
	case 'f':
		v, err = f.castFlt(a)
	case 's':
		v, err = f.castStr(a)
	case 'a':
		v, err = f.castArr(a)
	default:
		v = val.Nil{}
		err = errors.New("invalid field type")
	}

	return v, err
}

func (f *Field) Set(a any) error {
	v, err := f.Cast(a)

	if err == nil {
		f.val = v
		f.ts = time.Now().UnixMicro()
	}

	return err
}

func (f *Field) Change(req change.Request) (rep change.Report, err error) {
	newVal, err := f.Cast(req.Value)

	if err != nil {
		return rep, err
	}

	rep = req.BuildReport(f.ts, f.val, newVal, f.Name, f.Unit)
	f.val, f.ts = newVal, req.Timestamp

	return rep, err
}

func (f *Field) IsNil() bool { return f.val.IsNil() }

func (f *Field) Type() Type { return f.t }

func (f *Field) castInt(a any) (v val.Val, err error) {
	switch c := a.(type) {
	case int:
		v = val.IntAsInt(int64(c))
	case int64:
		v = val.IntAsInt(c)
	case float64:
		v = val.FltAsInt(c)
	case string:
		v, err = val.StrAsInt(c)
	case val.Int:
		v = c
	case val.Flt:
		v = val.IntAsInt(c.ToInt())
	case val.Str:
		v, err = val.StrAsInt(c.String())
	default:
		v = val.Nil{}
		err = errors.New("failed to cast field")
	}

	return v, err
}

func (f *Field) castFlt(a any) (v val.Val, err error) {
	switch c := a.(type) {
	case int:
		v = val.IntAsFlt(int64(c))
	case int64:
		v = val.IntAsFlt(c)
	case float64:
		v = val.FltAsFlt(c)
	case string:
		v, err = val.StrAsFlt(c)
	case val.Int:
		v = val.IntAsFlt(c.ToInt())
	case val.Flt:
		v = c
	case val.Str:
		v, err = val.StrAsFlt(c.String())
	default:
		v = val.Nil{}
		err = errors.New("failed to cast field")
	}

	return v, err
}

func (f *Field) castStr(a any) (v val.Val, err error) {
	switch c := a.(type) {
	case int:
		v = val.IntAsStr(int64(c))
	case int64:
		v = val.IntAsStr(c)
	case float64:
		v = val.FltAsStr(c)
	case string:
		v, err = val.StrAsStr(c)
	case []string:
		v, err = val.ArrAsStr(c)
	case val.Int:
		v, err = val.StrAsStr(c.String())
	case val.Flt:
		v, err = val.StrAsStr(c.String())
	case val.Str:
		v = c
	case val.Arr:
		v, err = val.StrAsStr(c.String())
	default:
		v = val.Nil{}
		err = errors.New("failed to cast field")
	}

	return v, err
}

func (f *Field) castArr(a any) (v val.Val, err error) {
	switch c := a.(type) {
	case int:
		v = val.IntAsArr(int64(c))
	case int64:
		v = val.IntAsArr(c)
	case float64:
		v = val.FltAsArr(c)
	case string:
		v, err = val.StrAsArr(c)
	case []string:
		v, err = val.ArrAsArr(c)
	case val.Int:
		v, err = val.StrAsArr(c.String())
	case val.Flt:
		v, err = val.StrAsArr(c.String())
	case val.Str:
		v, err = val.StrAsArr(c.String())
	case val.Arr:
		v = c
	default:
		v = val.Nil{}
		err = errors.New("failed to cast field")
	}

	return v, err
}
