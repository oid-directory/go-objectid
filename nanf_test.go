package objectid

import (
	"testing"
)

func TestNewNameAndNumberForm(t *testing.T) {
	if _, err := NewNameAndNumberForm("enterprise(1)"); err != nil {
		t.Errorf("%s failed: %s",
			t.Name(), err.Error())
	}
}

func TestBogusNameAndNumberForm(t *testing.T) {
	if _, err := NewNameAndNumberForm("Enterprise(1)"); err == nil {
		t.Errorf("%s failed: parsed bogus string value without error", t.Name())
	}
}
