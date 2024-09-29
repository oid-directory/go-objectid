package objectid

import (
	"fmt"
	"math/big"
	"testing"
)

const (
	testASN1NotationISO  = `{iso(1)}`
	testASN1JesseExample = `{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`
	testASN1Bogus        = `iso(1) identified-organization(3}`
)

func ExampleASN1Notation_Dot() {
	aNot, err := NewASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`)
	if err != nil {
		fmt.Println(err)
		return
	}
	dot := aNot.Dot()
	fmt.Println(dot)
	// Output: 1.3.6.1.4.1.56521.999
}

/*
This example demonstrates a bogus [DotNotation] output due to the presence
of less than two (2) [NameAndNumberForm] instances within the receiver.

[DotNotation] ALWAYS requires two (2) or more elements to be considered valid.
*/
func ExampleASN1Notation_Dot_bogus() {
	aNot, err := NewASN1Notation(`{iso(1)}`)
	if err != nil {
		fmt.Println(err)
		return
	}
	dot := aNot.Dot()
	fmt.Println(dot)
	// Output:
}

func ExampleASN1Notation_Index() {
	aNot, err := NewASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`)
	if err != nil {
		fmt.Println(err)
		return
	}

	nanf, _ := aNot.Index(1)
	fmt.Printf("%s", nanf)
	// Output: identified-organization(3)
}

func ExampleASN1Notation_Leaf() {
	aNot, err := NewASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", aNot.Leaf())
	// Output: example(999)
}

func ExampleASN1Notation_Parent() {
	aNot, err := NewASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", aNot.Parent())
	// Output: 56521
}

func ExampleASN1Notation_String() {
	aNot, err := NewASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", aNot)
	// Output: {iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}
}

func ExampleASN1Notation_Root() {
	aNot, err := NewASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", aNot.Root())
	// Output: iso(1)
}

func ExampleASN1Notation_Len() {
	aNot, err := NewASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Length: %d", aNot.Len())
	// Output: Length: 8
}

func ExampleASN1Notation_IsZero() {
	var aNot ASN1Notation
	fmt.Printf("Is Zero: %t", aNot.IsZero())
	// Output: Is Zero: true
}

func ExampleASN1Notation_Valid() {
	var aNot ASN1Notation
	fmt.Printf("Is Valid: %t", aNot.Valid())
	// Output: Is Valid: false
}

func TestASN1Notation001(t *testing.T) {
	asn, err := NewASN1Notation(testASN1NotationISO)
	if err != nil {
		t.Errorf("%s error: %s", t.Name(), err.Error())
		return
	}

	want := testASN1NotationISO
	got := asn.String()
	if want != got {
		t.Errorf("%s failed: want '%s', got '%s'",
			t.Name(), want, got)
		return
	}

	asn, err = NewASN1Notation([]string{
		`iso(1)`,
		`identified-organization(3)`,
		`dod(6)`,
		`internet(1)`,
		`private(4)`,
		`enterprise(1)`,
		`56521`,
		`example(999)`})

	if l := asn.Len(); l != 8 {
		t.Errorf("%s failed: want '%d', got '%d'",
			t.Name(), 8, l)
		return
	}

	if _, err = NewASN1Notation(float64(123)); err == nil {
		t.Errorf("%s failed; no error where one was expected", t.Name())
		return
	}

	if _, err = NewASN1Notation([]NameAndNumberForm{
		{identifier: `iso`, primaryIdentifier: NumberForm(*big.NewInt(1)), parsed: true},
		{identifier: `identified-organization`, primaryIdentifier: NumberForm(*big.NewInt(3)), parsed: true},
	}); err != nil {
		t.Errorf("%s error: %v", t.Name(), err)
		return
	}

	if _, err = NewASN1Notation([]NameAndNumberForm{}); err == nil {
		t.Errorf("%s error: %v", t.Name(), err)
		return
	}
}

