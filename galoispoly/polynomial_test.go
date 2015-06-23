package galoispoly

import (
	"testing"
)

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
