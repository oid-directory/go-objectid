package objectid

import (
	"testing"
)

const (
	testASN1NotationISO  = `{iso(1)}`
	testASN1JesseExample = `{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`
	testASN1Bogus        = `iso(1) identified-organization(3}`
)

func TestASN1Notation001(t *testing.T) {
	asn, err := NewASN1Notation(testASN1NotationISO)
	if err != nil {
		t.Errorf("%s error: %s\n", t.Name(), err.Error())
		return
	}

	want := testASN1NotationISO
	got := asn.String()
	if want != got {
		t.Errorf("%s failed: want '%s', got '%s'",
			t.Name(), want, got)
		return
	}
}

func TestASN1Notation_Index(t *testing.T) {
	asn, err := NewASN1Notation(testASN1JesseExample)
	if err != nil {
		t.Errorf("%s error: %s\n", t.Name(), err.Error())
		return
	}

	want := `example(999)`
	got := asn.Leaf()
	if want != got.String() {
		t.Errorf("%s failed: want '%s', got '%s'",
			t.Name(), want, got)
		return
	}

	want = `56521`
	got = asn.Parent()
	if want != got.String() {
		t.Errorf("%s failed: want '%s', got '%s'",
			t.Name(), want, got)
		return
	}

	want = `iso(1)`
	got = asn.Root()
	if want != got.String() {
		t.Errorf("%s failed: want '%s', got '%s'",
			t.Name(), want, got)
		return
	}
}

func TestASN1Notation_bogus(t *testing.T) {
	if _, err := NewASN1Notation(testASN1Bogus); err == nil {
		t.Errorf("%s successfully parsed bogus value; expected an error\n", t.Name())
		return
	}
}

func TestASN1Notation_Ancestry(t *testing.T) {
	asn, err := NewASN1Notation(testASN1JesseExample)
	if err != nil {
		t.Errorf("%s failed: %s",
			t.Name(), err.Error())
	}
	anc := asn.Ancestry()

	want := 8
	got := len(anc)

	if want != got {
		t.Errorf("%s failed: wanted length of %d, got %d",
			t.Name(), want, got)
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
	}
}

func TestASN1Notation_IsZero(t *testing.T) {
	var asn ASN1Notation
	if !asn.IsZero() {
		t.Errorf("%s failed: bogus IsZero return",
			t.Name())
	}
}

func TestASN1Notation_AncestorOf(t *testing.T) {
	asn, _ := NewASN1Notation(`{joint-iso-itu-t(2) asn1(1)}`)
	child, _ := NewASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1)}`)
	if asn.AncestorOf(child) {
		t.Errorf("%s failed: ancestry check returned bogus result",
			t.Name())
	}
}
