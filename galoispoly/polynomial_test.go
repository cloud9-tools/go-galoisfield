package galoispoly

import (
	"testing"

	"github.com/cloud9-tools/go-galoisfield"
)

func TestNewPolynomial(t *testing.T) {
	type testrow struct {
		input Polynomial
		str string
		gostr string
		field *galoisfield.GF
		deg uint
		coeff []byte
	}
	for idx, row := range []testrow{
		testrow{NewPolynomial(nil),
			"0",
			"NewPolynomial(Poly84320_g2)",
			galoisfield.Poly84320_g2, 0, nil},
		testrow{NewPolynomial(nil, 1),
			"1",
			"NewPolynomial(Poly84320_g2, 1)",
			galoisfield.Poly84320_g2, 0, []byte{1}},
		testrow{NewPolynomial(nil, 2),
			"2",
			"NewPolynomial(Poly84320_g2, 2)",
			galoisfield.Poly84320_g2, 0, []byte{2}},
		testrow{NewPolynomial(nil, 17),
			"17",
			"NewPolynomial(Poly84320_g2, 17)",
			galoisfield.Poly84320_g2, 0, []byte{17}},
		testrow{NewPolynomial(nil, 0, 2),
			"2x",
			"NewPolynomial(Poly84320_g2, 0, 2)",
			galoisfield.Poly84320_g2, 1, []byte{0, 2}},
		testrow{NewPolynomial(nil, 1, 2),
			"2x + 1",
			"NewPolynomial(Poly84320_g2, 1, 2)",
			galoisfield.Poly84320_g2, 1, []byte{1, 2}},
		testrow{NewPolynomial(nil, 1, 0, 1),
			"x^2 + 1",
			"NewPolynomial(Poly84320_g2, 1, 0, 1)",
			galoisfield.Poly84320_g2, 2, []byte{1, 0, 1}},
		testrow{NewPolynomial(nil, 0, 1, 1),
			"x^2 + x",
			"NewPolynomial(Poly84320_g2, 0, 1, 1)",
			galoisfield.Poly84320_g2, 2, []byte{0, 1, 1}},
		testrow{NewPolynomial(nil, 0, 1, 1, 0),
			"x^2 + x",
			"NewPolynomial(Poly84320_g2, 0, 1, 1)",
			galoisfield.Poly84320_g2, 2, []byte{0, 1, 1}},
		testrow{NewPolynomial(nil, 3, 1, 4),
			"4x^2 + x + 3",
			"NewPolynomial(Poly84320_g2, 3, 1, 4)",
			galoisfield.Poly84320_g2, 2, []byte{3, 1, 4}},
	} {
		str := row.input.String()
		if str != row.str {
			t.Errorf("[%2d] expected %q, got %q", idx, row.str, str)
		}
		gostr := row.input.GoString()
		if gostr != row.gostr {
			t.Errorf("[%2d] expected %q, got %q", idx, row.gostr, gostr)
		}
		field := row.input.Field()
		if field != row.field {
			t.Errorf("[%2d] expected %#v, got %#v", idx, row.field, field)
		}
		deg := row.input.Degree()
		if deg != row.deg {
			t.Errorf("[%2d] expected %d, got %d", idx, row.deg, deg)
		}
		coeff := row.input.Coefficients()
		if !equalBytes(coeff, row.coeff) {
			t.Errorf("[%2d] expected %v, got %v", idx, row.coeff, coeff)
		}
		for i, k := range row.coeff {
			actual := row.input.Coefficient(uint(i))
			if actual != k {
				t.Errorf("[%2d] expected %d, got %d", idx, k, actual)
			}
		}
		for i := 0; i < len(row.coeff); i++ {
			actual := row.input.Coefficient(uint(i+len(row.coeff)))
			if actual != 0 {
				t.Errorf("[%2d] expected 0, got %d", idx, actual)
			}
		}
	}
}

