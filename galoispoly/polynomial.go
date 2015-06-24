package galoispoly

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/cloud9-tools/go-galoisfield"
)

var (
	ErrIncompatibleFields = errors.New("cannot combine polynomials from different finite fields")
)

// Polynomial implements polynomials with coefficients drawn from a Galois field.
type Polynomial struct {
	field        *galoisfield.GF
	coefficients []byte
}

// NewPolynomial returns a new polynomial with the given coefficients.
// Coefficients are in little-endian order; that is, the first coefficient is
// the constant term, the second coefficient is the linear term, etc.
func NewPolynomial(field *galoisfield.GF, coefficients ...byte) Polynomial {
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
func (a Polynomial) Field() *galoisfield.GF { return a.field }

// IsZero returns true iff this polynomial has no terms.
func (a Polynomial) IsZero() bool { return a.coefficients == nil }

// Degree returns the degree of this polynomial, with the convention that the
// polynomial of zero terms has degree 0.
func (a Polynomial) Degree() uint {
	if a.IsZero() {
		return 0
	}
	return uint(len(a.coefficients) - 1)
}

// Coefficients returns the coefficients of the terms of this polynomial.  The
// result is in little-endian order; see NewPolynomial for details.
func (a Polynomial) Coefficients() []byte {
	var dup []byte
	if len(a.coefficients) > 0 {
		dup = make([]byte, len(a.coefficients))
		copy(dup, a.coefficients)
	}
	return dup
}

// Coefficient returns the coefficient of the i'th term.
func (a Polynomial) Coefficient(i uint) byte {
	if i >= uint(len(a.coefficients)) {
		return 0
	}
	return a.coefficients[i]
}

// Scale multiplies this polynomial by a scalar.
func (a Polynomial) Scale(s byte) Polynomial {
	if s == 0 {
		return Polynomial{a.field, nil}
	}
	if s == 1 {
		return a
	}
	coefficients := make([]byte, len(a.coefficients))
	for i, coeff_i := range a.coefficients {
		coefficients[i] = a.field.Mul(coeff_i, s)
	}
	return NewPolynomial(a.field, coefficients...)
}

// Add returns the sum of one or more polynomials.
func (first Polynomial) Add(rest ...Polynomial) Polynomial {
	sum := make([]byte, len(first.coefficients))
	copy(sum, first.coefficients)
	for _, next := range rest {
		if first.field != next.field {
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
	return NewPolynomial(first.field, sum...)
}

// Mul returns the product of one or more polynomials.
func (first Polynomial) Mul(rest ...Polynomial) Polynomial {
	prod := make([]byte, len(first.coefficients))
	copy(prod, first.coefficients)
	for _, next := range rest {
		if first.field != next.field {
			panic(ErrIncompatibleFields)
		}
		if first.IsZero() || next.IsZero() {
			continue
		}
		a, b := prod, next.coefficients
		if len(a) < len(b) {
			a, b = b, a
		}
		newprod := make([]byte, len(a)+len(b)-1)
		for bi := 0; bi < len(b); bi++ {
			for ai := 0; ai < len(a); ai++ {
				product := first.field.Mul(a[ai], b[bi])
				newprod[ai+bi] = first.field.Add(newprod[ai+bi], product)
			}
		}
		prod = newprod
	}
	return NewPolynomial(first.field, prod...)
}

// GoString returns a Go-syntax representation of this polynomial.
func (a Polynomial) GoString() string {
	var buf bytes.Buffer
	buf.WriteString("NewPolynomial(")
	buf.WriteString(a.field.GoString())
	for _, k := range a.coefficients {
		buf.WriteString(", ")
		buf.WriteString(strconv.Itoa(int(k)))
	}
	buf.WriteByte(')')
	return buf.String()
}

// String returns a human-readable algebraic representation of this polynomial.
func (a Polynomial) String() string {
	if a.IsZero() {
		return "0"
	}
	var buf bytes.Buffer
	for d := len(a.coefficients) - 1; d >= 0; d-- {
		k := a.coefficients[d]
		if k == 0 {
			continue
		}
		if buf.Len() > 0 {
			buf.WriteString(" + ")
		}
		if k > 1 || d == 0 {
			fmt.Fprintf(&buf, "%d", k)
		}
		if d > 1 {
			fmt.Fprintf(&buf, "x^%d", d)
		} else if d == 1 {
			buf.WriteByte('x')
		}
	}
	return buf.String()
}

// Compare defines a partial order for polynomials: -1 if a < b, 0 if a == b,
// +1 if a > b, or panic if a and b are drawn from different Galois fields.
func (a Polynomial) Compare(b Polynomial) int {
	if cmp := a.field.Compare(b.field); cmp != 0 {
		return cmp
	}
	if len(a.coefficients) < len(b.coefficients) {
		return -1
	}
	if len(a.coefficients) > len(b.coefficients) {
		return 1
	}
	for i := len(a.coefficients) - 1; i >= 0; i-- {
		pi := a.coefficients[i]
		qi := b.coefficients[i]
		if pi < qi {
			return -1
		}
		if pi > qi {
			return 1
		}
	}
	return 0
}

// Equal returns true iff a == b.
func (a Polynomial) Equal(b Polynomial) bool {
	return a.Compare(b) == 0
}

// Less returns true iff a < b.
func (a Polynomial) Less(b Polynomial) bool {
	return a.Compare(b) < 0
}

// Evaluate substitutes for x and returns the resulting value.
func (a Polynomial) Evaluate(x byte) byte {
	var sum byte = 0
	var pow byte = 1
	for _, k := range a.coefficients {
		sum = a.field.Add(sum, a.field.Mul(k, pow))
		pow = a.field.Mul(pow, x)
	}
	return sum
}
