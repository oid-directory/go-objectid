package objectid

/*
asn.go handles ASN1Notation operations. For ASN.1 encoding/decoding, see codec.go.
*/

/*
ASN1Notation contains an ordered sequence of [NameAndNumberForm] instances.
*/
type ASN1Notation []NameAndNumberForm

/*
String is a stringer method that returns a properly formatted ASN.1 string value.
*/
func (a ASN1Notation) String() string {
	var x []string
	for i := 0; i < len(a); i++ {
		x = append(x, a[i].String())
	}
	return `{` + join(x, ` `) + `}`
}

/*
Dot returns a [DotNotation] instance based on the contents of the receiver instance.

Note that at a receiver length of two (2) or more is required for successful output.
*/
func (a ASN1Notation) Dot() (d DotNotation) {
	if a.Len() < 2 {
		return
	}
	if !a.IsZero() {
		L := a.Len()
		d = make(DotNotation, L)
		for i := 0; i < L; i++ {
			d[i] = a[i].NumberForm()
		}
	}

	return
}

/*
Root returns the root node (0) string value from the receiver.
*/
func (a ASN1Notation) Root() NameAndNumberForm {
	x, _ := a.Index(0)
	return x
}

/*
Leaf returns the leaf node (-1) string value from the receiver.
*/
func (a ASN1Notation) Leaf() NameAndNumberForm {
	x, _ := a.Index(-1)
	return x
}

/*
Parent returns the leaf node's parent (-2) string value from the receiver.
*/
func (a ASN1Notation) Parent() NameAndNumberForm {
	x, _ := a.Index(-2)
	return x
}

/*
Len returns the integer length of the receiver.
*/
func (a ASN1Notation) Len() int { return len(a) }

/*
IsZero returns a Boolean indicative of whether the receiver is unset.
*/
func (a ASN1Notation) IsZero() (is bool) {
	if is = &a == nil; !is {
		is = a.Len() == 0
	}

	return
}

/*
Index returns the Nth index from the receiver, alongside a Boolean
value indicative of success. This method supports the use of negative
indices.
*/
func (a ASN1Notation) Index(idx int) (nanf NameAndNumberForm, ok bool) {
	L := a.Len()

	// Bail if receiver is empty.
	if L > 0 {
		if idx < 0 {
			var x int = L + idx
			if x < 0 {
				nanf = a[0]
			} else {
				nanf = a[x]
			}
		} else if idx > L {
			nanf = a[L-1]
		} else {
			nanf = a[idx]
		}
	}

	// Make sure the instance was produced
	// via recommended procedure.
	ok = nanf.parsed

	return
}

/*
NewASN1Notation returns an instance of *[ASN1Notation] alongside an error.

Valid input forms for ASN.1 values are string (e.g.: "{iso(1)}") and string
slices (e.g.: []string{"iso(1)", "identified-organization(3)" ...}).

[NumberForm] values CANNOT be negative.
*/
func NewASN1Notation(x any) (a *ASN1Notation, err error) {
	// prepare temporary instance
	t := new(ASN1Notation)

	var nfs []string
	switch tv := x.(type) {
	case string:
		nfs = fields(condenseWHSP(trimR(trimL(tv, `{`), `}`)))
	case []string:
		nfs = tv
	default:
		err = errorf("Unsupported %T input type: %#v", x, x)
		return
	}

	for i := 0; i < len(nfs) && err == nil; i++ {
		var nanf *NameAndNumberForm
		if nanf, err = NewNameAndNumberForm(nfs[i]); nanf != nil {
			*t = append(*t, *nanf)
		}
	}

	if err != nil {
		return
	}

	// verify content is valid
	err = errorf("%T instance did not pass validity checks: %#v", t, *t)
	if t.Valid() {
		// transfer temporary content
		// to return value instance.
		a = new(ASN1Notation)
		*a = *t
		err = nil
	}

	return
}

/*
Valid returns a Boolean value indicative of whether the receiver's
length is greater than or equal to one (1) slice member.
*/
func (a ASN1Notation) Valid() (is bool) {
	// Don't waste time on
	// zero instances.
	if L := a.Len(); L > 0 {
		if root, ok := a.Index(0); ok {
			// root cannot be greater than 2
			is = root.NumberForm().Lt(3)
		}
	}

	return
}

/*
Ancestry returns slices of [DotNotation] values ordered from leaf node
(first) to root node (last).

Empty slices of DotNotation are returned if the dotNotation value
within the receiver is less than two (2) [NumberForm] values in length.
*/
func (a ASN1Notation) Ancestry() (anc []ASN1Notation) {
	if a.Len() >= 2 {
		for i := a.Len(); i > 0; i-- {
			anc = append(anc, a[:i])
		}
	}

	return
}

/*
NewSubordinate returns a new instance of [ASN1Notation] based upon the
contents of the receiver as well as the input [NameAndNumberForm]
subordinate value. This creates a fully-qualified child [ASN1Notation]
value of the receiver.
*/
func (a ASN1Notation) NewSubordinate(nanf any) *ASN1Notation {
	var A ASN1Notation
	if a.Len() > 0 {
		// Prepare the new leaf numberForm, or die trying.
		if n, err := NewNameAndNumberForm(nanf); err == nil {
			A = make(ASN1Notation, a.Len()+1, a.Len()+1)
			for i := 0; i < a.Len(); i++ {
				A[i] = a[i]
			}
			A[A.Len()-1] = *n
		}
	}

	return &A
}

/*
AncestorOf returns a Boolean value indicative of whether the receiver
is an ancestor of the input value, which can be string or [ASN1Notation].
*/
func (a ASN1Notation) AncestorOf(asn any) (anc bool) {
	if !a.IsZero() {
		var A *ASN1Notation

		switch tv := asn.(type) {
		case string:
			A, _ = NewASN1Notation(tv)
		case *ASN1Notation:
			if tv != nil {
				A = tv
			}
		case ASN1Notation:
			if tv.Len() >= 0 {
				A = &tv
			}
		}

		if A.Len() > a.Len() {
			anc = a.matchASN1(A)
		}
	}

	return
}

func (a ASN1Notation) matchASN1(asn *ASN1Notation) (matched bool) {
	L := a.Len()
	ct := 0
	for i := 0; i < L; i++ {
		x, _ := a.Index(i)
		if y, ok := asn.Index(i); ok {
			if x.Equal(y) {
				ct++
			}
		}
	}

	return ct == L
}
