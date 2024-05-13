package objectid

import (
	"fmt"
	"math/big"
	"testing"
)

/*
This example demonstrates a means of creating a new instance of [DotNotation]
using variadic input comprised of mixed (supported) type instances.

This may be useful in cases where an instance of [DotNotation] is being created
using relative components derived elsewhere, or otherwise inferred incrementally
in some manner.
*/
func ExampleNewDotNotation_mixedVariadicInput() {
	dot, err := NewDotNotation(uint(1), uint(3), big.NewInt(6), `1`, `4`, `1`, uint64(56521), 999, 5)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s", dot)
	// Output: 1.3.6.1.4.1.56521.999.5
}

func ExampleDotNotation_Index() {
	dot, err := NewDotNotation(`1.3.6.1.4.1.56521.999.5`)
	if err != nil {
		fmt.Println(err)
		return
	}

	arc, _ := dot.Index(1)
	fmt.Printf("%s", arc)
	// Output: 3
}

/*
This example demonstrates use of the [DotNotation.Encode] method, which
returns the ASN.1 encoding of the [DotNotation] instance alongside error.

The result should yield a complete ASN.1 byte sequence, representing the
encoded form of the receiver instance of [DotNotation].
*/
func ExampleDotNotation_Encode() {
	dot, _ := NewDotNotation(`1.3.6.1.4.1.56521.999.5`)
	b, err := dot.Encode()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v", b)
	// Output: [6 11 43 6 1 4 1 131 185 73 135 103 5]
}

/*
This example demonstrates use of the [DotNotation.Decode] method, which
returns an error following an attempt to decode b ([]byte), which is ASN.1
encoded, into the receiver instance.

The result should yield a freshly populated [DotNotation] receiver instance
-- bearing the OID 1.3.6.1.4.1.56521.999.5 -- alongside a nil error.
*/
func ExampleDotNotation_Decode() {
	var d DotNotation

	// pre-encoded bytes for OID 1.3.6.1.4.1.56521.999.5
	b := []byte{0x6, 0xb, 0x2b, 0x6, 0x1, 0x4, 0x1, 0x83, 0xb9, 0x49, 0x87, 0x67, 0x5}

	err := d.Decode(b)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", d)
	// Output: 1.3.6.1.4.1.56521.999.5

}

func ExampleDotNotation_IsZero() {
	var dot DotNotation
	fmt.Printf("Is Zero: %t", dot.IsZero())
	// Output: Is Zero: true
}

func ExampleDotNotation_Valid() {
	var dot DotNotation
	fmt.Printf("Is Valid: %t", dot.Valid())
	// Output: Is Valid: false
}

func ExampleDotNotation_Len() {
	dot, err := NewDotNotation(`1.3.6.1.4.1.56521.999.5`)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Length: %d", dot.Len())
	// Output: Length: 9
}

func ExampleDotNotation_String() {
	dot, err := NewDotNotation(`1.3.6.1.4.1.56521.999.5`)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", dot)
	// Output: 1.3.6.1.4.1.56521.999.5
}

func TestDotNotation_encCodecov(t *testing.T) {
	bi := big.NewInt(0)
	d := DotNotation{NumberForm(*bi)}
	if _, err := d.Encode(); err == nil {
		t.Errorf("%s failed: expected error, got nothing", t.Name())
	}

	bi2 := big.NewInt(7430)
	d = append(d, NumberForm(*bi2))
	if _, err := d.Encode(); err == nil {
		t.Errorf("%s failed: expected error, got nothing", t.Name())
	}

	bi3 := big.NewInt(2)
	bi4 := big.NewInt(7430)
	d = DotNotation{NumberForm(*bi3), NumberForm(*bi4)}
	if _, err := d.Encode(); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
	}
}

func TestDotNotation_decCodecov(t *testing.T) {
	bad := []byte{0x05, 0x00}
	var dot DotNotation
	if err := dot.Decode(bad); err == nil {
		t.Errorf("%s failed: expected error, got nothing", t.Name())
	}
	bad = []byte{0x05, 0x01, 0x01}
	if err := dot.Decode(bad); err == nil {
		t.Errorf("%s failed: expected error, got nothing", t.Name())
	}
	bad = []byte{0x06, 0x02, 0x01}
	if err := dot.Decode(bad); err == nil {
		t.Errorf("%s failed: expected error, got nothing", t.Name())
	}
	bad = []byte{0x06, 0x02, 0x87, 0x67}
	if err := dot.Decode(bad); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
	}
	bad = []byte{0x06, 0x01, 0x51}
	if err := dot.Decode(bad); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
	}
}

func TestDotNotation_badInit(t *testing.T) {
	var d DotNotation
	want := false
	got := d.Valid()
	if want != got {
		t.Errorf("%s failed: wanted validity of %t, got %t",
			t.Name(), want, got)
		return
	}
}

func TestDotNotation_Ancestry(t *testing.T) {
	dot, err := NewDotNotation(`1.3.6.1.4.1.56521.999.5`)
	if err != nil {
		t.Errorf("%s failed: %s",
			t.Name(), err.Error())
		return
	}
	anc := dot.Ancestry()

	want := 9
	got := len(anc)

	if want != got {
		t.Errorf("%s failed: wanted length of %d, got %d",
			t.Name(), want, got)
		return
	}
}

