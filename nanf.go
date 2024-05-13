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
func (nanf NameAndNumberForm) IsZero() bool {
	return !nanf.parsed
}

/*
Identifier returns the string-based nameForm
value assigned to the receiver instance.
*/
func (nanf NameAndNumberForm) Identifier() string {
	return nanf.identifier
}

/*
NumberForm returns the underlying [NumberForm]
value assigned to the receiver instance.
*/
func (nanf NameAndNumberForm) NumberForm() NumberForm {
	return nanf.primaryIdentifier
}

/*
String is a stringer method that returns the properly
formatted [NameAndNumberForm] string value.
*/
func (nanf NameAndNumberForm) String() (val string) {
	n := nanf.primaryIdentifier.String()
	if len(nanf.identifier) == 0 {
		return n
	}
	return sprintf("%s(%s)", nanf.identifier, n)
}

/*
Equal returns a Boolean value indicative of whether instance
n of [NameAndNumberForm] matches the receiver.
*/
func (nanf NameAndNumberForm) Equal(n any) (is bool) {
	switch tv := n.(type) {
	case NameAndNumberForm:
		is = nanf.identifier == tv.identifier &&
			nanf.primaryIdentifier.Equal(tv.primaryIdentifier)
	case *NameAndNumberForm:
		is = nanf.identifier == tv.identifier &&
			nanf.primaryIdentifier.Equal(tv.primaryIdentifier)
	}

	return
}

func parseRootNameOnly(x string) (nanf *NameAndNumberForm, err error) {
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
		nanf = &NameAndNumberForm{
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
func parseNaNFstr(x string) (nanf *NameAndNumberForm, err error) {
	// Don't waste time on bogus values.
	if len(x) == 0 {
		err = errorf("No content for parseNaNFstr to read")
		return
	} else if x[len(x)-1] != ')' {
		nanf, err = parseRootNameOnly(x)
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
		nanf = new(NameAndNumberForm)
		nanf.parsed = true
		nanf.primaryIdentifier = prid
		nanf.identifier = x[:idx]
	}

	return
}

func parseNaNFOrNF(tv string) (nanf *NameAndNumberForm, err error) {
	nanf = new(NameAndNumberForm)

	if !isNumber(tv) {
		nanf, err = parseNaNFstr(tv)
	} else {
		var a NumberForm
		if a, err = NewNumberForm(tv); err == nil {
			nanf = &NameAndNumberForm{primaryIdentifier: a}
		}
	}

	return
}

func parseNaNFBig(tv *big.Int) (nanf *NameAndNumberForm, err error) {
	nanf = new(NameAndNumberForm)

	if tv.Int64() < 0 {
		err = errorf("NameAndNumberForm cannot contain a negative NumberForm")
		return
	}

	if len(tv.Bytes()) == 0 {
		err = errorf("NumberForm (%T) is nil", tv)
		return
	}

	nanf.primaryIdentifier = NumberForm(*tv)

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
func NewNameAndNumberForm(x any) (nanf *NameAndNumberForm, err error) {

	switch tv := x.(type) {
	case string:
		nanf, err = parseNaNFOrNF(tv)
	case *big.Int:
		nanf, err = parseNaNFBig(tv)
	case NumberForm:
		nanf = new(NameAndNumberForm)
		nanf.primaryIdentifier = tv
	case uint64:
		u, _ := NewNumberForm(tv) // skip error checking, we know it won't overflow.
		nanf = new(NameAndNumberForm)
		nanf.primaryIdentifier = u
	case uint:
		nanf = new(NameAndNumberForm)
		nanf, err = NewNameAndNumberForm(uint64(tv))
	case int:
		nanf = new(NameAndNumberForm)
		if tv < 0 {
			err = errorf("NumberForm cannot be negative")
			break
		}
		nanf, err = NewNameAndNumberForm(uint64(tv))
	default:
		err = errorf("Unsupported %T input type '%T'", nanf, tv)
	}

	// mark this instance as complete,
	// assuming no errors were found.
	if err == nil {
		nanf.parsed = true
	}

	return
}
