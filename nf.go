package objectid

/*
nf.go provides NumberForm methods and types.

NOTE: uint128-related code written by Luke Champine, per https://github.com/lukechampine/uint128.
It has been incorporated into this package (unexported), and simplified to serve in the capacity
of OID numberForm storage.
*/

import (
	"fmt"
	"math/bits"
)

var nilNF NumberForm

/*
NumberForm is an unsigned 128-bit number. This type is based on
github.com/lukechampine/uint128. It has been incorporated into
this package to produce unsigned 128-bit [OID] numberForm support
(i.e.: UUID-based [OID]s).
*/
type NumberForm struct {
	lo, hi uint64
	parsed bool
}

/*
IsZero returns a Boolean value indicative of whether the
receiver instance is nil, or unset.
*/
func (a *NumberForm) IsZero() (is bool) {
	if is = a == nil; !is {
		// NOTE: we do not compare against Zero, because that
		// is a global variable that could be modified.
		is = (a.lo == uint64(0) && a.hi == uint64(0))
	}
	return
}

/*
Equal returns a boolean value indicative of whether the receiver is equal to
the value provided. Valid input types are string, uint64, int and [NumberForm].

Any input that represents a negative number guarantees a false return.
*/
func (a NumberForm) Equal(n any) (is bool) {
	switch tv := n.(type) {
	case NumberForm:
		is = a == tv
	case string:
		if nf, err := NewNumberForm(tv); err == nil {
			is = a == nf
		}
	case uint64:
		is = a.lo == tv && a.hi == 0
	case int:
		if 0 <= tv {
			is = a.lo == uint64(tv) && a.hi == 0
		}
	}

	return
}

/*
Gt returns a boolean value indicative of whether the receiver is greater than
the value provided. Valid input types are string, uint64, int and [NumberForm].

Any input that represents a negative number guarantees a false return.
*/
func (a NumberForm) Gt(n any) (is bool) {
	switch tv := n.(type) {
	case NumberForm, string:
		is = a.gtLt(tv, false)
	case uint64:
		is = a.lo > tv && a.hi == uint64(0)
	case int:
		if 0 <= tv {
			is = a.lo > uint64(tv) && a.hi == uint64(0)
		}
	}
	return
}

/*
Ge returns a boolean value indicative of whether the receiver is greater than
or equal to the value provided. Valid input types are string, uint64, int and
[NumberForm]. This method is merely a convenient wrapper to an ORed call of the
[NumberForm.Gt] and [NumberForm.Equal] methods.

Any input that represents a negative number guarantees a false return.
*/
func (a NumberForm) Ge(n any) (is bool) {
	return a.Gt(n) || a.Equal(n)
}

/*
Lt returns a boolean value indicative of whether the receiver is less than
the value provided. Valid input types are string, uint64, int and [NumberForm].

Any input that represents a negative number guarantees a false return.
*/
func (a NumberForm) Lt(n any) (is bool) {
	switch tv := n.(type) {
	case NumberForm, string:
		is = a.gtLt(tv, true)
	case uint64:
		is = a.lo < tv && a.hi == uint64(0)
	case int:
		if 0 <= tv {
			is = a.lo < uint64(tv) && a.hi == uint64(0)
		}
	}
	return
}

/*
Le returns a boolean value indicative of whether the receiver is less than or
equal to the value provided. Valid input types are string, uint64, int and
[NumberForm]. This method is merely a convenient wrapper to an ORed call of the
[NumberForm.Lt] and [NumberForm.Equal] methods.

Any input that represents a negative number guarantees a false return.
*/
func (a NumberForm) Le(n any) (is bool) {
	return a.Lt(n) || a.Equal(n)
}

func (a NumberForm) gtLt(x any, lt bool) bool {
	var nf NumberForm

	switch tv := x.(type) {
	case string:
		nf, _ = NewNumberForm(tv)
	case NumberForm:
		nf = tv
	default:
		return false
	}

	if lt {
		return a.hi < nf.hi || (a.hi == nf.hi && a.lo < nf.lo)
	}
	return a.hi > nf.hi || (a.hi == nf.hi && a.lo > nf.lo)
}

/*
Valid returns a Boolean value indicative of proper initialization.
*/
func (a NumberForm) Valid() bool {
	return a.parsed
}

/*
leadingZeros returns the number of leading zero bits in u;
the result is 128 for a == 0.
*/
func (a NumberForm) leadingZeros() int {
	if a.hi > 0 {
		return bits.LeadingZeros64(a.hi)
	}
	return 64 + bits.LeadingZeros64(a.lo)
}

/*
Len returns the minimum number of bits required to represent u;
the result is 0 for a == 0.
*/
func (a NumberForm) len() int {
	return 128 - a.leadingZeros()
}

/*
String returns the base-10 string representation of the receiver
instance.
*/
func (a NumberForm) String() string {
	if !a.IsZero() {
		buf := []byte("0000000000000000000000000000000000000000") // log10(2^128) < 40
		for i := len(buf); ; i -= 19 {
			q, r := a.quoRem64(1e19) // largest power of 10 that fits in a uint64
			var n int
			for ; r != 0; r /= 10 {
				n++
				buf[i-n] += byte(r % 10)
			}
			if q.IsZero() {
				return string(buf[i-n:])
			}
			a = q
		}
	}

	return "-1"
}

/*
Scan implements fmt.Scanner, and is only present to allow conversion
of an [NumberForm] into a string value per [fmt.Sscan].  Users need not
execute this method directly.
*/
func (a *NumberForm) Scan(s fmt.ScanState, ch rune) error {
	return sscan(s, ch, a)
}

// quoRem64 returns q = u/v and r = u%v.
// Credit: Luke Champine
func (a NumberForm) quoRem64(v uint64) (q NumberForm, r uint64) {
	if a.hi < v {
		q.lo, r = bits.Div64(a.hi, a.lo, v)
	} else {
		q.hi, r = bits.Div64(0, a.hi, v)
		q.lo, r = bits.Div64(r, a.lo, v)
	}
	return
}

// NewNumberForm returns the [NumberForm] value.
func newNumberForm(lo, hi uint64) NumberForm {
	return NumberForm{lo: lo, hi: hi, parsed: true}
}

/*
NewNumberForm converts v into an instance of [NumberForm], which
is returned alongside an error.

Acceptable input types are string, int and uint64. No decimal value,
whether string or int, can ever be negative.
*/
func NewNumberForm(v any) (a NumberForm, err error) {
	switch tv := v.(type) {
	case string:
		_, err = fmt.Sscan(tv, &a)
		a.parsed = err == nil
	case int:
		if tv < 0 {
			err = errorf("A NumberForm cannot be negative")
			break
		}
		a = newNumberForm(uint64(tv), 0)
	case uint64:
		a = newNumberForm(tv, 0)
	}

	return
}
