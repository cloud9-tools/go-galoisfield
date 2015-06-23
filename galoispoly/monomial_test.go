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
