package objectid

import (
	"testing"
)

func TestNewNumberForm(t *testing.T) {
	nf := `3849141823758536772162786183725055278`
	if _, err := NewNumberForm(nf); err != nil {
		t.Errorf("%s failed: %s", t.Name(), err.Error())
	}
}

func TestBogusNewNumberForm(t *testing.T) {
	bogus := `-48675`
	if _, err := NewNumberForm(bogus); err == nil {
		t.Errorf("%s failed: bogus NumberForm '%v' accepted without error",
			t.Name(), bogus)
	}
}

func TestNumberForm_Gt(t *testing.T) {
	nf, _ := NewNumberForm(4658)
	if !nf.Gt(3700) {
		t.Errorf("%s failed: Gt evaluation returned a bogus value", t.Name())
	}

	if nf, _ = NewNumberForm(`437829765`); nf.Gt(500000000) {
		t.Errorf("%s failed: Gt evaluation returned a bogus value", t.Name())
	}

	if nf, _ = NewNumberForm(uint64(829765)); !nf.Lt(500000000) {
		t.Errorf("%s failed: Gt evaluation returned a bogus value", t.Name())
	}
}

func TestNumberForm_Lt(t *testing.T) {
	nf, _ := NewNumberForm(4658)
	if nf.Lt(3700) {
		t.Errorf("%s failed: Gt evaluation returned a bogus value", t.Name())
	}

	if nf, _ = NewNumberForm(`437829765`); !nf.Lt(500000000) {
		t.Errorf("%s failed: Gt evaluation returned a bogus value", t.Name())
	}

	if nf, _ = NewNumberForm(uint64(329856)); !nf.Lt(500000000) {
		t.Errorf("%s failed: Gt evaluation returned a bogus value", t.Name())
	}
}

func TestNumberForm_Equal(t *testing.T) {
	nf, _ := NewNumberForm(4658)
	if nf.Equal(3700) {
		t.Errorf("%s failed: Gt evaluation returned a bogus value", t.Name())
	}

	if nf, _ = NewNumberForm(`437829765`); !nf.Equal(437829765) {
		t.Errorf("%s failed: Gt evaluation returned a bogus value", t.Name())
	}

	if nf, _ = NewNumberForm(uint64(329856)); !nf.Equal(uint64(329856)) {
		t.Errorf("%s failed: Gt evaluation returned a bogus value", t.Name())
	}
}
