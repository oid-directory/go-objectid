package objectid

import "math/big"

/*
DotNotation contains an ordered sequence of [NumberForm] instances.
*/
type DotNotation []NumberForm

/*
String is a stringer method that returns the dot notation form of the
receiver (e.g.: "1.3.6.1").
*/
func (r DotNotation) String() (s string) {
	if !r.IsZero() {
		var x []string
		for i := 0; i < len(r); i++ {
			x = append(x, r[i].String())
		}

		s = join(x, `.`)
	}
	return
}

/*
Root returns the root node (0) [NumberForm] instance.
*/
func (r DotNotation) Root() NumberForm {
	x, _ := r.Index(0)
	return x
}

func (r DotNotation) Len() int {
	return len(r)
}

/*
Leaf returns the leaf-node (-1) [NumberForm] instance.
*/
func (r DotNotation) Leaf() NumberForm {
	x, _ := r.Index(-1)
	return x
}

/*
Parent returns the leaf-node's parent (-2) [NumberForm] instance.
*/
func (r DotNotation) Parent() NumberForm {
	x, _ := r.Index(-2)
	return x
}

/*
IsZero returns a Boolean indicative of whether the receiver is unset.
*/
func (r *DotNotation) IsZero() (is bool) {
	if r != nil {
		is = r.Len() == 0
	}
	return
}

/*
NewDotNotation returns an instance of *[DotNotation] alongside a Boolean
value indicative of success.

Variadic input allows for slice mixtures of all of the following types,
each treated as an individual [NumberForm] instance:

  - *[math/big.Int]
  - [NumberForm]
  - string
  - uint64
  - uint
  - int

If a string primitive is the only input option, it will be treated as a
complete [DotNotation] (e.g.: "1.3.6").
*/
func NewDotNotation(x ...any) (r *DotNotation, err error) {
	var _d DotNotation = make(DotNotation, 0)

	if len(x) == 1 {
		if slice, ok := x[0].(string); ok {
			r, err = newDotNotationStr(slice)
			return
		}
	}

	for i := 0; i < len(x) && err == nil; i++ {
		var nf NumberForm
		switch tv := x[i].(type) {
		case NumberForm:
			if !tv.Valid() {
				err = errorf("Unsupported %T value", tv)
				break
			}
			nf = tv
		case *big.Int, string, uint64, uint, int:
			nf, err = NewNumberForm(tv)
		default:
			err = errorf("Unsupported slice type '%T' for OID", tv)
		}

		_d = append(_d, nf)
	}

	if err == nil {
		r = new(DotNotation)
		*r = _d
	}

	return
}

func newDotNotationStr(dot string) (r *DotNotation, err error) {
	if !isNumericOID(dot) {
		err = errorf("Invalid OID '%s' cannot be processed", dot)
		return
	}
	z := split(dot, `.`)

	_d := make(DotNotation, 0)
	for j := 0; j < len(z) && err == nil; j++ {
		var nf NumberForm
		if nf, err = NewNumberForm(z[j]); err == nil {
			_d = append(_d, nf)
		}
	}

	if err == nil {
		r = new(DotNotation)
		*r = _d
	}

	return
}

/*
Encode returns the ASN.1 encoding of the receiver instance alongside an error.
*/
func (r DotNotation) Encode() (b []byte, err error) {
	if r.Len() < 2 {
		err = errorf("Length below encoding minimum")
		return
	}

	var start int
	firstArc := r[0].cast()
	forty := big.NewInt(40)
	firstArc.Mul(firstArc, forty)

	if r[1].cast().Cmp(forty) < 1 {
		// this is meant for second level arcs <= 39
		firstArc.Add(firstArc, r[1].cast()) // first + arc2
		if firstBytes := firstArc.Bytes(); len(firstBytes) == 0 {
			// We need the explicit zero byte, not an
			// empty []byte{} instance.  This effort
			// is really needed for OID "0.0".
			b = append([]byte{0x00}, b...)
		} else {
			b = append(firstBytes, b...)
		}
		start = 2
	} else {
		if r[0].cast().Uint64() != 2 {
			err = errorf("Only joint-iso-itu-t(2) OIDs allow second-level arcs > 39")
			return
		}

		// Multi-Byte encoding for second-level arcs
		// below joint-iso-itu-t(2), such as "999" for
		// "2.999", that are larger than 39.  Instead,
		// skip the Addition operation we'd normally
		// perform, and just begin VLQ encoding each
		// subsequent byte.
		start = 1
	}

	if len(r) > start {
		for i := start; i < len(r); i++ {
			b = append(b, encodeVLQ(r[i].cast().Bytes())...)
		}
	}

	b = append([]byte{byte(len(b))}, b...) // byte representation of int length of byte slice b
	b = append([]byte{0x06}, b...)         // ASN.1 Object Identifier Tag (0x06)

	return
}