func TestPolynomial_Scale(t *testing.T) {
	type testrow struct {
		scalar   byte
		input    Polynomial
		expected Polynomial
	}
	for _, row := range []testrow{
		testrow{5,
			NewPolynomial(nil, 3, 0, 1),
			NewPolynomial(nil, 15, 0, 5)},
		testrow{1,
			NewPolynomial(nil, 3, 0, 1),
			NewPolynomial(nil, 3, 0, 1)},
		testrow{0,
			NewPolynomial(nil, 3, 0, 1),
			NewPolynomial(nil)},
	} {
		actual := row.input.Scale(row.scalar)
		if !actual.Equal(row.expected) {
			t.Errorf("expected %d*(%v)=(%v), got (%v)",
				row.scalar, row.input, row.expected, actual)
		}
	}
}

func TestPolynomial_Compare(t *testing.T) {
	type testrow struct {
		a        Polynomial
		b        Polynomial
		expected int
	}
	for _, row := range []testrow{
		testrow{NewPolynomial(nil), NewPolynomial(nil), 0},
		testrow{NewPolynomial(nil, 5), NewPolynomial(nil, 5), 0},
		testrow{NewPolynomial(nil, 3, 5), NewPolynomial(nil, 3, 5), 0},
		testrow{NewPolynomial(nil), NewPolynomial(nil, 1), -1},
		testrow{NewPolynomial(nil, 0), NewPolynomial(nil, 1), -1},
		testrow{NewPolynomial(nil, 2, 1), NewPolynomial(nil, 1, 2), -1},
		testrow{NewPolynomial(galoisfield.Poly310_g2), NewPolynomial(galoisfield.Poly210_g2), 1},
	} {
		a, b, expected := row.a, row.b, row.expected
		actual := a.Compare(b)
		if actual != expected {
			t.Errorf("expected %#v cmp %#v == %d, got %d", a, b, expected, actual)
		}
		checkCompareAxioms(
			t, a, b, actual,
			a.Less(b),
			b.Less(a),
			a.Equal(b),
			b.Equal(a))
	}
}

func TestPolynomial_Add(t *testing.T) {
	type testrow struct {
		a, b Polynomial
		expected Polynomial
	}
	for _, row := range []testrow{
		testrow{NewPolynomial(nil, 1, 0, 0, 1),
			NewPolynomial(nil),
			NewPolynomial(nil, 1, 0, 0, 1)},
		testrow{NewPolynomial(nil, 1, 0, 0, 1),
			NewPolynomial(nil, 0, 1),
			NewPolynomial(nil, 1, 1, 0, 1)},
		testrow{NewPolynomial(nil, 1, 0, 0, 1),
			NewPolynomial(nil, 0, 0, 1, 1),
			NewPolynomial(nil, 1, 0, 1)},
	} {
		actual := row.a.Add(row.b)
		if !actual.Equal(row.expected) {
			t.Errorf("expected (%v)+(%v)=(%v), got %v", row.a, row.b, row.expected, actual)
		}
	}
}

func TestPolynomial_Add_axioms(t *testing.T) {
	zero := NewPolynomial(nil)
	add := func(x, y interface{}) interface{} {
		return x.(Polynomial).Add(y.(Polynomial))
	}
	eq := func(x, y interface{}) bool {
		return x.(Polynomial).Equal(y.(Polynomial))
	}
	type testrow struct {
		a, b, c Polynomial
	}
	for _, row := range []testrow{
		testrow{NewPolynomial(nil, 5, 0, 0, 7),
			NewPolynomial(nil, 0, 3),
			NewPolynomial(nil, 0, 0, 1)},
		testrow{NewPolynomial(nil, 5, 0, 0, 7),
			NewPolynomial(nil, 0, 3),
			NewPolynomial(nil, 1, 0, 1, 2)},
		testrow{NewPolynomial(nil, 5, 0, 0, 7),
			NewPolynomial(nil, 0, 3),
			NewPolynomial(nil, 5, 0, 1, 7)},
	} {
		checkAddAxioms(t, row.a, row.b, row.c, zero, add, eq)
	}
}

