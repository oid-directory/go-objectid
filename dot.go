package objectid

/*
DotNotation contains an ordered sequence of NumberForm instances.
*/
type DotNotation []NumberForm

/*
String is a stringer method that returns the dotNotation
form of the receiver (e.g.: "1.3.6.1").
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
Root returns the root node (0) NumberForm instance.
*/
func (d DotNotation) Root() NumberForm {
	x, _ := d.Index(0)
	return x
}

func (d DotNotation) Len() int {
	return len(d)
}

/*
Leaf returns the leaf-node (-1) NumberForm instance.
*/
func (d DotNotation) Leaf() NumberForm {
	x, _ := d.Index(-1)
	return x
}

/*
Parent returns the leaf-node's parent (-2) NumberForm instance.
*/
func (d DotNotation) Parent() NumberForm {
	x, _ := d.Index(-2)
	return x
}

/*
IsZero returns a boolean indicative of whether the receiver
is unset.
*/
func (d *DotNotation) IsZero() (is bool) {
	if d != nil {
		is = d.Len() == 0
	}
	return
}

/*
NewDotNotation returns an instance of *DotNotation alongside a boolean
value indicative of success.
*/
func NewDotNotation(id string) (d *DotNotation, err error) {
	if len(id) == 0 {
		err = errorf("input value is zero-length")
		return
	}

	ids := split(id, `.`)
	var t *DotNotation = new(DotNotation)
	for i := 0; i < len(ids) && err == nil; i++ {
		var a NumberForm
		a, err = NewNumberForm(ids[i])
		*t = append(*t, a)
	}

	if err == nil {
		d = new(DotNotation)
		*d = *t
	}

	return
}

/*
IntSlice returns slices of integer values and an error. The integer values are based
upon the contents of the receiver. Note that if any single arc number overflows int,
a zero slice is returned.

Successful output can be cast as an instance of asn1.ObjectIdentifier, if desired.
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
Index returns the Nth index from the receiver, alongside a boolean
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
Ancestry returns slices of DotNotation values ordered from leaf
node (first) to root node (last).

Empty slices of DotNotation are returned if the dotNotation value
within the receiver is less than two (2) NumberForm values in length.
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
AncestorOf returns a boolean value indicative of whether the receiver
is an ancestor of the input value, which can be string or DotNotation.
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
NewSubordinate returns a new instance of DotNotation based upon the
contents of the receiver as well as the input NumberForm subordinate
value. This creates a fully-qualified child DotNotation value of the
receiver.
*/
func (d DotNotation) NewSubordinate(nf any) (dot *DotNotation) {
	if d.Len() > 0 {
		// Prepare the new leaf numberForm,
		// or die trying.
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
Valid returns a boolean value indicative of the following:

• Receiver's length is greater than or equal to one (1) slice member, and ...
• The first slice in the receiver contains a decimal value that is less than three (3)
*/
func (d DotNotation) Valid() (is bool) {
	if !d.IsZero() {
		is = d.Root().Lt(3)
	}

	return
}