/*
Decode returns an error following an attempt to parse b, which must be
the ASN.1 encoding of an OID, into the receiver instance. The receiver
instance is reinitialized at runtime.
*/
func (r *DotNotation) Decode(b []byte) (err error) {
	if len(b) < 3 {
		err = errorf("Truncated OID encoding")
		return
	}

	if b[0] != 0x06 {
		err = errorf("Invalid ASN.1 Tag; want: 0x06")
		return
	}

	length := int(b[1])
	b = b[2:]

	if length != len(b) {
		err = errorf("Length of bytes does not match with the indicated length")
		return
	}

	var (
		i             int
		subidentifier *big.Int = big.NewInt(0)
	)

	*r = make(DotNotation, 0)

	for i < len(b) {
		for {
			subidentifier.Lsh(subidentifier, 7)
			subidentifier.Add(subidentifier, big.NewInt(int64(b[i]&0x7F)))
			if b[i]&0x80 == 0 {
				break
			}
			i++
		}

		i++
		*r = append(*r, NumberForm(*subidentifier))
		subidentifier = big.NewInt(0)
	}

	if len(*r) > 0 {
		r.decodeFirstArcs(b[0])
	}

	return
}

func (r *DotNotation) decodeFirstArcs(b byte) {
	var firstArc *big.Int
	var secondArc *big.Int

	var forty *big.Int = big.NewInt(40)
	var eighty *big.Int = big.NewInt(80)

	if (*r)[0].cast().Cmp(big.NewInt(80)) < 0 {
		firstArc = big.NewInt(0).Div((*r)[0].cast(), forty)
		secondArc = big.NewInt(0).Mod((*r)[0].cast(), forty)
	} else {
		firstArc = big.NewInt(2)
		if b >= 0x80 {
			// Handle large second-level arcs
			secondArc = big.NewInt(0).Sub((*r)[0].cast(), firstArc)
			secondArc.Add(secondArc, firstArc)
		} else {
			secondArc = big.NewInt(0).Sub((*r)[0].cast(), eighty)
		}
	}

	(*r)[0] = NumberForm(*secondArc)
	*r = append(DotNotation{NumberForm(*firstArc)}, *r...)
}

/*
IntSlice returns slices of integer values and an error. The integer values are based
upon the contents of the receiver. Note that if any single arc number overflows int,
a zero slice is returned.

Successful output can be cast as an instance of [encoding/asn1.ObjectIdentifier], if desired.
*/
func (r DotNotation) IntSlice() (slice []int, err error) {
	if r.IsZero() {
		return
	}

	var t []int
	for i := 0; i < len(r); i++ {
		var n int
		if n, err = atoi(r[i].String()); err != nil {
			return
		}
		t = append(t, n)
	}
	if len(t) > 0 {
		slice = t[:]
	}

	return
}

/*
Uint64Slice returns slices of uint64 values and an error. The uint64
values are based upon the contents of the receiver.

Note that if any single arc number overflows uint64, a zero slice is
returned alongside an error.

Successful output can be cast as an instance of [crypto/x509.OID], if
desired.
*/
func (r DotNotation) Uint64Slice() (slice []uint64, err error) {
	if r.IsZero() {
		return
	}

	var t []uint64
	for i := 0; i < len(r); i++ {
		var n uint64
		if n, err = puint64(r[i].String(), 10, 64); err != nil {
			return
		}
		t = append(t, n)
	}
	if len(t) > 0 {
		slice = t[:]
	}

	return
}

