package val

import (
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
	ToStr() string
	IsNil() bool
	IsNum() bool
}

type Num interface {
	Val
	ToFlt() float64
}

func (v Int) ToInt() int64 { return v.i }
func (v Flt) AsInt() int64 { return v.i }

func (v Int) ToFlt() float64 { return float64(v.i) }
func (v Flt) ToFlt() float64 { return float64(v.i / 100) }

func (v Int) ToStr() string { return v.s }
func (v Flt) ToStr() string { return v.s }
func (v Str) ToStr() string { return v.s }
func (v Nil) ToStr() string { return "" }

func (v Int) IsNil() bool { return v.s == "" }
func (v Flt) IsNil() bool { return v.s == "" }
func (v Str) IsNil() bool { return v.s == "" }
func (v Nil) IsNil() bool { return true }

func (v Int) IsNum() bool { return true }
func (v Flt) IsNum() bool { return true }
func (v Str) IsNum() bool { return false }
func (v Nil) IsNum() bool { return false }

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
	return strings.TrimRightFunc(s, func(r rune) bool {
		return r == '0' || r == '.'
	})
}
