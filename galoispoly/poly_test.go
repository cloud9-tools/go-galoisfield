package galoispoly

import "testing"

func TestNewMonomial(t *testing.T) {
	type testrow struct {
		coefficient byte
		degree uint
		expected string
	}
	for _, row := range []testrow{
		testrow{ 0, 0, "0"},
		testrow{ 1, 0, "1"},
		testrow{ 2, 0, "2"},
		testrow{17, 0, "17"},
		testrow{ 0, 1, "0"},
		testrow{ 1, 1, "x"},
		testrow{ 2, 1, "2x"},
		testrow{17, 1, "17x"},
		testrow{ 0, 2, "0"},
		testrow{ 1, 2, "x^2"},
		testrow{ 2, 2, "2x^2"},
		testrow{17, 2, "17x^2"},
		testrow{ 0, 4, "0"},
		testrow{17, 4, "17x^4"},
		testrow{ 0, 5, "0"},
		testrow{ 1, 5, "x^5"},
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
		degree uint
		scalar byte
		expected string
	}
	for _, row := range []testrow{
		testrow{ 0, 0, 0, "0"},
		testrow{ 1, 0, 0, "0"},
		testrow{ 1, 1, 0, "0"},
		testrow{ 1, 5, 0, "0"},

		testrow{ 0, 0, 1, "0"},
		testrow{ 1, 0, 1, "1"},
		testrow{ 1, 1, 1, "x"},
		testrow{ 1, 5, 1, "x^5"},

		testrow{ 0, 0, 2, "0"},
		testrow{ 1, 0, 2, "2"},
		testrow{ 1, 1, 2, "2x"},
		testrow{ 1, 5, 2, "2x^5"},
	} {
		m := NewMonomial(nil, row.coefficient, row.degree)
		m = m.Scale(row.scalar)
		actual := m.String()
		if actual != row.expected {
			t.Errorf("%#v, got %q", row, actual)
		}
	}
}

func TestNewPolynomial(t *testing.T) {
	type testrow struct {
		coefficients []byte
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
		p := NewPolynomial(nil, row.coefficients)
		actual := p.String()
		if actual != row.expected {
			t.Errorf("%#v, got %q", row, actual)
		}
	}
}
