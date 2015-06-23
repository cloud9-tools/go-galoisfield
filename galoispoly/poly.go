package galoispoly

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/cloud9-tools/go-galoisfield"
)

var ErrIncompatibleFields = errors.New("cannot combine polynomials from different finite fields")

type Monomial struct {
	field       *galoisfield.GF
	degree      uint
	coefficient byte
}

// NewMonomial returns coefficient*(x**degree) in field.
func NewMonomial(field *galoisfield.GF, coefficient byte, degree uint) Monomial {
	if coefficient == 0 {
		return Monomial{field, 0, 0}
	}
	return Monomial{field: field, degree: degree, coefficient: coefficient}
}

func (p Monomial) Field() *galoisfield.GF { return p.field }
func (p Monomial) Degree() uint           { return p.degree }
func (p Monomial) Coefficient() byte      { return p.coefficient }
func (p Monomial) IsZero() bool           { return p.coefficient == 0 }

func (p Monomial) Scale(s byte) Monomial {
	deg := p.degree
	coeff := p.field.Mul(p.coefficient, s)
	if coeff == 0 {
		deg = 0
	}
	return Monomial{field: p.field, degree: deg, coefficient: coeff}
}

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

func (p Monomial) GoString() string {
	return fmt.Sprintf("galois.NewMonomial(%#v, %d, %d)", p.field, p.coefficient, p.degree)
}

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

func (p Monomial) Polynomial() Polynomial {
	coefficients := make([]byte, p.degree+1)
	coefficients[p.degree] = p.coefficient
	return NewPolynomial(p.field, coefficients)
}

type Polynomial struct {
	field        *galoisfield.GF
	coefficients []byte
}

func NewPolynomial(field *galoisfield.GF, coefficients []byte) Polynomial {
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

func (p Polynomial) IsZero() bool           { return p.coefficients == nil }
func (p Polynomial) Field() *galoisfield.GF { return p.field }

func (p Polynomial) Degree() uint {
	if p.IsZero() {
		return 0
	}
	return uint(len(p.coefficients) - 1)
}

func (p Polynomial) Coefficients() []byte {
	var dup []byte
	if len(p.coefficients) > 0 {
		dup = make([]byte, len(p.coefficients))
		copy(dup, p.coefficients)
	}
	return dup
}

func (p Polynomial) Coefficient(i uint) byte {
	if i >= uint(len(p.coefficients)) {
		return 0
	}
	return p.coefficients[i]
}

func (p Polynomial) Term(i uint) Monomial {
	return NewMonomial(p.field, p.Coefficient(i), i)
}

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

func Add(p, q Polynomial) Polynomial {
	if !galoisfield.Equal(p.field, q.field) {
		panic(ErrIncompatibleFields)
	}
	if p.IsZero() || q.IsZero() {
		return Polynomial{p.field, nil}
	}
	d := maxint(len(p.coefficients), len(q.coefficients))
	coefficients := make([]byte, d)
	for i := range coefficients {
		a_i := p.Coefficient(uint(i))
		b_i := q.Coefficient(uint(i))
		coefficients[i] = p.field.Add(a_i, b_i)
	}
	return NewPolynomial(p.field, coefficients)
}

func Mul(p, q Polynomial) Polynomial {
	if !galoisfield.Equal(p.field, q.field) {
		panic(ErrIncompatibleFields)
	}
	if p.IsZero() || q.IsZero() {
		return Polynomial{p.field, nil}
	}
	coefficients := make([]byte, len(p.coefficients)+len(q.coefficients)-1)
	if p.Degree() < q.Degree() {
		p, q = q, p
	}
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
	for j := 0; j < len(q.coefficients); j++ {
		for i := 0; i <= len(p.coefficients); i++ {
			product := p.field.Mul(p.coefficients[i], q.coefficients[j])
			coefficients[i+j] = p.field.Add(coefficients[i+j], product)
		}
	}
	return NewPolynomial(p.field, coefficients)
}

func (p Polynomial) GoString() string {
	return fmt.Sprintf("galois.NewPolynomial(%#v, %#v)", p.field, p.coefficients)
}

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

func maxint(p, q int) int {
	if p > q {
		return p
	}
	return q
}
