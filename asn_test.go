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