func TestASN1Notation_Index(t *testing.T) {
	aNot, err := NewASN1Notation(testASN1JesseExample)
	if err != nil {
		fmt.Println(err)
		return
	}

	for key, value := range map[int]string{
		1:   `identified-organization`,
		-1:  `example`,
		-18: `iso`,
		100: `example`,
	} {
		if nanf, _ := aNot.Index(key); nanf.Identifier() != value {
			t.Errorf("%s failed: unable to call index %d from %T;\nwant '%s'\ngot '%s'",
				t.Name(), key, nanf, value, nanf.Identifier())
			return
		}
	}
}

func TestASN1Notation_bogus(t *testing.T) {
	if _, err := NewASN1Notation(testASN1Bogus); err == nil {
		t.Errorf("%s successfully parsed bogus value; expected an error", t.Name())
		return
	}

	if _, err := NewASN1Notation(`iso(3) identified-organization(3)`); err == nil {
		t.Errorf("%s successfully parsed bogus value; expected an error", t.Name())
		return
	}

	if _, err := NewASN1Notation(`itu-t recommendation(-3)`); err == nil {
		t.Errorf("%s successfully parsed bogus value; expected an error", t.Name())
		return
	}

	if _, err := NewASN1Notation(`joint-iso-itu-t thing`); err == nil {
		t.Errorf("%s successfully parsed bogus value; expected an error", t.Name())
		return
	}
}

func TestASN1Notation_Ancestry(t *testing.T) {
	asn, err := NewASN1Notation(testASN1JesseExample)
	if err != nil {
		t.Errorf("%s failed: %s",
			t.Name(), err.Error())
		return
	}
	anc := asn.Ancestry()

	want := 8
	got := len(anc)

	if want != got {
		t.Errorf("%s failed: wanted length of %d, got %d",
			t.Name(), want, got)
		return
	}
}

func TestASN1Notation_NewSubordinate(t *testing.T) {
	asn, _ := NewASN1Notation(testASN1JesseExample)
	leaf := asn.NewSubordinate(`friedChicken(5)`)

	want := `{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999) friedChicken(5)}`
	got := leaf.String()

	if want != got {
		t.Errorf("%s failed: wanted %s, got %s",
			t.Name(), want, got)
		return
	}
}

func TestASN1Notation_IsZero(t *testing.T) {
	var asn ASN1Notation
	if !asn.IsZero() {
		t.Errorf("%s failed: bogus IsZero return",
			t.Name())
		return
	}
}

func TestASN1Notation_AncestorChildOf(t *testing.T) {
	asn, _ := NewASN1Notation(`{iso(1)}`)
	chstr := `{iso(1) identified-organization(3)}`
	child, err := NewASN1Notation(chstr)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for _, d := range []any{
		chstr,
		child,
		*child,
	} {
		if !asn.AncestorOf(d) {
			t.Errorf("%s failed: ancestor check returned bogus result",
				t.Name())
			return
		}
	}
}

func TestASN1Notation_ChildOf(t *testing.T) {
	asn, _ := NewASN1Notation(`{iso(1)}`)
	chstr := `{iso(1) identified-organization(3)}`
	child, err := NewASN1Notation(chstr)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for _, d := range []any{
		chstr,
		child,
		*child,
	} {
		if !asn.ChildOf(d) {
			t.Errorf("%s failed: child check returned bogus result",
				t.Name())
			return
		}
	}
}

func TestASN1Notation_SiblingOf(t *testing.T) {
	asn, _ := NewASN1Notation(`{joint-iso-itu-t(2) uuid(25)}`)
	sibstr := `{joint-iso-itu-t(2) asn1(1)}`
	sib, err := NewASN1Notation(sibstr)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for _, d := range []any{
		sibstr,
		sib,
		*sib,
	} {
		if !asn.SiblingOf(d) {
			t.Errorf("%s failed: sibling check returned bogus result",
				t.Name())
			return
		}
	}
}
