package objectid

import "fmt"

func ExampleParseNumberForm() {
	// single UUID integer parse example
	arc, err := ParseNumberForm(`987895962269883002155146617097157934`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", arc)
	// Output: 987895962269883002155146617097157934
}

func ExampleNewOID() {
	// UUID-based (uint128) OID example
	a := `{joint-iso-itu-t(2) uuid(25) ans(987895962269883002155146617097157934)}`
	id, err := NewOID(a)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("ASN.1 Notation: %s", id.ASN())
	// Output: ASN.1 Notation: {joint-iso-itu-t(2) uuid(25) ans(987895962269883002155146617097157934)}
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

func ExampleNewDotNotation() {
	a := `2.25.987895962269883002155146617097157934`
	id, err := NewDotNotation(a)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("dotNotation: %s", id)
	// Output: dotNotation: 2.25.987895962269883002155146617097157934
}

func ExampleDotNotation_Leaf() {
	a := `2.25.987895962269883002155146617097157934`
	id, err := NewDotNotation(a)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Leaf node: %s", id.Leaf())
	// Output: Leaf node: 987895962269883002155146617097157934
}

func ExampleDotNotation_Parent() {
	a := `2.25.987895962269883002155146617097157934`
	id, err := NewDotNotation(a)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Leaf node parent: %s", id.Parent())
	// Output: Leaf node parent: 25
}

func ExampleDotNotation_Root() {
	a := `2.25.987895962269883002155146617097157934`
	id, err := NewDotNotation(a)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Root node: %s", id.Root())
	// Output: Root node: 2
}

func ExampleNewASN1Notation() {
	a := `{iso(1) identified-organization(3) dod(6)}`
	id, err := NewASN1Notation(a)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Leaf node: %s", id.Leaf())
	// Output: Leaf node: dod(6)
}

func ExampleDotNotation_IntSlice() {
	a := `1.3.6.1.4.1.56521.999.5`
	dot, _ := NewDotNotation(a)

	// If needed, slice instance can be
	// cast as an asn1.ObjectIdentifier.
	slice, err := dot.IntSlice()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v", slice)
	// Output: [1 3 6 1 4 1 56521 999 5]
}

func ExampleDotNotation_IntSlice_overflow() {
	a := `2.25.987895962269883002155146617097157934`
	dot, _ := NewDotNotation(a)
	if _, err := dot.IntSlice(); err != nil {
		fmt.Println(err)
		return
	}
	// Output: strconv.Atoi: parsing "987895962269883002155146617097157934": value out of range
}