func TestPolynomial_Add_incompatible(t *testing.T) {
	e := panicValue(func() {
		_ = NewPolynomial(galoisfield.Poly210_g2).
			Add(NewPolynomial(galoisfield.Poly310_g2))
	})
	if e != ErrIncompatibleFields {
		t.Errorf("expected ErrIncompatibleFields, got %q", e.Error())
	}
}

func checkCompareAxioms(t *testing.T, a, b interface{}, cmp int, lt, gt, eq, qe bool) {
	if eq != qe {
		t.Errorf("equality not commutative for %#v and %#v", a, b)
	}
	if eq && lt {
		t.Errorf("equality and lessthan not disjoint for %#v and %#v", a, b)
	}
	if eq && gt {
		t.Errorf("equality and greaterthan not disjoint for %#v and %#v", a, b)
	}
	if cmp < 0 && !lt {
		t.Errorf("expected %#v < %#v, got >=", a, b)
	}
	if cmp == 0 && !eq {
		t.Errorf("expected %#v == %#v, got !=", a, b)
	}
	if cmp > 0 && !gt {
		t.Errorf("expected %#v > %#v, got <=", a, b)
	}
}

func checkAddAxioms(t *testing.T, a, b, c, z interface{}, add func(_, _ interface{}) interface{}, eq func(_, _ interface{}) bool) {
	// a+0 = 0+a = a
	az := add(a, z)
	za := add(z, a)
	if !eq(az, za) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", a, z, az, za)
	} else if !eq(a, az) {
		t.Errorf("additive 'identity' isn't for %#v and %#v; got x+0=0+x=%#v", a, z, az)
	}

	// b+0 = 0+b = b
	bz := add(b, z)
	zb := add(z, b)
	if !eq(bz, zb) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", b, z, bz, zb)
	} else if !eq(b, bz) {
		t.Errorf("additive 'identity' isn't for %#v and %#v; got x+0=0+x=%#v", b, z, bz)
	}

	// c+0 = 0+c = b
	cz := add(c, z)
	zc := add(z, c)
	if !eq(cz, zc) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", c, z, cz, zc)
	} else if !eq(c, cz) {
		t.Errorf("additive 'identity' isn't for %#v and %#v; got x+0=0+x=%#v", c, z, cz)
	}

	// 0+0 = 0
	zz := add(z, z)
	if !eq(z, zz) {
		t.Errorf("additive 'identity' isn't for %#v and itself; got 0+0=%#v", z, zz)
	}

	// a+b = b+a
	ab := add(a, b)
	ba := add(b, a)
	if !eq(ab, ba) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", a, b, ab, ba)
	}
	// a+c = c+a
	ac := add(a, c)
	ca := add(c, a)
	if !eq(ac, ca) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", a, c, ac, ca)
	}
	// b+c = c+b
	bc := add(b, c)
	cb := add(c, b)
	if !eq(bc, cb) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", b, c, bc, cb)
	}

	// (a+b)+c = c+(a+b)
	ab_c := add(ab, c)
	c_ab := add(c, ab)
	if !eq(ab_c, c_ab) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", ab, c, ab_c, c_ab)
	}
	// a+(b+c) = (b+c)+a
	a_bc := add(a, bc)
	bc_a := add(bc, a)
	if !eq(a_bc, bc_a) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", a, bc, a_bc, bc_a)
	}
	// (a+b)+c = a+(b+c)
	if !eq(ab_c, a_bc) {
		t.Errorf("addition not associative; got (a+b)+c=%#v, a+(b+c)=%#v", ab_c, a_bc)
	}
}

func panicValue(f func()) (value error) {
	defer func() {
		if e, ok := recover().(error); ok {
			value = e
		}
	}()
	f()
	return
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
