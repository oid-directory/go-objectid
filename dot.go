package objectid

/*
DotNotation contains an ordered sequence of NumberForm instances.
*/
type DotNotation []NumberForm

/*
String is a stringer method that returns the dotNotation
form of the receiver (e.g.: "1.3.6.1").
*/
func (d DotNotation) String() string {
	var x []string
	for i := 0; i < len(d); i++ {
		x = append(x, d[i].String())
	}
	return join(x, `.`)
}

/*
Root returns the root node (0) NumberForm instance.
*/
func (d DotNotation) Root() NumberForm {
	x, _ := d.Index(0)
	return x
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
NewDotNotation returns an instance of *DotNotation alongside a boolean value
indicative of success.
*/
func NewDotNotation(id string) (d *DotNotation, err error) {
	if len(id) == 0 {
		return
	}

	ids := split(id, `.`)
	var t *DotNotation = new(DotNotation)
	for i := 0; i < len(ids); i++ {
		var a NumberForm
		if a, err = ParseNumberForm(ids[i]); err != nil {
			break
		}
		*t = append(*t, a)
	}

	if err == nil {
		d = new(DotNotation)
		*d = *t
	}

	return
}

/*
IntSlice returns slices of integer values based upon the contents of the receiver
*/
func (d DotNotation) IntSlice() (slice []int, err error) {
	if len(d) == 0 {
		err = errorf("%T instance is nil", d)
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
Index returns the nth index from the receiver, alongside a boolean value
indicative of success. This method supports the use of negative indices.
*/
func (d DotNotation) Index(idx int) (a NumberForm, ok bool) {
	L := len(d)

	// Bail if receiver is empty.
	if L == 0 {
		return
	}

	if idx < 0 {
		var x int = L + idx
		if x < 0 {
			a = d[0]
		} else {
			a = d[x]
		}
	} else if idx > L {
		a = d[L-1]
	} else {
		a = d[idx]
	}

	ok = true
	return
}
