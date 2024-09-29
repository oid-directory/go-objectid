package objectid

/*
asn.go handles ASN1Notation operations. For object
identifier encoding/decoding, see dot.go.
*/

/*
ASN1Notation contains an ordered sequence of [NameAndNumberForm] instances.
*/
type ASN1Notation []NameAndNumberForm

/*
String is a stringer method that returns a properly formatted ASN.1 string value.
*/
func (r ASN1Notation) String() string {
	var x []string
	for i := 0; i < len(r); i++ {
		x = append(x, r[i].String())
	}
	return `{` + join(x, ` `) + `}`
}

/*
Dot returns a [DotNotation] instance based on the contents of the receiver instance.

Note that at a receiver length of two (2) or more is required for successful output.
*/
func (r ASN1Notation) Dot() (d DotNotation) {
	if r.Len() < 2 {
		return
	}
	if !r.IsZero() {
		L := r.Len()
		d = make(DotNotation, L)
		for i := 0; i < L; i++ {
			d[i] = r[i].NumberForm()
		}
	}

	return
}

/*
Root returns the root node (0) string value from the receiver.
*/
func (r ASN1Notation) Root() NameAndNumberForm {
	x, _ := r.Index(0)
	return x
}

/*
Leaf returns the leaf node (-1) string value from the receiver.
*/
func (r ASN1Notation) Leaf() NameAndNumberForm {
	x, _ := r.Index(-1)
	return x
}

/*
Parent returns the leaf node's parent (-2) string value from the receiver.
*/
func (r ASN1Notation) Parent() NameAndNumberForm {
	x, _ := r.Index(-2)
	return x
}

/*
Len returns the integer length of the receiver.
*/
func (r ASN1Notation) Len() int { return len(r) }

/*
IsZero returns a Boolean indicative of whether the receiver is unset.
*/
func (r ASN1Notation) IsZero() (is bool) {
	if is = &r == nil; !is {
		is = r.Len() == 0
	}

	return
}

/*
Index returns the Nth index from the receiver, alongside a Boolean
value indicative of success. This method supports the use of negative
indices.
*/
func (r ASN1Notation) Index(idx int) (nanf NameAndNumberForm, ok bool) {
	L := r.Len()

	// Bail if receiver is empty.
	if L > 0 {
		if idx < 0 {
			var x int = L + idx
			if x < 0 {
				nanf = r[0]
			} else {
				nanf = r[x]
			}
		} else if idx > L {
			nanf = r[L-1]
		} else if idx < L {
			nanf = r[idx]
		}
	}

	// Make sure the instance was produced
	// via recommended procedure.
	ok = nanf.parsed

	return
}

/*
NewASN1Notation returns an instance of *[ASN1Notation] alongside an error.

Valid input forms for ASN.1 values are:

  - string (e.g.: "{iso(1) ... }")
  - string slices (e.g.: []string{"iso(1)", "identified-organization(3)" ...})
  - [NameAndNumberForm] slices ([][NameAndNumberForm]{...})

Note that the following root node abbreviations are supported:

  - `itu-t` resolves to itu-t(0)
  - `iso` resolves to iso(1)
  - `joint-iso-itu-t` resolves to joint-iso-itu-t(2)

Case is significant during processing of the above abbreviations.  Note that it is
inappropriate to utilize these abbreviations for any portion of an [ASN1Notation]
instance other than as the respective root node.

[NumberForm] values CANNOT be negative, but are unbounded in their magnitude.
*/
func NewASN1Notation(x any) (r *ASN1Notation, err error) {
	// prepare temporary instance
	t := make(ASN1Notation, 0)
	r = new(ASN1Notation)

	var nfs []string
	switch tv := x.(type) {
	case []NameAndNumberForm:
		t = ASN1Notation(tv)
		if !t.Valid() {
			err = errorf("%T instance did not pass validity checks: %#v", t, t)
			break
		}
		*r = t
		return
	case string:
		nfs = fields(condenseWHSP(trimR(trimL(tv, `{`), `}`)))
	case []string:
		nfs = tv
	default:
		err = errorf("Unsupported %T input type: %#v", x, x)
		return
	}

	for i := 0; i < len(nfs); i++ {
		var nanf *NameAndNumberForm
		if nanf, err = NewNameAndNumberForm(nfs[i]); err != nil {
			break
		}
		t = append(t, *nanf)
	}

	if err == nil {
		// verify content is valid
		if !t.Valid() {
			err = errorf("%T instance did not pass validity checks: %#v", t, t)
			return
		}

		// transfer temporary content
		// to return value instance.
		*r = t
	}

	return
}

/*
Valid returns a Boolean value indicative of whether the receiver's
length is greater than or equal to one (1) slice member.
*/
func (r ASN1Notation) Valid() (is bool) {
	// Don't waste time on
	// zero instances.
	if L := r.Len(); L > 0 {
		if root, ok := r.Index(0); ok {
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
func (r ASN1Notation) Ancestry() (anc []ASN1Notation) {
	if r.Len() >= 2 {
		for i := r.Len(); i > 0; i-- {
			anc = append(anc, r[:i])
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
func (r ASN1Notation) NewSubordinate(nanf any) *ASN1Notation {
	var A ASN1Notation
	if r.Len() > 0 {
		// Prepare the new leaf numberForm, or die trying.
		if n, err := NewNameAndNumberForm(nanf); err == nil {
			A = make(ASN1Notation, r.Len()+1, r.Len()+1)
			for i := 0; i < r.Len(); i++ {
				A[i] = r[i]
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
func (r ASN1Notation) AncestorOf(asn any) (anc bool) {
	if !r.IsZero() {
		if A := assertASN1Notation(asn); !A.IsZero() {
			if A.Len() > r.Len() {
				anc = r.matchASN1(A, 0)
			}
		}
	}

	return
}

/*
ChildOf returns a Boolean value indicative of whether the receiver is
a direct superior (parent) of the input value, which can be string or
[ASN1Notation].
*/
func (r ASN1Notation) ChildOf(asn any) (cof bool) {
	if !r.IsZero() {
		if A := assertASN1Notation(asn); !A.IsZero() {
			if A.Len()-1 == r.Len() {
				cof = r.matchASN1(A, 0)
			}
		}
	}

	return
}

/*
SiblingOf returns a Boolean value indicative of whether the receiver is
a sibling of the input value, which can be string or [ASN1Notation].
*/
func (r ASN1Notation) SiblingOf(asn any) (sof bool) {
	if !r.IsZero() {
		if A := assertASN1Notation(asn); !A.IsZero() {
			if A.Len() == r.Len() && !A.Leaf().Equal(r.Leaf()) {
				sof = r.matchASN1(A, -1)
			}
		}
	}

	return
}

func (r ASN1Notation) matchASN1(asn *ASN1Notation, off int) (matched bool) {
	L := r.Len()
	ct := 0
	for i := 0; i < L; i++ {
		x, _ := r.Index(i)
		if y, ok := asn.Index(i); ok {
			if x.Equal(y) {
				ct++
			} else if off == -1 && L-1 == i {
				// sibling check should end in
				// a FAILED match for the final
				// arcs.
				ct++
			}
		}
	}

	return ct == L
}

func assertASN1Notation(asn any) (A *ASN1Notation) {
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

	return
}
