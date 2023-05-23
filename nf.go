package objectid

/*
nf.go provides NumberForm methods and types.

NOTE: uint128-related code written by Luke Champine, per https://github.com/lukechampine/uint128.
It has been incorporated into this package (unexported), and simplified to serve in the capacity
of OID arc storage.
*/

import (
	"fmt"
	"math/bits"
)

var nilNF NumberForm

/*
NumberForm is an unsigned 128-bit number. This type is based on
github.com/lukechampine/uint128. It has been incorporated
into this package to produce unsigned 128-bit OID numberForm
support (i.e.: UUID-based OIDs).
*/
type NumberForm struct {
	lo, hi uint64
	parsed bool
}

// isZero returns true if a == 0.
func (a *NumberForm) IsZero() bool {
	if a == nil {
		return true
	}

	// NOTE: we do not compare against Zero, because that
	// is a global variable that could be modified.
	return a.lo == uint64(0) && a.hi == uint64(0)
}

/*
Equal returns true if a == v.

NumberForm values can be compared directly with ==, but use of the
Equals method // is preferred for consistency.
*/
func (a NumberForm) Equal(v NumberForm) bool {
	return a == v
}

// Equal64 returns true if a == v.
func (a NumberForm) Equal64(v uint64) bool {
	return a.lo == v && a.hi == 0
}

/*
Compare compares a and v and returns:

	-1 if a <  v
	 0 if a == v
	+1 if a >  v
*/
func (a NumberForm) Compare(v NumberForm) int {
	if a == v {
		return 0
	} else if a.hi < v.hi || (a.hi == v.hi && a.lo < v.lo) {
		return -1
	} else {
		return 1
	}
}

/*
Compare64 compares a and v and returns:

	-1 if a <  v
	 0 if a == v
	+1 if a >  v
*/
func (a NumberForm) Compare64(v uint64) int {
	if a.hi == 0 && a.lo == v {
		return 0
	} else if a.hi == 0 && a.lo < v {
		return -1
	} else {
		return 1
	}
}

/*
Valid returns a boolean valud indicative of proper instantiation.
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

// String returns the base-10 representation of a as a string.
func (a NumberForm) String() string {
	if a.IsZero() {
		return "0"
	} else if !a.parsed {
		return "0"
	}

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

/*
Scan implements fmt.Scanner, and is only present to allow conversion
of an NumberForm into a string value per fmt.Sscan. Users need not execute
this method directly.
*/
func (a *NumberForm) Scan(s fmt.ScanState, ch rune) error {
	return sscan(s, ch, a)
}

// quoRem64 returns q = u/v and r = u%v.
func (a NumberForm) quoRem64(v uint64) (q NumberForm, r uint64) {
	if a.hi < v {
		q.lo, r = bits.Div64(a.hi, a.lo, v)
	} else {
		q.hi, r = bits.Div64(0, a.hi, v)
		q.lo, r = bits.Div64(r, a.lo, v)
	}
	return
}

// NewNumberForm returns the NumberForm value.
func newNumberForm(lo, hi uint64) NumberForm {
	return NumberForm{lo: lo, hi: hi, parsed: true}
}

/*
ParseNumberForm converts v into an instance of NumberForm, which
is returned alongside an error.

Acceptable input types are string, int and uint64. No decimal value,
whether string or int, can ever be negative.
*/
func ParseNumberForm(v any) (a NumberForm, err error) {
	switch tv := v.(type) {
	case string:
		_, err = fmt.Sscan(tv, &a)
	case int:
		if tv < 0 {
			err = errorf("A NumberForm cannot be negative")
			break
		}
		a = newNumberForm(uint64(tv), 0)
	case uint64:
		a = newNumberForm(tv, 0)
	}

	a.parsed = err == nil

	return
}
