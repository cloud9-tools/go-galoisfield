package galoispoly

import (
	"fmt"
	"strconv"

	"github.com/cloud9-tools/go-galoisfield"
)

// Monomial implements monomials with coefficients drawn from a Galois field.
type Monomial struct {
	field       *galoisfield.GF
	degree      uint
	coefficient byte
}

// NewMonomial returns coefficient*(x**degree) in field.
func NewMonomial(field *galoisfield.GF, coefficient byte, degree uint) Monomial {
	if field == nil {
		field = galoisfield.Default
	}
	if coefficient == 0 {
		return Monomial{field, 0, 0}
	}
	return Monomial{field: field, degree: degree, coefficient: coefficient}
}

// Field returns the Galois field from which this monomial's coefficient is drawn.
func (a Monomial) Field() *galoisfield.GF { return a.field }

// Degree returns the degree of this monomial, with the convention that a
// monomial with a zero coefficient has degree 0.
func (a Monomial) Degree() uint { return a.degree }

// Coefficient returns the coefficient of this monomial.
func (a Monomial) Coefficient() byte { return a.coefficient }

// IsZero returns true iff this monomial has a zero coefficient.
func (a Monomial) IsZero() bool { return a.coefficient == 0 }

// Scale multiplies this monomial by a scalar.
func (a Monomial) Scale(s byte) Monomial {
	deg := a.degree
	coeff := a.field.Mul(a.coefficient, s)
	if coeff == 0 {
		deg = 0
	}
	return Monomial{field: a.field, degree: deg, coefficient: coeff}
}

// Mul multiplies this monomial by another monomial.
func (a Monomial) Mul(b Monomial) Monomial {
	if a.field != b.field {
		panic(ErrIncompatibleFields)
	}
	deg := a.degree + b.degree
	coeff := a.field.Mul(a.coefficient, b.coefficient)
	if coeff == 0 {
		deg = 0
	}
	return Monomial{field: a.field, degree: deg, coefficient: coeff}
}

// Polynomial returns the polynomial whose sole term is this monomial.
func (a Monomial) Polynomial() Polynomial {
	coefficients := make([]byte, a.degree+1)
	coefficients[a.degree] = a.coefficient
	return NewPolynomial(a.field, coefficients)
}

// GoString returns a Go-syntax representation of this monomial.
func (a Monomial) GoString() string {
	return fmt.Sprintf("NewMonomial(%#v, %d, %d)",
		a.field, a.coefficient, a.degree)
}

// String returns a human-readable algebraic representation of this monomial.
func (a Monomial) String() string {
	if a.IsZero() {
		return "0"
	} else if a.degree == 0 {
		return strconv.Itoa(int(a.coefficient))
	} else if a.degree == 1 && a.coefficient == 1 {
		return "x"
	} else if a.degree == 1 {
		return fmt.Sprintf("%dx", a.coefficient)
	} else if a.coefficient == 1 {
		return fmt.Sprintf("x^%d", a.degree)
	} else {
		return fmt.Sprintf("%dx^%d", a.coefficient, a.degree)
	}
}

// Compare defines a total order for monomials: -1 if a < b, 0 if a == b, or
// +1 if a > b.
func (a Monomial) Compare(b Monomial) int {
	if cmp := a.field.Compare(b.field); cmp != 0 {
		return cmp
	}
	if a.degree < b.degree {
		return -1
	}
	if a.degree > b.degree {
		return 1
	}
	if a.coefficient < b.coefficient {
		return -1
	}
	if a.coefficient > b.coefficient {
		return 1
	}
	return 0
}

// Equal returns true iff a == b.
func (a Monomial) Equal(b Monomial) bool {
	return a.Compare(b) == 0
}

// Less returns true iff a < b.
func (a Monomial) Less(b Monomial) bool {
	return a.Compare(b) < 0
}

// Evaluate substitutes for x and returns the resulting value.
func (a Monomial) Evaluate(x byte) byte {
	pow := a.coefficient
	for d := uint(0); d < a.degree; d++ {
		pow = a.field.Mul(pow, x)
	}
	return pow
}