func TestDotNotation_NewSubordinate(t *testing.T) {
	dot, _ := NewDotNotation(`1.3.6.1.4.1.56521.999.5`)
	leaf := dot.NewSubordinate(`10001`)

	want := `1.3.6.1.4.1.56521.999.5.10001`
	got := leaf.String()

	if want != got {
		t.Errorf("%s failed: wanted %s, got %s",
			t.Name(), want, got)
		return
	}

	if !dot.Valid() {
		t.Errorf("%s failed %T validity checks", t.Name(), dot)
		return
	}
}

func TestDotNotation_IsZero(t *testing.T) {
	var dot DotNotation
	if !dot.IsZero() {
		t.Errorf("%s failed: bogus IsZero return",
			t.Name())
		return
	}
}

func TestDotNotation_AncestorOf(t *testing.T) {
	dot, _ := NewDotNotation(`1.3.6`)
	chstr := `1.3.6.1.4`
	child, err := NewDotNotation(chstr)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for _, d := range []any{
		chstr,
		child,
		*child,
	} {
		if !dot.AncestorOf(d) {
			t.Errorf("%s failed: ancestry check returned bogus result",
				t.Name())
			return
		}
	}
}

func TestDotNotation_codecov(t *testing.T) {
	_, err := NewDotNotation(``)
	if err == nil {
		t.Errorf("%s failed: zero length OID parsed without error", t.Name())
		return
	}
	if _, err = NewDotNotation(uint(1), uint(3), rune(6)); err == nil {
		t.Errorf("%s failed: unsupported type accepted by NewDotNotation without error", t.Name())
		return
	}

	var nf1, nf2 NumberForm
	if nf1, err = NewNumberForm(1); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}
	if nf2, err = NewNumberForm(3); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	if _, err = NewDotNotation(nf1, nf2); err != nil {
		t.Errorf("%s failed: NumberForm error: %v", t.Name(), err)
		return
	}

	var nf3, nf4 NumberForm
	if _, err = NewDotNotation(nf3, nf4); err == nil {
		t.Errorf("%s failed: no error where one was expected", t.Name())
		return
	}

	var X DotNotation
	_, _ = X.IntSlice()
	_, _ = X.Uint64Slice()
}

func TestDotNotation_Index(t *testing.T) {
	dot, err := NewDotNotation(`1.3.6.1.4.1.56521.999.5`)
	if err != nil {
		fmt.Println(err)
		return
	}

	nf, _ := dot.Index(1)
	if nf.String() != `3` {
		t.Errorf("%s failed: unable to call index 1 from %T", t.Name(), nf)
		return
	}

	nf, _ = dot.Index(-1)
	if nf.String() != `5` {
		t.Errorf("%s failed: unable to call index -1 from %T", t.Name(), nf)
		return
	}

	nf, _ = dot.Index(100)
	if nf.String() != `5` {
		t.Errorf("%s failed: unable to call index 100 from %T", t.Name(), nf)
		return
	}
}

func TestDotNotation_Codec(t *testing.T) {
	for key, slice := range map[string][]byte{
		`0.0`:     {0x06, 0x01, 0x0},
		`0`:       []byte(`bogus`),
		`2.0`:     {0x06, 0x01, 0x50},
		`2`:       []byte(`bogus`),
		`1.0`:     {0x06, 0x01, 0x28},
		`1`:       []byte(`bogus`),
		`t.1.0..`: []byte(`bogus`),
		`1.t`:     []byte(`bogus`),
		`1.0.`:    []byte(`bogus`),
		`1.0..`:   []byte(`bogus`),
		`3.1`:     []byte(`bogus`),
		`1.2`:     {0x06, 0x01, 0x2a},
		`1.765`:   []byte(`bogus`),
		`2.25`:    {0x06, 0x01, 0x69},
		`2.-25`:   []byte(`bogus`),
		`2.999`:   {0x06, 0x02, 0x87, 0x67},
		`2.`:      []byte(`bogus`),
		`1.3.6.1.4.1.56521.999`: {
			0x06, 0x0a, 0x2b, 0x06, 0x01, 0x04,
			0x01, 0x83, 0xb9, 0x49, 0x87, 0x67},
		``: []byte(`bogus`),
		`2.25.987895962269883002155146617097157934`: {
			0x06, 0x13, 0x69, 0x81, 0xbe, 0xa1, 0xc2,
			0x81, 0xc8, 0xc8, 0xc2, 0x8b, 0x8e, 0xd0,
			0x80, 0x80, 0xaa, 0xae, 0xd7, 0xfa, 0x2e},
	} {
		dot, err := NewDotNotation(key)
		if err != nil {
			if string(slice) != `bogus` {
				t.Errorf("%s failed: valid DotNotation not parsed: %v", t.Name(), err)
			}
			continue
		} else if string(slice) == `bogus` {
			t.Errorf("%s failed: bogus DotNotation (%s) parsed without error", t.Name(), key)
			continue
		}

		dot.encode2Decode(key, slice, t)
	}
}

func (r *DotNotation) encode2Decode(key string, slice []byte, t *testing.T) {
	b, err := r.Encode()
	if err != nil && string(slice) != `bogus` {
		t.Errorf("%s failed: valid DotNotation not encoded: %v", t.Name(), err)
	} else if err == nil && string(slice) == `bogus` {
		t.Errorf("%s failed: bogus DotNotation encoded without error", t.Name())
	} else if string(slice) != `bogus` {
		var r2 DotNotation
		if err = r2.Decode(b); err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
		} else if d2 := r2.String(); d2 != key {
			t.Errorf("%s failed: want '%s', got '%s'", t.Name(), key, d2)
		}
	}
}
