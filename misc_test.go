package objectid

import (
	"testing"
)

var strInSliceMap map[int]map[int][]bool = map[int]map[int][]bool{
	// case match
	0: {
		0: {true, true, true, true, true},
		1: {true, true, true, true, true},
	},

	// case fold
	1: {
		0: {true, true, true, true, true},
		1: {true, true, true, true, true},
	},
}

func TestStrInSlice(t *testing.T) {
	for idx, fn := range []func(string, []string) bool{
		strInSlice,
		strInSliceFold,
	} {
		for i, values := range [][]string{
			{`cAndidate1`, `blarGetty`, `CANndidate7`, `squatcobbler`, `<censored>`},
			{`Ó-aîï4Åø´øH«w%);<wÃ¯`, `piles`, `4378295fmitty`, string(rune(0)), `broccolI`},
		} {
			for j, val := range values {
				result_expected := strInSliceMap[idx][i][j]

				// warp the candidate value such that
				// it no longer matches the slice from
				// whence it originates. j² is used as
				// its quicker and less stupid than
				// adding a rand generator.
				if isPowerOfTwo(j) {
					var R []rune = []rune(val)
					for g, h := 0, len(R)-1; g < h; g, h = g+1, h-1 {
						R[g], R[h] = R[h], R[g]
					}
					val = string(R)
					result_expected = !result_expected // invert
				}

				result_received := fn(val, values)
				if result_expected != result_received {
					t.Errorf("%s[%d->%d] failed; []byte(%v) in %v: %t (wanted %t)",
						t.Name(), i, j, []byte(val), values, result_received, result_expected)
					return
				}
			}
		}
	}
}

func TestMisc_codecov(t *testing.T) {
	_ = errorf("this is a string %s", `error`)
	_ = errorf(errorf("this is an error"))
}

func TestIntSlices(t *testing.T) {
	for idx, pair := range [][][]int{
		{
			{1, 2, 3, 4},
			{1, 2, 3, 4},
		},
		{
			{1, 2, 3},
			{4, 5, 6, 7},
		},
		{
			{8, 8},
			{8, 8},
		},
		{
			{1, 2, 3, 5},
			{1, 5, 3, 7},
		},
	} {
		eql := intSliceEqual(pair[0], pair[1])
		if !eql && idx%2 == 0 {
			t.Errorf("%s failed: matching int slices not deemed equal", t.Name())
			return
		} else if eql && idx%2 != 0 {
			t.Errorf("%s failed: non-matching int slices deemed equal", t.Name())
			return
		}
	}
}

func TestStrSlices(t *testing.T) {
	for idx, pair := range [][][]string{
		{
			{`1`, `2`, `3`, `4`},
			{`1`, `2`, `3`, `4`},
		},
		{
			{`1`, `2`, `3`},
			{`4`, `5`, `6`, `7`},
		},
		{
			{`8`, `8`},
			{`8`, `8`},
		},
		{
			{`1`, `2`, `3`, `5`},
			{`1`, `5`, `3`, `7`},
		},
	} {
		eql := strSliceEqual(pair[0], pair[1])
		if !eql && idx%2 == 0 {
			t.Errorf("%s failed: matching str slices not deemed equal", t.Name())
			return
		} else if eql && idx%2 != 0 {
			t.Errorf("%s failed: non-matching str slices deemed equal", t.Name())
			return
		}
	}
}

/*
func TestStrSlices(t *testing.T) {

        strSlices := [][]string{
                {`1`, `2`, `3`, `4`},
                {`1`, `3`, `2`, `4`},
        }

	if strSliceEqual([]string{`1`, `2`, `3`}, []string{`4`, `5`, `6`, `7`}) {
		t.Errorf("%s failed: non-matching int slices deemed equal", t.Name())
		return
	}

}
*/

func TestIsNumber(t *testing.T) {
	for idx, candidate := range []string{
		`1`,
		`t18`,
		`17`,
		`S`,
		`8`,
		`~ej`,
		`11`,
		`S&*(D`,
		`100`,
		``,
		`9919387`,
	} {
		var err error
		is := isNumber(candidate)
		if is && idx%2 != 0 {
			err = errorf("%s failed: good value [%s] not cleared as number", t.Name(), candidate)
		} else if !is && idx%2 == 0 {
			err = errorf("%s failed: bogus value [%s] cleared as number", t.Name(), candidate)
		}

		if err != nil {
			t.Errorf("%v", err)
			return
		}
	}
}

func TestIsIdentifier(t *testing.T) {
	for idx, candidate := range []string{
		`enterprise`,
		`Enterprise`,
		`iso`,
		`itu-`,
		`telcoCompany`,
		`-enterprise`,
		`identified-organization`,
		`100`,
		`joint-iso-itu-t`,
		``,
		`itu-t`,
		`itu?t`,
	} {
		var err error
		is := isIdentifier(candidate)
		if is && idx%2 != 0 {
			err = errorf("%s failed: good value [%s] not cleared as an identifier", t.Name(), candidate)
		} else if !is && idx%2 == 0 {
			err = errorf("%s failed: bogus value [%s] cleared as an identifier", t.Name(), candidate)
		}

		if err != nil {
			t.Errorf("%v", err)
			return
		}
	}
}
