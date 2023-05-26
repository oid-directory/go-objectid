package objectid

import (
	"testing"
)

func TestNewOID(t *testing.T) {
	if _, err := NewOID(`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`); err != nil {
		t.Errorf("%s failed: %s", t.Name(), err.Error())
	}
}
