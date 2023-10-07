package objectid

import (
	"fmt"
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

		_, err := NewOID(typ)
		if err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
			return
		}
	}
}

func ExampleOID_Dot() {
	raw := `{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`
	id, err := NewOID(raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s", id.Dot())
	// Output: 1.3.6.1.4.1.56521.999
}

func ExampleOID_Len() {
	raw := `{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`
	id, err := NewOID(raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%d", id.Len())
	// Output: 8
}

func ExampleOID_ASN() {
	raw := `{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}`
	id, err := NewOID(raw)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", id.ASN())
	// Output: {iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999)}
}

func ExampleOID_IsZero() {
	var z OID
	fmt.Printf("Zero: %t", z.IsZero())
	// Output: Zero: true
}

func ExampleOID_Valid() {
	var o OID
	fmt.Printf("Valid: %t", o.Valid())
	// Output: Valid: false
}

func ExampleOID_Leaf() {
	a := `{joint-iso-itu-t(2) uuid(25) ans(987895962269883002155146617097157934)}`
	id, err := NewOID(a)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Leaf node: %s", id.Leaf())
	// Output: Leaf node: ans(987895962269883002155146617097157934)
}

func ExampleOID_Parent() {
	a := `{joint-iso-itu-t(2) uuid(25) ans(987895962269883002155146617097157934)}`
	id, err := NewOID(a)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Leaf node parent: %s", id.Parent())
	// Output: Leaf node parent: uuid(25)
}

func ExampleOID_Root() {
	a := `{joint-iso-itu-t(2) uuid(25) ans(987895962269883002155146617097157934)}`
	id, err := NewOID(a)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Root node: %s", id.Root())
	// Output: Root node: joint-iso-itu-t(2)
}