/*
Index returns the Nth index from the receiver, alongside a Boolean
value indicative of success. This method supports the use of negative
indices.
*/
func (r DotNotation) Index(idx int) (a NumberForm, ok bool) {
	if L := len(r); L > 0 {
		if idx < 0 {
			a = r[0]
			if x := L + idx; x >= 0 {
				a = r[x]
			}
		} else if idx > L {
			a = r[L-1]
		} else {
			a = r[idx]
		}
		ok = true
	}

	return
}

/*
Ancestry returns slices of [DotNotation] values ordered from leaf node
(first) to root node (last).

Empty slices of [DotNotation] are returned if the dot notation value
within the receiver is less than two (2) [NumberForm] values in length.
*/
func (r DotNotation) Ancestry() (anc []DotNotation) {
	if r.Len() > 0 {
		for i := r.Len(); i > 0; i-- {
			anc = append(anc, r[:i])
		}
	}

	return
}

/*
AncestorOf returns a Boolean value indicative of whether the receiver
is an ancestor of the input value, which can be string or [DotNotation].
*/
func (r DotNotation) AncestorOf(dot any) (is bool) {
	if !r.IsZero() {
		var D *DotNotation

		switch tv := dot.(type) {
		case string:
			D, _ = NewDotNotation(tv)
		case *DotNotation:
			if tv != nil {
				D = tv
			}
		case DotNotation:
			if tv.Len() >= 0 {
				D = &tv
			}
		}

		if D.Len() > r.Len() {
			is = r.matchDotNot(D)
		}
	}

	return
}

func (r DotNotation) matchDotNot(dot *DotNotation) bool {
	L := r.Len()
	ct := 0
	for i := 0; i < L; i++ {
		x, _ := r.Index(i)
		if y, ok := dot.Index(i); ok {
			if x.Equal(y) {
				ct++
			}
		}
	}

	return ct == L
}

/*
NewSubordinate returns a new instance of [DotNotation] based upon the
contents of the receiver as well as the input [NumberForm] subordinate
value. This creates a fully-qualified child [DotNotation] value of the
receiver.
*/
func (r DotNotation) NewSubordinate(nf any) (dot *DotNotation) {
	if r.Len() > 0 {
		// Prepare the new leaf numberForm, or die trying.
		if a, err := NewNumberForm(nf); err == nil {
			D := make(DotNotation, r.Len()+1, r.Len()+1)
			for i := 0; i < r.Len(); i++ {
				D[i] = r[i]
			}
			D[D.Len()-1] = a
			dot = &D
		}
	}

	return
}

/*
Valid returns a Boolean value indicative of the following:

  - Receiver's length is greater than or equal to two (2) slice members, AND ...
  - The first slice in the receiver contains an unsigned decimal value that is less than three (3)
*/
func (r DotNotation) Valid() (is bool) {
	if !r.IsZero() {
		is = r.Root().Lt(3) && r.Len() >= 2
	}

	return
}

/*
encodeVLQ returns the VLQ -- or Variable Length Quantity -- encoding of
the raw input value.
*/
func encodeVLQ(b []byte) []byte {
	var oid []byte
	n := big.NewInt(0).SetBytes(b)

	for n.Cmp(big.NewInt(0)) > 0 {
		temp := new(big.Int)
		temp.Mod(n, big.NewInt(128))
		if len(oid) > 0 {
			temp.Add(temp, big.NewInt(128))
		}

		oid = append([]byte{byte(temp.Uint64())}, oid...)
		n.Div(n, big.NewInt(128))
	}
	return oid
}

func isNumericOID(id string) bool {
	if !isValidOIDPrefix(id) {
		return false
	}

	var last rune
	for i, c := range id {
		switch {
		case c == '.':
			if last == c {
				return false
			} else if i == len(id)-1 {
				return false
			}
			last = '.'
		case ('0' <= c && c <= '9'):
			last = c
			continue
		}
	}

	return true
}

func isValidOIDPrefix(id string) bool {
	slices := split(id, `.`)
	if len(slices) < 2 {
		return false
	}

	root, err := atoi(slices[0])
	if err != nil {
		return false
	}
	if !(0 <= root && root <= 2) {
		return false
	}

	var sub int
	if sub, err = atoi(slices[1]); err != nil {
		return false
	} else if !(0 <= sub && sub <= 39) && root != 2 {
		return false
	}

	return true
}
