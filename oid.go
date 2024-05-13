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
func (id *OID) IsZero() (is bool) {
	if id != nil {
		is = len(id.nanf) == 0
	}
	return
}

/*
Dot returns a [DotNotation] instance based on the contents of the underlying [ASN1Notation]
instance found within the receiver.

Note that at a receiver length of two (2) or more is required for successful output.
*/
func (id OID) Dot() (d DotNotation) {
	if id.Len() < 2 {
		return
	}
	if !id.IsZero() {
		d = make(DotNotation, len(id.nanf))
		for i := 0; i < len(id.nanf); i++ {
			d[i] = id.nanf[i].NumberForm()
		}
	}

	return
}

/*
ASN returns the underlying [ASN1Notation] instance found within the receiver.
*/
func (id OID) ASN() (a ASN1Notation) {
	if !id.IsZero() {
		a = id.nanf
	}
	return
}

/*
Valid returns a Boolean value indicative of whether the receiver's state is considered value.
*/
func (id OID) Valid() (ok bool) {
	if !id.IsZero() {
		var nanf NameAndNumberForm
		if nanf, ok = id.nanf.Index(0); ok {
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
func (id OID) Len() (i int) {
	if !id.IsZero() {
		i = len(id.nanf)
	}

	return
}

/*
Leaf returns the leaf-node instance of [NameAndNumberForm].
*/
func (id OID) Leaf() (nanf NameAndNumberForm) {
	if !id.IsZero() {
		nanf, _ = id.nanf.Index(-1)
	}
	return
}

/*
Parent returns the leaf-node's Parent instance of [NameAndNumberForm].
*/
func (id OID) Parent() (nanf NameAndNumberForm) {
	if !id.IsZero() {
		nanf, _ = id.nanf.Index(-2)
	}
	return
}

/*
Root returns the root node instance of [NameAndNumberForm].
*/
func (id OID) Root() (nanf NameAndNumberForm) {
	if !id.IsZero() {
		nanf, _ = id.nanf.Index(0)
	}
	return
}

/*
NewOID creates an instance of [OID] and returns it alongside an error.

The correct raw input syntax is the ASN.1 [NameAndNumberForm] sequence syntax, i.e.:

	{iso(1) identified-organization(3) dod(6)}

Not all [NameAndNumberForm] values (arcs) require actual names; they can be numbers alone or in the so-called nameAndNumber syntax (name(Number)). For example:

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
func NewOID(x any) (o *OID, err error) {
	t := new(OID)

	var nfs []string
	switch tv := x.(type) {
	case string:
		nfs = fields(condenseWHSP(trimR(trimL(tv, `{`), `}`)))
	case []string:
		nfs = tv
	default:
		err = errorf("Unsupported %T input type: %#v\n", x, x)
		return
	}

	for i := 0; i < len(nfs) && err == nil; i++ {
		var nanf *NameAndNumberForm
		if nanf, err = NewNameAndNumberForm(nfs[i]); nanf != nil {
			t.nanf = append(t.nanf, *nanf)
		}
	}

	if err == nil {
		if !t.Valid() {
			err = errorf("%T instance did not pass validity checks: %#v", t, *t)
			return
		}

		o = new(OID)
		o.parsed = true
		*o = *t
	}

	return
}
