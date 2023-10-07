package objectid

import (
	"fmt"
	"testing"
)

const (
	testASN1NotationISO  = `{iso(1)}`
	testASN1JesseExample = `{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`
	testASN1Bogus        = `iso(1) identified-organization(3}`
)

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
}

func TestASN1Notation_Index(t *testing.T) {
	aNot, err := NewASN1Notation(testASN1JesseExample)
	if err != nil {
		fmt.Println(err)
		return
	}

	nanf, _ := aNot.Index(1)
	if nanf.Identifier() != `identified-organization` {
		t.Errorf("%s failed: unable to call index 1 from %T; got '%s'", t.Name(), nanf, nanf.Identifier())
		return
	}

	nanf, _ = aNot.Index(-1)
	if nanf.NumberForm().String() != `999` {
		t.Errorf("%s failed: unable to call index -1 from %T", t.Name(), nanf)
		return
	}

	nanf, _ = aNot.Index(100)
	if nanf.NumberForm().String() != `999` {
		t.Errorf("%s failed: unable to call index 100 from %T", t.Name(), nanf)
		return
	}
}

func TestASN1Notation_bogus(t *testing.T) {
	if _, err := NewASN1Notation(testASN1Bogus); err == nil {
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

func TestASN1Notation_AncestorOf(t *testing.T) {
	asn, err := NewASN1Notation(`{joint-iso-itu-t(2) asn1(1)}`)
	if err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	child, err := NewASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1)}`)
	if err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}
	if asn.AncestorOf(child) {
		t.Errorf("%s failed: ancestry check returned bogus result",
			t.Name())
		return
	}

	if asn.AncestorOf(*child) {
		t.Errorf("%s failed: ancestry check returned bogus result",
			t.Name())
		return
	}

	if asn.AncestorOf(child.String()) {
		t.Errorf("%s failed: ancestry check returned bogus result",
			t.Name())
		return
	}
}
