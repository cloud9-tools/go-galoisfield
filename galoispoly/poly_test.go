package galoispoly

import (
	"testing"

	"github.com/cloud9-tools/go-galoisfield"
)

func TestNewMonomial(t *testing.T) {
	type testrow struct {
		x      Monomial
		field  *galoisfield.GF
		coeff  byte
		degree uint
		gostr  string
		str    string
	}
	for _, row := range []testrow{
		testrow{NewMonomial(nil, 0, 0),
			galoisfield.Poly84320_g2, 0, 0,
			"NewMonomial(Poly84320_g2, 0, 0)",
			"0"},
		testrow{NewMonomial(nil, 1, 0),
			galoisfield.Poly84320_g2, 1, 0,
			"NewMonomial(Poly84320_g2, 1, 0)",
			"1"},
		testrow{NewMonomial(nil, 2, 0),
			galoisfield.Poly84320_g2, 2, 0,
			"NewMonomial(Poly84320_g2, 2, 0)",
			"2"},
		testrow{NewMonomial(nil, 17, 0),
			galoisfield.Poly84320_g2, 17, 0,
			"NewMonomial(Poly84320_g2, 17, 0)",
			"17"},
		testrow{NewMonomial(nil, 0, 1),
			galoisfield.Poly84320_g2, 0, 0,
			"NewMonomial(Poly84320_g2, 0, 0)",
			"0"},
		testrow{NewMonomial(nil, 1, 1),
			galoisfield.Poly84320_g2, 1, 1,
			"NewMonomial(Poly84320_g2, 1, 1)",
			"x"},
		testrow{NewMonomial(nil, 2, 1),
			galoisfield.Poly84320_g2, 2, 1,
			"NewMonomial(Poly84320_g2, 2, 1)",
			"2x"},
		testrow{NewMonomial(nil, 17, 1),
			galoisfield.Poly84320_g2, 17, 1,
			"NewMonomial(Poly84320_g2, 17, 1)",
			"17x"},
		testrow{NewMonomial(nil, 0, 2),
			galoisfield.Poly84320_g2, 0, 0,
			"NewMonomial(Poly84320_g2, 0, 0)",
			"0"},
		testrow{NewMonomial(nil, 1, 2),
			galoisfield.Poly84320_g2, 1, 2,
			"NewMonomial(Poly84320_g2, 1, 2)",
			"x^2"},
		testrow{NewMonomial(nil, 2, 2),
			galoisfield.Poly84320_g2, 2, 2,
			"NewMonomial(Poly84320_g2, 2, 2)",
			"2x^2"},
		testrow{NewMonomial(nil, 17, 2),
			galoisfield.Poly84320_g2, 17, 2,
			"NewMonomial(Poly84320_g2, 17, 2)",
			"17x^2"},
		testrow{NewMonomial(nil, 0, 4),
			galoisfield.Poly84320_g2, 0, 0,
			"NewMonomial(Poly84320_g2, 0, 0)",
			"0"},
		testrow{NewMonomial(nil, 17, 4),
			galoisfield.Poly84320_g2, 17, 4,
			"NewMonomial(Poly84320_g2, 17, 4)",
			"17x^4"},
		testrow{NewMonomial(nil, 0, 5),
			galoisfield.Poly84320_g2, 0, 0,
			"NewMonomial(Poly84320_g2, 0, 0)",
			"0"},
		testrow{NewMonomial(nil, 1, 5),
			galoisfield.Poly84320_g2, 1, 5,
			"NewMonomial(Poly84320_g2, 1, 5)",
			"x^5"},
	} {
		if field := row.x.Field(); field != row.field {
			t.Errorf("expected %#v, got %#v", row.field, field)
		}
		if coeff := row.x.Coefficient(); coeff != row.coeff {
			t.Errorf("expected coeff=%d, got %d", row.coeff, coeff)
		}
		if degree := row.x.Degree(); degree != row.degree {
			t.Errorf("expected degree=%d, got %d", row.degree, degree)
		}
		if str := row.x.String(); str != row.str {
			t.Errorf("expected %q, got %q", row.str, str)
		}
		if gostr := row.x.GoString(); gostr != row.gostr {
			t.Errorf("expected %q, got %q", row.gostr, gostr)
		}
	}
}

func TestMonomial_Scale(t *testing.T) {
	type testrow struct {
		coeff    byte
		degree   uint
		scalar   byte
		expected string
	}
	for _, row := range []testrow{
		testrow{0, 0, 0, "0"},
		testrow{1, 0, 0, "0"},
		testrow{1, 1, 0, "0"},
		testrow{1, 5, 0, "0"},

		testrow{0, 0, 1, "0"},
		testrow{1, 0, 1, "1"},
		testrow{1, 1, 1, "x"},
		testrow{1, 5, 1, "x^5"},

		testrow{0, 0, 2, "0"},
		testrow{1, 0, 2, "2"},
		testrow{1, 1, 2, "2x"},
		testrow{1, 5, 2, "2x^5"},
	} {
		m := NewMonomial(nil, row.coeff, row.degree)
		m = m.Scale(row.scalar)
		actual := m.String()
		if actual != row.expected {
			t.Errorf("%#v, got %q", row, actual)
		}
	}
}

