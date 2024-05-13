package objectid

/*
OID contains an underlying [ASN1Notation] value, and extends convenient methods allowing
interrogation and verification.
*/
type OID struct {
	nanf   ASN1Notation
	parsed bool
}

/*
IsZero checks the receiver for nilness and returns a Boolean indicative of the result.
*/
func (r *OID) IsZero() (is bool) {
	if r != nil {
		is = len(r.nanf) == 0
	}
	return
}

/*
Dot returns a [DotNotation] instance based on the contents of the underlying [ASN1Notation]
instance found within the receiver.

Note that at a receiver length of two (2) or more is required for successful output.
*/
func (r OID) Dot() (d DotNotation) {
	if r.Len() < 2 {
		return
	}
	if !r.IsZero() {
		d = make(DotNotation, len(r.nanf))
		for i := 0; i < len(r.nanf); i++ {
			d[i] = r.nanf[i].NumberForm()
		}
	}

	return
}

/*
ASN returns the underlying [ASN1Notation] instance found within the receiver.
*/
func (r OID) ASN() (a ASN1Notation) {
	if !r.IsZero() {
		a = r.nanf
	}
	return
}

/*
Valid returns a Boolean value indicative of whether the receiver's state is considered value.
*/
func (r OID) Valid() (ok bool) {
	if !r.IsZero() {
		var nanf NameAndNumberForm
		if nanf, ok = r.nanf.Index(0); ok {
			var found bool
			for i := 0; i < 3; i++ {
				if nanf.NumberForm().Equal(i) {
					found = true
					break
				}
			}
			ok = found
		}
	}
	return
}

/*
Len returns the integer length of all underlying [NumberForm] values present within the receiver.
*/
func (r OID) Len() (i int) {
	if !r.IsZero() {
		i = len(r.nanf)
	}

	return
}

/*
Leaf returns the leaf-node instance of [NameAndNumberForm].
*/
func (r OID) Leaf() (nanf NameAndNumberForm) {
	if !r.IsZero() {
		nanf, _ = r.nanf.Index(-1)
	}
	return
}

/*
Parent returns the leaf-node's Parent instance of [NameAndNumberForm].
*/
func (r OID) Parent() (nanf NameAndNumberForm) {
	if !r.IsZero() {
		nanf, _ = r.nanf.Index(-2)
	}
	return
}

/*
Root returns the root node instance of [NameAndNumberForm].
*/
func (r OID) Root() (nanf NameAndNumberForm) {
	if !r.IsZero() {
		nanf, _ = r.nanf.Index(0)
	}
	return
}

/*
NewOID creates an instance of [OID] and returns it alongside an error.

Valid input forms for ASN.1 values are:

  - string (e.g.: "{iso(1) ... }")
  - string slices (e.g.: []string{"iso(1)", "identified-organization(3)" ...})
  - [NameAndNumberForm] slices ([][NameAndNumberForm]{...})

Not all [NameAndNumberForm] values (arcs) require actual names; they can be
numbers alone or in the so-called nameAndNumber syntax (name(Number)). For example:

	{iso(1) identified-organization(3) 6}

... is perfectly valid, but generally NOT recommended when clarity or precision is desired.

Note that the following root node abbreviations are supported:

  - `itu-t` resolves to itu-t(0)
  - `iso` resolves to iso(1)
  - `joint-iso-itu-t` resolves to joint-iso-itu-t(2)

Case is significant during processing of the above abbreviations.  Note that it is
inappropriate to utilize these abbreviations for any portion of an [OID] instance
other than as the respective root node.

[NumberForm] values CANNOT be negative, but are unbounded in their magnitude.
*/
func NewOID(x any) (r *OID, err error) {
	// prepare temporary instance
	t := new(OID)
	r = new(OID)

	var nfs []string
	switch tv := x.(type) {
	case []NameAndNumberForm:
		t.nanf = ASN1Notation(tv)
		if !t.Valid() {
			err = errorf("%T instance did not pass validity checks: %#v", t, t)
			break
		}
		r.nanf = t.nanf
		r.parsed = true
		return
	case string:
		nfs = fields(condenseWHSP(trimR(trimL(tv, `{`), `}`)))
	case []string:
		nfs = tv
	default:
		err = errorf("Unsupported %T input type: %#v\n", x, x)
		return
	}

	for i := 0; i < len(nfs); i++ {
		var nanf *NameAndNumberForm
		if nanf, err = NewNameAndNumberForm(nfs[i]); err != nil {
			break
		}
		t.nanf = append(t.nanf, *nanf)
	}

	if err == nil {
		if !t.Valid() {
			err = errorf("%T instance did not pass validity checks: %#v", t, *t)
			return
		}

		r.parsed = true
		r.nanf = t.nanf
	}

	return
}
