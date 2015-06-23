package galoispoly

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/cloud9-tools/go-galoisfield"
)

var ErrIncompatibleFields = errors.New("cannot combine polynomials from different finite fields")

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

// Polynomial implements polynomials with coefficients drawn from a Galois field.
type Polynomial struct {
	field        *galoisfield.GF
	coefficients []byte
}

// NewPolynomial returns a new polynomial with the given coefficients.
// Coefficients are in little-endian order; that is, the first coefficient is
// the constant term, the second coefficient is the linear term, etc.
func NewPolynomial(field *galoisfield.GF, coefficients []byte) Polynomial {
	if field == nil {
		field = galoisfield.Default
	}
	for i := len(coefficients) - 1; i >= 0; i-- {
		if coefficients[i] != 0 {
			break
		}
		coefficients = coefficients[:i]
	}
	var dup []byte
	if len(coefficients) > 0 {
		dup = make([]byte, len(coefficients))
		copy(dup, coefficients)
	}
	return Polynomial{field, dup}
}

// Field returns the Galois field from which this polynomial's coefficients are drawn.
func (p Polynomial) Field() *galoisfield.GF { return p.field }

// IsZero returns true iff this polynomial has no terms.
func (p Polynomial) IsZero() bool { return p.coefficients == nil }

// Degree returns the degree of this polynomial, with the convention that the
// polynomial of zero terms has degree 0.
func (p Polynomial) Degree() uint {
	if p.IsZero() {
		return 0
	}
	return uint(len(p.coefficients) - 1)
}

// Coefficients returns the coefficients of the terms of this polynomial.  The
// result is in little-endian order; see NewPolynomial for details.
func (p Polynomial) Coefficients() []byte {
	var dup []byte
	if len(p.coefficients) > 0 {
		dup = make([]byte, len(p.coefficients))
		copy(dup, p.coefficients)
	}
	return dup
}

// Coefficient returns the coefficient of the i'th term.
func (p Polynomial) Coefficient(i uint) byte {
	if i >= uint(len(p.coefficients)) {
		return 0
	}
	return p.coefficients[i]
}

// Term returns the i'th term.
func (p Polynomial) Term(i uint) Monomial {
	return NewMonomial(p.field, p.Coefficient(i), i)
}

// Scale multiplies this polynomial by a scalar.
func (p Polynomial) Scale(s byte) Polynomial {
	if s == 0 {
		return Polynomial{p.field, nil}
	}
	if s == 1 {
		return p
	}
	coefficients := make([]byte, len(p.coefficients))
	for i, coeff_i := range p.coefficients {
		coefficients[i] = p.field.Mul(coeff_i, s)
	}
	return NewPolynomial(p.field, coefficients)
}

// Add returns the sum of one or more polynomials.
func Add(first Polynomial, rest ...Polynomial) Polynomial {
	sum := make([]byte, len(first.coefficients))
	copy(sum, first.coefficients)
	for _, next := range rest {
		if !galoisfield.Equal(first.field, next.field) {
			panic(ErrIncompatibleFields)
		}
		if next.IsZero() {
			continue
		}
		if len(next.coefficients) > len(sum) {
			newsum := make([]byte, len(next.coefficients))
			copy(newsum[:len(sum)], sum)
			sum = newsum
		}
		for i, ki := range next.coefficients {
			sum[i] = first.field.Add(sum[i], ki)
		}
	}
	return NewPolynomial(first.field, sum)
}

// Mul returns the product of one or more polynomials.
func Mul(first Polynomial, rest ...Polynomial) Polynomial {
	prod := make([]byte, len(first.coefficients))
	copy(prod, first.coefficients)
	for _, next := range rest {
		if !galoisfield.Equal(first.field, next.field) {
			panic(ErrIncompatibleFields)
		}
		if first.IsZero() || next.IsZero() {
			continue
		}
		p, q := prod, next.coefficients
		if len(p) < len(q) {
			p, q = q, p
		}
		newprod := make([]byte, len(p)+len(q)-1)
		//                                                 a3*x^3 + a2*x^2 + a1*x^1 + a0*x^0
		//                                               × b3*x^3 + b2*x^2 + b1*x^1 + b0*x^0
		// ---------------------------------------------------------------------------------
		// b3*x^3*(a3*x^3 + a2*x^2 + a1*x^1 + a0*x^0)
		//           + b2*x^2*(a3*x^3 + a2*x^2 + a1*x^1 + a0*x^0)
		//                       + b1*x^1*(a3*x^3 + a2*x^2 + a1*x^1 + a0*x^0)
		//                                   + b0*x^0*(a3*x^3 + a2*x^2 + a1*x^1 + a0*x^0)
		// ---------------------------------------------------------------------------------
		// a3*b3*x^6 + a2*b3*x^5 + a1*b3*x^4 + a0*b3*x^3
		//           + a3*b2*x^5 + a2*b2*x^4 + a1*b2*x^3 + a0*b2*x^2
		//                       + a3*b1*x^4 + a2*b1*x^3 + a1*b1*x^2 + a0*b1*x^1
		//                                   + a3*b0*x^3 + a2*b0*x^2 + a1*b0*x^1 + a0*b0*x^0
		// ---------------------------------------------------------------------------------
		// (a3*b3)x^6
		//           + (a2*b3+a3*b2)x^5
		//                       + (a1*b3+a2*b2+a3*b1)x^4
		//                                   + (a0*b3+a1*b2+a2*b1+a3*b0)x^3
		//                                               + (a0*b2+a1*b1+a2*b0)x^2
		//                                                           + (a0*b1+a1*b0)x^1
		//                                                                         + (a0*b0)
		for j := 0; j < len(q); j++ {
			for i := 0; i < len(p); i++ {
				product := first.field.Mul(p[i], q[j])
				newprod[i+j] = first.field.Add(newprod[i+j], product)
			}
		}
		prod = newprod
	}
	return NewPolynomial(first.field, prod)
}

// GoString returns a Go-syntax representation of this polynomial.
func (p Polynomial) GoString() string {
	return fmt.Sprintf("NewPolynomial(%#v, %#v)", p.field, p.coefficients)
}

// String returns a human-readable algebraic representation of this polynomial.
func (p Polynomial) String() string {
	if p.IsZero() {
		return "0"
	}
	var buf bytes.Buffer
	for i := len(p.coefficients) - 1; i >= 0; i-- {
		term := p.Term(uint(i))
		if !term.IsZero() {
			buf.WriteString(term.String())
			buf.WriteString(" + ")
		}
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 3)
	}
	return buf.String()
}

// Compare defines a partial order for polynomials: -1 if p < q, 0 if p == q,
// +1 if p > q, or panic if p and q are drawn from different Galois fields.
func (p Polynomial) Compare(q Polynomial) int {
	if !galoisfield.Equal(p.field, q.field) {
		panic(ErrIncompatibleFields)
	}
	switch {
	case len(p.coefficients) < len(q.coefficients):
		return -1
	case len(p.coefficients) > len(q.coefficients):
		return 1
	}
	for i := len(p.coefficients) - 1; i >= 0; i-- {
		pi := p.coefficients[i]
		qi := q.coefficients[i]
		if pi < qi {
			return -1
		}
		if pi > qi {
			return 1
		}
	}
	return 0
}

// Equals returns true iff p == q.
func (p Polynomial) Equals(q Polynomial) bool {
	if !galoisfield.Equal(p.field, q.field) {
		return false
	}
	return p.Compare(q) == 0
}

// Less returns true iff p < q.
func (p Polynomial) Less(q Polynomial) bool {
	return p.Compare(q) < 0
}
