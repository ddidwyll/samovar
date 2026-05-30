package val

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Flt struct {
	i int64
	s string
}

type Int struct {
	i int64
	s string
}

type Str struct {
	s string
}

type Nil struct{}

type Val interface {
	json.Marshaler
	String() string
	IsNil() bool
	IsNum() bool
	IsInt() bool
	IsFlt() bool
	ToFlt() float64
	ToInt() int64
}

func (v Int) ToInt() int64 { return v.i }
func (v Flt) ToInt() int64 { return v.i / 100 }
func (v Str) ToInt() int64 { panic("Val.str_to_int") }
func (v Nil) ToInt() int64 { panic("Val.nil_to_int") }
func (v Flt) AsInt() int64 { return v.i }

func (v Int) ToFlt() float64 { return float64(v.i) }
func (v Flt) ToFlt() float64 { return float64(v.i / 100) }
func (v Str) ToFlt() float64 { panic("Val.str_to_flt") }
func (v Nil) ToFlt() float64 { panic("Val.nil_to_flt") }

func (v Int) String() string { return v.s }
func (v Flt) String() string { return v.s }
func (v Str) String() string { return v.s }
func (v Nil) String() string { return "" }

func (v Int) IsNil() bool { return v.s == "" }
func (v Flt) IsNil() bool { return v.s == "" }
func (v Str) IsNil() bool { return v.s == "" }
func (v Nil) IsNil() bool { return true }

func (v Int) IsNum() bool { return true }
func (v Flt) IsNum() bool { return true }
func (v Str) IsNum() bool { return false }
func (v Nil) IsNum() bool { return false }

func (v Int) IsInt() bool { return true }
func (v Flt) IsInt() bool { return false }
func (v Str) IsInt() bool { return false }
func (v Nil) IsInt() bool { return false }

func (v Int) IsFlt() bool { return false }
func (v Flt) IsFlt() bool { return true }
func (v Str) IsFlt() bool { return false }
func (v Nil) IsFlt() bool { return false }

func FltAsFlt(f float64) Val {
	i := int64(math.Round(f * 100))
	s := fmt.Sprintf("%f", float64(i)/100.0)
	return Flt{i, trimZeros(s)}
}

func FltAsInt(f float64) Val {
	i := int64(math.Round(f))
	s := fmt.Sprintf("%d", i)
	return Int{i, s}
}

func FltAsStr(f float64) Val {
	i := int64(math.Round(f * 100))
	s := fmt.Sprintf("%f", float64(i)/100.0)
	return Str{trimZeros(s)}
}

func IntAsInt(i int64) Val {
	s := fmt.Sprintf("%d", i)
	return Int{i, s}
}

func IntAsFlt(i int64) Val {
	s := fmt.Sprintf("%d", i)
	return Flt{i * 100, s}
}

func IntAsStr(i int64) Val {
	s := fmt.Sprintf("%d", i)
	return Str{s}
}

func StrAsFlt(s string) (Val, error) {
	s = strings.TrimSpace(s)
	f, err := strconv.ParseFloat(s, 64)
	return FltAsFlt(f), err
}

func StrAsInt(s string) (Val, error) {
	s = strings.TrimSpace(s)
	f, err := strconv.ParseInt(s, 10, 64)
	return IntAsInt(f), err
}

func StrAsStr(s string) (v Val, err error) {
	if s = strings.TrimSpace(s); s == "" {
		err = errors.New("empty value")
	} else {
		v = Str{s}
	}
	return v, err
}

func trimZeros(s string) string {
	s = strings.TrimRight(s, "0")
	return strings.TrimSuffix(s, ".")
}

func (v Str) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.s)
}
func (v Flt) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.s)
}
func (v Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.s)
}
func (v Nil) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
