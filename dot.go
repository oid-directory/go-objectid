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
func (d DotNotation) String() (s string) {
	if !d.IsZero() {
		var x []string
		for i := 0; i < len(d); i++ {
			x = append(x, d[i].String())
		}

		s = join(x, `.`)
	}
	return
}

/*
Root returns the root node (0) [NumberForm] instance.
*/
func (d DotNotation) Root() NumberForm {
	x, _ := d.Index(0)
	return x
}

func (d DotNotation) Len() int {
	return len(d)
}

/*
Leaf returns the leaf-node (-1) [NumberForm] instance.
*/
func (d DotNotation) Leaf() NumberForm {
	x, _ := d.Index(-1)
	return x
}

/*
Parent returns the leaf-node's parent (-2) [NumberForm] instance.
*/
func (d DotNotation) Parent() NumberForm {
	x, _ := d.Index(-2)
	return x
}

/*
IsZero returns a Boolean indicative of whether the receiver is unset.
*/
func (d *DotNotation) IsZero() (is bool) {
	if d != nil {
		is = d.Len() == 0
	}
	return
}

/*
NewDotNotation returns an instance of *[DotNotation] alongside a Boolean
value indicative of success.

Variadic input allows for slice mixtures of all of the following types,
each treated as an individual [NumberForm] instance:

  - *[math/big.Int]
  - string
  - uint64
  - uint
  - int

If a string primitive is the only input option, it will be treated as a
complete [DotNotation] (e.g.: "1.3.6").
*/
func NewDotNotation(x ...any) (d *DotNotation, err error) {
	var _d DotNotation = make(DotNotation, 0)

	if len(x) == 1 {
		if slice, ok := x[0].(string); ok {
			d, err = newDotNotationStr(slice)
			return
		}
	}

	for i := 0; i < len(x) && err == nil; i++ {
		var nf NumberForm
		switch tv := x[i].(type) {
		case *big.Int, string, uint64, uint, int:
			nf, err = NewNumberForm(tv)
		default:
			err = errorf("Unsupported slice type '%T' for OID", tv)
		}

		if err == nil {
			_d = append(_d, nf)
		}
	}

	if err == nil {
		d = new(DotNotation)
		*d = _d
	}

	return
}

func newDotNotationStr(dot string) (d *DotNotation, err error) {
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
		d = new(DotNotation)
		*d = _d
	}

	return
}

/*
Encode returns the ASN.1 encoding of the receiver instance alongside an error.
*/
func (d DotNotation) Encode() (b []byte, err error) {
	if d.Len() < 2 {
		err = errorf("Length below encoding minimum")
		return
	}

	var start int
	firstArc := d[0].cast()
	forty := big.NewInt(40)
	firstArc.Mul(firstArc, forty)

	if d[1].cast().Cmp(forty) < 1 {
		// this is meant for second level arcs <= 39
		firstArc.Add(firstArc, d[1].cast()) // first + arc2
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
		if d[0].cast().Uint64() != 2 {
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

	if len(d) > start {
		for i := start; i < len(d); i++ {
			b = append(b, encodeVLQ(d[i].cast().Bytes())...)
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
func (d *DotNotation) Decode(b []byte) (err error) {
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

	*d = make(DotNotation, 0)

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
		*d = append(*d, NumberForm(*subidentifier))
		subidentifier = big.NewInt(0)
	}

	if len(*d) > 0 {
		d.decodeFirstArcs(b[0])
	}

	return
}

func (d *DotNotation) decodeFirstArcs(b byte) {
	var firstArc *big.Int
	var secondArc *big.Int

	var forty *big.Int = big.NewInt(40)
	var eighty *big.Int = big.NewInt(80)

	if (*d)[0].cast().Cmp(big.NewInt(80)) < 0 {
		firstArc = big.NewInt(0).Div((*d)[0].cast(), forty)
		secondArc = big.NewInt(0).Mod((*d)[0].cast(), forty)
	} else {
		firstArc = big.NewInt(2)
		if b >= 0x80 {
			// Handle large second-level arcs
			secondArc = big.NewInt(0).Sub((*d)[0].cast(), firstArc)
			secondArc.Add(secondArc, firstArc)
		} else {
			secondArc = big.NewInt(0).Sub((*d)[0].cast(), eighty)
		}
	}

	(*d)[0] = NumberForm(*secondArc)
	*d = append(DotNotation{NumberForm(*firstArc)}, *d...)
}

/*
IntSlice returns slices of integer values and an error. The integer values are based
upon the contents of the receiver. Note that if any single arc number overflows int,
a zero slice is returned.

Successful output can be cast as an instance of [encoding/asn1.ObjectIdentifier], if desired.
*/
func (d DotNotation) IntSlice() (slice []int, err error) {
	if d.IsZero() {
		return
	}

	var t []int
	for i := 0; i < len(d); i++ {
		var n int
		if n, err = atoi(d[i].String()); err != nil {
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
func (d DotNotation) Uint64Slice() (slice []uint64, err error) {
	if d.IsZero() {
		return
	}

	var t []uint64
	for i := 0; i < len(d); i++ {
		var n uint64
		if n, err = puint64(d[i].String(), 10, 64); err != nil {
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
func (d DotNotation) Index(idx int) (a NumberForm, ok bool) {
	if L := len(d); L > 0 {
		if idx < 0 {
			a = d[0]
			if x := L + idx; x >= 0 {
				a = d[x]
			}
		} else if idx > L {
			a = d[L-1]
		} else {
			a = d[idx]
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
func (d DotNotation) Ancestry() (anc []DotNotation) {
	if d.Len() > 0 {
		for i := d.Len(); i > 0; i-- {
			anc = append(anc, d[:i])
		}
	}

	return
}

/*
AncestorOf returns a Boolean value indicative of whether the receiver
is an ancestor of the input value, which can be string or [DotNotation].
*/
func (d DotNotation) AncestorOf(dot any) (is bool) {
	if !d.IsZero() {
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

		if D.Len() > d.Len() {
			is = d.matchDotNot(D)
		}
	}

	return
}

func (d DotNotation) matchDotNot(dot *DotNotation) bool {
	L := d.Len()
	ct := 0
	for i := 0; i < L; i++ {
		x, _ := d.Index(i)
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
func (d DotNotation) NewSubordinate(nf any) (dot *DotNotation) {
	if d.Len() > 0 {
		// Prepare the new leaf numberForm, or die trying.
		if a, err := NewNumberForm(nf); err == nil {
			D := make(DotNotation, d.Len()+1, d.Len()+1)
			for i := 0; i < d.Len(); i++ {
				D[i] = d[i]
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
func (d DotNotation) Valid() (is bool) {
	if !d.IsZero() {
		is = d.Root().Lt(3) && d.Len() >= 2
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
