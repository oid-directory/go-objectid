package objectid

import "math/big"

/*
nanf.go deals with NameAndNumberForm syntax and viability
*/

/*
NameAndNumberForm contains either an identifier with a parenthesis-enclosed
decimal value, or a decimal value alone. An ordered sequence of instances of
this type comprise an instance of [ASN1Notation].
*/
type NameAndNumberForm struct {
	identifier        string
	primaryIdentifier NumberForm
	parsed            bool
}

/*
IsZero returns a Boolean valu indicative of whether
the receiver is considered nil.
*/
func (r NameAndNumberForm) IsZero() bool {
	return !r.parsed
}

/*
Identifier returns the string-based nameForm
value assigned to the receiver instance.
*/
func (r NameAndNumberForm) Identifier() string {
	return r.identifier
}

/*
NumberForm returns the underlying [NumberForm]
value assigned to the receiver instance.
*/
func (r NameAndNumberForm) NumberForm() NumberForm {
	return r.primaryIdentifier
}

/*
String is a stringer method that returns the properly
formatted [NameAndNumberForm] string value.
*/
func (r NameAndNumberForm) String() (val string) {
	n := r.primaryIdentifier.String()
	if len(r.identifier) == 0 {
		return n
	}
	return sprintf("%s(%s)", r.identifier, n)
}

/*
Equal returns a Boolean value indicative of whether instance
n of [NameAndNumberForm] matches the receiver.
*/
func (r NameAndNumberForm) Equal(n any) (is bool) {
	switch tv := n.(type) {
	case NameAndNumberForm:
		is = r.identifier == tv.identifier &&
			r.primaryIdentifier.Equal(tv.primaryIdentifier)
	case *NameAndNumberForm:
		is = r.identifier == tv.identifier &&
			r.primaryIdentifier.Equal(tv.primaryIdentifier)
	}

	return
}

func parseRootNameOnly(x string) (r *NameAndNumberForm, err error) {
	var root *big.Int
	switch x {
	case `itu-t`:
		root = big.NewInt(0)
	case `iso`:
		root = big.NewInt(1)
	case `joint-iso-itu-t`:
		root = big.NewInt(2)
	default:
		err = errorf("Unknown root abbreviation, or no closing NumberForm parenthesis to read")
	}

	if err == nil {
		r = &NameAndNumberForm{
			parsed:            true,
			identifier:        x,
			primaryIdentifier: NumberForm(*root),
		}
	}

	return
}

/*
parseNaNFstr returns an instance of *[NameAndNumberForm] alongside an error.
*/
func parseNaNFstr(x string) (r *NameAndNumberForm, err error) {
	// Don't waste time on bogus values.
	if len(x) == 0 {
		err = errorf("No content for parseNaNFstr to read")
		return
	} else if x[len(x)-1] != ')' {
		r, err = parseRootNameOnly(x)
		return
	}

	// index the rune for '(', indicating the
	// identifier (nameForm) has ended, and the
	// numberForm is beginning.
	idx := indexRune(x, '(')
	if idx == -1 {
		err = errorf("No opening parenthesis for parseNaNFstr to read")
		return
	}

	// select the numerical characters,
	// or bail out ...
	n := x[idx+1 : len(x)-1]
	if !isNumber(n) {
		err = errorf("Bad numberForm")
		return
	}
	// Parse/verify what appears to be the
	// identifier string value.
	var identifier string = x[:idx]
	if !isIdentifier(identifier) {
		err = errorf("Invalid identifier [%s]; syntax must conform to: LOWER *[ [-] +[ UPPER / LOWER / DIGIT ] ]", identifier)
		return
	}

	// parse the string numberForm value into
	// an instance of NumberForm, or bail out.
	var prid NumberForm
	if prid, err = NewNumberForm(n); err == nil {
		// Prepare to return valid information.
		r = new(NameAndNumberForm)
		r.parsed = true
		r.primaryIdentifier = prid
		r.identifier = x[:idx]
	}

	return
}

func parseNaNFOrNF(tv string) (r *NameAndNumberForm, err error) {
	r = new(NameAndNumberForm)

	if !isNumber(tv) {
		r, err = parseNaNFstr(tv)
	} else {
		var a NumberForm
		if a, err = NewNumberForm(tv); err == nil {
			r = &NameAndNumberForm{primaryIdentifier: a}
		}
	}

	return
}

func parseNaNFBig(tv *big.Int) (r *NameAndNumberForm, err error) {
	r = new(NameAndNumberForm)

	if tv.Int64() < 0 {
		err = errorf("NameAndNumberForm cannot contain a negative NumberForm")
		return
	}

	if len(tv.Bytes()) == 0 {
		err = errorf("NumberForm (%T) is nil", tv)
		return
	}

	r.primaryIdentifier = NumberForm(*tv)

	return
}

/*
NewNameAndNumberForm returns an instance of *[NameAndNumberForm]
alongside an error. Valid input forms are:

• nameAndNumberForm (e.g.: "enterprise(1)"), or ...

• numberForm (e.g.: 1)

[NumberForm] components CANNOT be negative. Permitted input types are
string, uint, uint64, [NumberForm], *[math/big.Int] and (non-negative) int.
*/
func NewNameAndNumberForm(x any) (r *NameAndNumberForm, err error) {

	switch tv := x.(type) {
	case string:
		r, err = parseNaNFOrNF(tv)
	case *big.Int:
		r, err = parseNaNFBig(tv)
	case NumberForm:
		r = new(NameAndNumberForm)
		r.primaryIdentifier = tv
	case uint64:
		u, _ := NewNumberForm(tv) // skip error checking, we know it won't overflow.
		r = new(NameAndNumberForm)
		r.primaryIdentifier = u
	case uint:
		r = new(NameAndNumberForm)
		r, err = NewNameAndNumberForm(uint64(tv))
	case int:
		r = new(NameAndNumberForm)
		if tv < 0 {
			err = errorf("NumberForm cannot be negative")
			break
		}
		r, err = NewNameAndNumberForm(uint64(tv))
	default:
		err = errorf("Unsupported %T input type '%T'", r, tv)
	}

	// mark this instance as complete,
	// assuming no errors were found.
	if err == nil {
		r.parsed = true
	}

	return
}
