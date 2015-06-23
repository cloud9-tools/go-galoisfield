package galoispoly

import "testing"

func TestNewMonomial(t *testing.T) {
	type testrow struct {
		coefficient byte
		degree      uint
		expected    string
	}
	for _, row := range []testrow{
		testrow{0, 0, "0"},
		testrow{1, 0, "1"},
		testrow{2, 0, "2"},
		testrow{17, 0, "17"},
		testrow{0, 1, "0"},
		testrow{1, 1, "x"},
		testrow{2, 1, "2x"},
		testrow{17, 1, "17x"},
		testrow{0, 2, "0"},
		testrow{1, 2, "x^2"},
		testrow{2, 2, "2x^2"},
		testrow{17, 2, "17x^2"},
		testrow{0, 4, "0"},
		testrow{17, 4, "17x^4"},
		testrow{0, 5, "0"},
		testrow{1, 5, "x^5"},
	} {
		m := NewMonomial(nil, row.coefficient, row.degree)
		actual := m.String()
		if actual != row.expected {
			t.Errorf("%#v, got %q", row, actual)
		}
	}
}

func TestMonomial_Scale(t *testing.T) {
	type testrow struct {
		coefficient byte
		degree      uint
		scalar      byte
		expected    string
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
		m := NewMonomial(nil, row.coefficient, row.degree)
		m = m.Scale(row.scalar)
		actual := m.String()
		if actual != row.expected {
			t.Errorf("%#v, got %q", row, actual)
		}
	}
}

func TestMonomial_Compare(t *testing.T) {
	type testrow struct {
		a Monomial
		b Monomial
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
			t.Errorf("expected (%v)cmp(%v)=%d, got %d", a, b, expected, actual)
		}
	}
}

func TestMonomial_Add(t *testing.T) {
	type testrow struct {
		a Monomial
		b Monomial
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

func TestNewPolynomial(t *testing.T) {
	type testrow struct {
		coefficients []byte
		expected     string
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
		p := NewPolynomial(nil, row.coefficients)
		actual := p.String()
		if actual != row.expected {
			t.Errorf("%#v, got %q", row, actual)
		}
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
