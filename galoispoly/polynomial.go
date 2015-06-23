package galoispoly

import (
	"bytes"
	"fmt"

	"github.com/cloud9-tools/go-galoisfield"
)

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
