package objectid

import (
	"testing"
)

func TestNewOID(t *testing.T) {
	for _, typ := range []any{
		[]string{
			`iso(1)`,
			`identified-organization(3)`,
			`dod(6)`,
			`internet(1)`,
			`private(4)`,
			`enterprise(1)`,
			`56521`,
			`example(999)`,
		},
		`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`,
	} {

		if _, err := NewOID(typ); err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
			return
		}
	}
}
