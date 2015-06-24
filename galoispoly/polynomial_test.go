package galoispoly

import (
	"testing"
)

func TestNewPolynomial(t *testing.T) {
	type testrow struct {
		input Polynomial
		str string
		gostr string
	}
	for idx, row := range []testrow{
		testrow{NewPolynomial(nil),
			"0",
			"NewPolynomial(Poly84320_g2)"},
		testrow{NewPolynomial(nil, 1),
			"1",
			"NewPolynomial(Poly84320_g2, 1)"},
		testrow{NewPolynomial(nil, 2),
			"2",
			"NewPolynomial(Poly84320_g2, 2)"},
		testrow{NewPolynomial(nil, 17),
			"17",
			"NewPolynomial(Poly84320_g2, 17)"},
		testrow{NewPolynomial(nil, 0, 2),
			"2x",
			"NewPolynomial(Poly84320_g2, 0, 2)"},
		testrow{NewPolynomial(nil, 1, 2),
			"2x + 1",
			"NewPolynomial(Poly84320_g2, 1, 2)"},
		testrow{NewPolynomial(nil, 1, 0, 1),
			"x^2 + 1",
			"NewPolynomial(Poly84320_g2, 1, 0, 1)"},
		testrow{NewPolynomial(nil, 0, 1, 1),
			"x^2 + x",
			"NewPolynomial(Poly84320_g2, 0, 1, 1)"},
		testrow{NewPolynomial(nil, 3, 1, 4),
			"4x^2 + x + 3",
			"NewPolynomial(Poly84320_g2, 3, 1, 4)"},
	} {
		str := row.input.String()
		if str != row.str {
			t.Errorf("[%2d] expected %q, got %q", idx, row.str, str)
		}
		gostr := row.input.GoString()
		if gostr != row.gostr {
			t.Errorf("[%2d] expected %q, got %q", idx, row.gostr, gostr)
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

func TestAdd(t *testing.T) {
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

func TestAdd_axioms(t *testing.T) {
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
