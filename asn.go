package objectid

/*
ASN1Notation contains an ordered sequence of NameAndNumberForm instances.
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
Index returns the nth index from the receiver, alongside a boolean value
indicative of success. This method supports the use of negative indices.
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
Valid returns a boolean value indicative of whether the receiver's length is greater than or equal to one (1) slice member.
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
	return root.NumberForm().Compare64(uint64(2)) != +1
}
