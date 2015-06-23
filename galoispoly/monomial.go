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
func (p Monomial) Field() *galoisfield.GF { return p.field }

// Degree returns the degree of this monomial, with the convention that a
// monomial with a zero coefficient has degree 0.
func (p Monomial) Degree() uint { return p.degree }

// Coefficient returns the coefficient of this monomial.
func (p Monomial) Coefficient() byte { return p.coefficient }

// IsZero returns true iff this monomial has a zero coefficient.
func (p Monomial) IsZero() bool { return p.coefficient == 0 }

// Scale multiplies this monomial by a scalar.
func (p Monomial) Scale(s byte) Monomial {
	deg := p.degree
	coeff := p.field.Mul(p.coefficient, s)
	if coeff == 0 {
		deg = 0
	}
	return Monomial{field: p.field, degree: deg, coefficient: coeff}
}

// Mul multiplies this monomial by another monomial.
func (p Monomial) Mul(q Monomial) Monomial {
	if !galoisfield.Equal(p.field, q.field) {
		panic(ErrIncompatibleFields)
	}
	deg := p.degree + q.degree
	coeff := p.field.Mul(p.coefficient, q.coefficient)
	if coeff == 0 {
		deg = 0
	}
	return Monomial{field: p.field, degree: deg, coefficient: coeff}
}

// GoString returns a Go-syntax representation of this monomial.
func (p Monomial) GoString() string {
	return fmt.Sprintf("NewMonomial(%#v, %d, %d)",
		p.field, p.coefficient, p.degree)
}

// String returns a human-readable algebraic representation of this monomial.
func (p Monomial) String() string {
	if p.IsZero() {
		return "0"
	} else if p.degree == 0 {
		return strconv.Itoa(int(p.coefficient))
	} else if p.degree == 1 && p.coefficient == 1 {
		return "x"
	} else if p.degree == 1 {
		return fmt.Sprintf("%dx", p.coefficient)
	} else if p.coefficient == 1 {
		return fmt.Sprintf("x^%d", p.degree)
	} else {
		return fmt.Sprintf("%dx^%d", p.coefficient, p.degree)
	}
}

// Compare defines a partial order for monomials: -1 if p < q, 0 if p == q,
// +1 if p > q, or panic if p and q are drawn from different Galois fields.
func (p Monomial) Compare(q Monomial) int {
	if !galoisfield.Equal(p.field, q.field) {
		panic(ErrIncompatibleFields)
	}
	switch {
	case p.degree < q.degree:
		return -1
	case p.degree > q.degree:
		return 1
	case p.coefficient < q.coefficient:
		return -1
	case p.coefficient > q.coefficient:
		return 1
	default:
		return 0
	}
}

// Equals returns true iff p == q.
func (p Monomial) Equals(q Monomial) bool {
	if !galoisfield.Equal(p.field, q.field) {
		return false
	}
	return p.Compare(q) == 0
}

// Less returns true iff p < q.
func (p Monomial) Less(q Monomial) bool {
	return p.Compare(q) < 0
}

// Polynomial returns the polynomial whose sole term is this monomial.
func (p Monomial) Polynomial() Polynomial {
	coefficients := make([]byte, p.degree+1)
	coefficients[p.degree] = p.coefficient
	return NewPolynomial(p.field, coefficients)
}