func TestMonomial_Add(t *testing.T) {
	type testrow struct {
		a        Monomial
		b        Monomial
		expected Monomial
	}
	for _, row := range []testrow{
		testrow{NewMonomial(nil, 5, 0),
			NewMonomial(nil, 1, 1),
			NewMonomial(nil, 5, 1)},
		testrow{NewMonomial(nil, 5, 1),
			NewMonomial(nil, 3, 1),
			NewMonomial(nil, 15, 2)},
		testrow{NewMonomial(nil, 5, 1),
			NewMonomial(nil, 1, 5),
			NewMonomial(nil, 5, 6)},
		testrow{NewMonomial(nil, 5, 1),
			NewMonomial(nil, 0, 5),
			NewMonomial(nil, 0, 0)},
	} {
		a, b, expected := row.a, row.b, row.expected
		actual := a.Mul(b)
		if !actual.Equals(expected) {
			t.Errorf("expected (%v)*(%v)=(%v), got %v", a, b, expected, actual)
		}
	}
}

func TestMonomial_Compare(t *testing.T) {
	type testrow struct {
		a        Monomial
		b        Monomial
		expected int
	}
	for _, row := range []testrow{
		testrow{NewMonomial(nil, 5, 0), NewMonomial(nil, 3, 1), -1},
		testrow{NewMonomial(nil, 5, 2), NewMonomial(nil, 3, 1), 1},
		testrow{NewMonomial(nil, 5, 1), NewMonomial(nil, 3, 1), 1},
		testrow{NewMonomial(nil, 3, 1), NewMonomial(nil, 3, 1), 0},
		testrow{NewMonomial(nil, 2, 1), NewMonomial(nil, 3, 1), -1},
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
			a.Equals(b),
			b.Equals(a))
	}
}

func TestNewPolynomial(t *testing.T) {
	type testrow struct {
		coeffs   []byte
		expected string
	}
	for _, row := range []testrow{
		testrow{[]byte{}, "0"},
		testrow{[]byte{1}, "1"},
		testrow{[]byte{2}, "2"},
		testrow{[]byte{17}, "17"},
		testrow{[]byte{0, 2}, "2x"},
		testrow{[]byte{1, 2}, "2x + 1"},
		testrow{[]byte{1, 0, 1}, "x^2 + 1"},
		testrow{[]byte{0, 1, 1}, "x^2 + x"},
		testrow{[]byte{3, 1, 4}, "4x^2 + x + 3"},
	} {
		p := NewPolynomial(nil, row.coeffs)
		actual := p.String()
		if actual != row.expected {
			t.Errorf("%#v, got %q", row, actual)
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
			NewPolynomial(nil, []byte{3, 0, 1}),
			NewPolynomial(nil, []byte{15, 0, 5})},
	} {
		actual := row.input.Scale(row.scalar)
		if !actual.Equals(row.expected) {
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
		testrow{NewPolynomial(nil, []byte{}), NewPolynomial(nil, []byte{}), 0},
		testrow{NewPolynomial(nil, []byte{5}), NewPolynomial(nil, []byte{5}), 0},
		testrow{NewPolynomial(nil, []byte{3, 5}), NewPolynomial(nil, []byte{3, 5}), 0},
		testrow{NewPolynomial(nil, []byte{}), NewPolynomial(nil, []byte{1}), -1},
		testrow{NewPolynomial(nil, []byte{0}), NewPolynomial(nil, []byte{1}), -1},
		testrow{NewPolynomial(nil, []byte{2, 1}), NewPolynomial(nil, []byte{1, 2}), -1},
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
			a.Equals(b),
			b.Equals(a))
	}
}

func TestAdd(t *testing.T) {
	type testrow struct {
		a        []byte
		b        []byte
		expected []byte
	}
	for _, row := range []testrow{
		testrow{[]byte{1, 0, 0, 1},
			[]byte{},
			[]byte{1, 0, 0, 1}},
		testrow{[]byte{1, 0, 0, 1},
			[]byte{0, 1},
			[]byte{1, 1, 0, 1}},
		testrow{[]byte{1, 0, 0, 1},
			[]byte{0, 0, 1, 1},
			[]byte{1, 0, 1}},
	} {
		a := NewPolynomial(nil, row.a)
		b := NewPolynomial(nil, row.b)
		expected := NewPolynomial(nil, row.expected)
		actual := Add(a, b)
		if !actual.Equals(expected) {
			t.Errorf("expected (%v)+(%v)=(%v), got %v", a, b, expected, actual)
		}
	}
}

func TestAdd_axioms(t *testing.T) {
	zero := NewPolynomial(nil, nil)
	add := func(x, y interface{}) interface{} {
		return Add(x.(Polynomial), y.(Polynomial))
	}
	eq := func(x, y interface{}) bool {
		return x.(Polynomial).Equals(y.(Polynomial))
	}
	type testrow struct {
		a, b, c []byte
	}
	for _, row := range []testrow{
		testrow{[]byte{5, 0, 0, 7},
			[]byte{0, 3},
			[]byte{0, 0, 1}},
		testrow{[]byte{5, 0, 0, 7},
			[]byte{0, 3},
			[]byte{1, 0, 1, 2}},
		testrow{[]byte{5, 0, 0, 7},
			[]byte{0, 3},
			[]byte{5, 0, 1, 7}},
	} {
		a := NewPolynomial(nil, row.a)
		b := NewPolynomial(nil, row.b)
		c := NewPolynomial(nil, row.c)
		checkAddAxioms(t, a, b, c, zero, add, eq)
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
