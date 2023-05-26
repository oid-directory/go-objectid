package objectid

/*
ASN1Notation contains an ordered sequence of NameAndNumberForm
instances.
*/
type ASN1Notation []NameAndNumberForm

/*
String is a stringer method that returns a properly formatted
ASN.1 string value.
*/
func (a ASN1Notation) String() string {
	var x []string
	for i := 0; i < len(a); i++ {
		x = append(x, a[i].String())
	}
	return `{` + join(x, ` `) + `}`
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
Parent returns the leaf node's parent (-2) string value from
the receiver.
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
IsZero returns a boolean indicative of whether the receiver
is unset.
*/
func (a ASN1Notation) IsZero() bool {
	if &a == nil {
		return true
	}
	return a.Len() == 0
}

/*
Index returns the Nth index from the receiver, alongside a boolean
value indicative of success. This method supports the use of negative
indices.
*/
func (a ASN1Notation) Index(idx int) (nanf NameAndNumberForm, ok bool) {
	L := a.Len()

	// Bail if receiver is empty.
	if L == 0 {
		return
	}

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

	// Make sure the instance was produced
	// via recommended procedure.
	ok = nanf.parsed

	return
}

/*
NewASN1Notation returns an instance of *ASN1Notation alongside an error.

Valid input forms for ASN.1 values are string (e.g.: "{iso(1)}") and string
slices (e.g.: []string{"iso(1)", "identified-organization(3)" ...}).

NumberForm values CANNOT be negative, and CANNOT overflow NumberForm (uint128).
*/
func NewASN1Notation(x any) (a *ASN1Notation, err error) {
	// prepare temporary instance
	t := new(ASN1Notation)

	switch tv := x.(type) {
	case string:
		f := fields(condenseWHSP(trimR(trimL(tv, `{`), `}`)))
		for i := 0; i < len(f); i++ {
			var nanf *NameAndNumberForm
			if nanf, err = NewNameAndNumberForm(f[i]); err != nil {
				return
			}
			*t = append(*t, *nanf)
		}
	case []string:
		for i := 0; i < len(tv); i++ {
			var nanf *NameAndNumberForm
			if nanf, err = NewNameAndNumberForm(condenseWHSP(tv[i])); err != nil {
				return
			}
			*t = append(*t, *nanf)
		}
	default:
		err = errorf("Unsupported %T input type: %#v\n", x, x)
		return
	}

	// verify content is valid
	if !t.Valid() {
		err = errorf("%T instance did not pass validity checks: %#v", t, *t)
		return
	}

	// transfer temporary content
	// to return value instance.
	a = new(ASN1Notation)
	*a = *t

	return
}

/*
Valid returns a boolean value indicative of whether the receiver's
length is greater than or equal to one (1) slice member.
*/
func (a ASN1Notation) Valid() bool {
	// Don't waste time on
	// zero instances.
	if a.Len() == 0 {
		return false
	}

	// bail out if any of the slice
	// values are unparsed.
	for i := 0; i < a.Len(); i++ {
		if !a[i].parsed {
			return false
		}
	}

	root, ok := a.Index(0)
	if !ok {
		return false
	}

	// root cannot be greater than 2
	return root.NumberForm().Lt(3)
}

/*
Ancestry returns slices of DotNotation values ordered from leaf node
(first) to root node (last).

Empty slices of DotNotation are returned if the dotNotation value
within the receiver is less than two (2) NumberForm values in length.
*/
func (a ASN1Notation) Ancestry() (anc []ASN1Notation) {
	if a.Len() < 2 {
		return
	}

	for i := a.Len(); i > 0; i-- {
		anc = append(anc, a[:i])
	}

	return
}

/*
NewSubordinate returns a new instance of ASN1Notation based upon the
contents of the receiver as well as the input NameAndNumberForm
subordinate value. This creates a fully-qualified child ASN1Notation
value of the receiver.
*/
func (a ASN1Notation) NewSubordinate(nanf any) *ASN1Notation {
	// Don't bother processing
	if a.Len() == 0 {
		return nil
	}

	// Prepare the new leaf numberForm,
	// or die trying.
	n, err := NewNameAndNumberForm(nanf)
	if err != nil {
		return nil
	}

	A := make(ASN1Notation, a.Len()+1, a.Len()+1)
	for i := 0; i < a.Len(); i++ {
		A[i] = a[i]
	}
	A[A.Len()-1] = *n

	return &A
}

/*
AncestorOf returns a boolean value indicative of whether the receiver
is an ancestor of the input value, which can be string or ASN1Notation.
*/
func (a ASN1Notation) AncestorOf(asn any) bool {
	if a.IsZero() {
		return false
	}

	var A *ASN1Notation

	switch tv := asn.(type) {
	case string:
		var err error
		if A, err = NewASN1Notation(tv); err != nil {
			return false
		}
	case *ASN1Notation:
		if tv == nil {
			return false
		}
		A = tv
	case ASN1Notation:
		if tv.Len() == 0 {
			return false
		}
		A = &tv
	default:
		return false
	}

	if A.Len() < a.Len() {
		return false
	}

	for i := 0; i < a.Len(); i++ {
		x, _ := a.Index(i)
		y, ok := A.Index(i)
		if !ok {
			return false
		}
		if !x.Equal(y) {
			return false
		}
	}

	return true
}
