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

type Arr struct {
	a []string
	s string
}

type Nil struct{}

type Val interface {
	json.Marshaler
	T() rune
	Eq(Val) bool
	EqStr(string) bool
	String() string
	IsNil() bool
	IsNum() bool
	IsInt() bool
	IsFlt() bool
	ToFlt() float64
	ToInt() int64
	ToArr() []string
}

func (v Int) T() rune { return 'i' }
func (v Flt) T() rune { return 'f' }
func (v Str) T() rune { return 's' }
func (v Nil) T() rune { return 'n' }
func (v Arr) T() rune { return 'a' }

func (v Int) ToInt() int64 { return v.i }
func (v Flt) ToInt() int64 { return v.i / 100 }
func (v Str) ToInt() int64 { panic("Val.str_to_int") }
func (v Nil) ToInt() int64 { panic("Val.nil_to_int") }
func (v Arr) ToInt() int64 { panic("Val.arr_to_int") }
func (v Flt) AsInt() int64 { return v.i }

func (v Int) ToFlt() float64 { return float64(v.i) }
func (v Flt) ToFlt() float64 { return float64(v.i) / 100.0 }
func (v Arr) ToFlt() float64 { panic("Val.arr_to_flt") }
func (v Str) ToFlt() float64 { panic("Val.str_to_flt") }
func (v Nil) ToFlt() float64 { panic("Val.nil_to_flt") }

func (v Int) ToArr() []string { return []string{v.s} }
func (v Flt) ToArr() []string { return []string{v.s} }
func (v Str) ToArr() []string { return []string{v.s} }
func (v Arr) ToArr() []string { return v.a }
func (v Nil) ToArr() []string { return []string{} }

func (v Int) String() string { return v.s }
func (v Flt) String() string { return v.s }
func (v Str) String() string { return v.s }
func (v Arr) String() string { return v.s }
func (v Nil) String() string { return "" }

func (v Int) IsNil() bool { return v.s == "" }
func (v Flt) IsNil() bool { return v.s == "" }
func (v Str) IsNil() bool { return v.s == "" }
func (v Arr) IsNil() bool { return v.s == "" }
func (v Nil) IsNil() bool { return true }

func (v Int) IsNum() bool { return !v.IsNil() }
func (v Flt) IsNum() bool { return !v.IsNil() }
func (v Str) IsNum() bool { return false }
func (v Arr) IsNum() bool { return false }
func (v Nil) IsNum() bool { return false }

func (v Int) IsInt() bool { return !v.IsNil() }
func (v Flt) IsInt() bool { return false }
func (v Str) IsInt() bool { return false }
func (v Arr) IsInt() bool { return false }
func (v Nil) IsInt() bool { return false }

func (v Int) IsFlt() bool { return false }
func (v Flt) IsFlt() bool { return !v.IsNil() }
func (v Str) IsFlt() bool { return false }
func (v Arr) IsFlt() bool { return false }
func (v Nil) IsFlt() bool { return false }

func (a Int) Eq(b Val) bool { return Eq(a, b) }
func (a Flt) Eq(b Val) bool { return Eq(a, b) }
func (a Str) Eq(b Val) bool { return Eq(a, b) }
func (a Arr) Eq(b Val) bool { return Eq(a, b) }
func (a Nil) Eq(b Val) bool { return Eq(a, b) }

func (v Int) EqStr(s string) bool { return v.String() == s }
func (v Flt) EqStr(s string) bool { return v.String() == s }
func (v Str) EqStr(s string) bool { return v.String() == s }
func (v Arr) EqStr(s string) bool { return v.String() == s }
func (v Nil) EqStr(s string) bool { return v.String() == s }

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

func FltAsArr(f float64) Val {
	i := int64(math.Round(f * 100))
	s := fmt.Sprintf("%f", float64(i)/100.0)
	v, _ := StrAsArr(trimZeros(s))
	return v
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

func IntAsArr(i int64) Val {
	s := fmt.Sprintf("%d", i)
	v, _ := StrAsArr(s)
	return v
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

func StrAsArr(s string) (v Val, err error) {
	arr := strings.Split(s, ";")
	return ArrAsArr(arr)
}

func ArrAsArr(arr []string) (Val, error) {
	for i, str := range arr {
		if str = strings.TrimSpace(str); str == "" {
			return Nil{}, errors.New("empty value")
		} else {
			arr[i] = str
		}
	}
	if len(arr) == 0 {
		return Nil{}, errors.New("empty value")
	}
	v := Arr{arr, strings.Join(arr, ";")}
	return v, nil
}

func ArrAsStr(arr []string) (Val, error) {
	if v, err := ArrAsArr(arr); err == nil {
		return Str{v.String()}, nil
	} else {
		return Nil{}, err
	}
}

func ArrAsInt(_ []string) (Val, error) {
	return Nil{}, errors.New("Val.arr_as_int")
}

func ArrAsFlt(_ []string) (Val, error) {
	return Nil{}, errors.New("Val.arr_as_int")
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

func (v Arr) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.s)
}

func (v Nil) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}

func Eq(a Val, b Val) bool {
	return a.T() == b.T() && a.String() == b.String()
}
