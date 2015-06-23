package galoisfield

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrFieldSize      = errors.New("only field sizes 4, 8, 16, 32, 64, 128, and 256 are permitted")
	ErrPolyOutOfRange = errors.New("polynomial is out of range")
	ErrReduciblePoly  = errors.New("polynomial is reducible")
	ErrNotGenerator   = errors.New("value is not a generator")
	ErrDivByZero      = errors.New("division by zero")
	ErrLogZero        = errors.New("logarithm of zero")
)

var Poly84310_g3, Poly84320_g2, Default *GF

type params struct {
	n uint
	p uint
	g uint
}

type GF struct {
	params
	m uint
	log []byte
	exp []byte
}

var (
	mu     sync.Mutex
	global map[params]*GF
)

func init() {
	global = make(map[params]*GF)
	Poly84310_g3 = New(256, 0x11b, 0x03)
	Poly84320_g2 = New(256, 0x11d, 0x02)
	Default = Poly84320_g2
}

// New takes n (a power of 2), p (a polynomial), and g (a generator), then uses
// them to construct an instance of GF(n).  This comes complete with
// precomputed g**x and log_g(x) tables, so that all operations take O(1) time.
//
// If n isn't a supported power of 2, if p is reducible or of the wrong degree,
// or if g isn't actually a generator for the field, this function will panic.
//
// In the following, let k := log_2(n).
//
// The "p" argument describes a polynomial of the form
//	x**k + ∑_i: p_i*x**i; i ∈ [0..(k-1)]
// where the coefficient p_i is ((p>>i)&1), i.e. the i-th bit counting from the
// LSB.  The k-th bit MUST be 1, and all higher bits MUST be 0.
// Thus, n ≤ p < 2n.
//
// The "g" argument determines the permutation of field elements.  The value g
// chosen must be a generator for the field, i.e. the sequence
//	g**0, g**1, g**2, ... g**(n-1)
// must be a complete list of all elements in the field.  The field is small
// enough that the easiest way to discover generators is trial-and-error.
//
// The "p" and "g" arguments both have no effect on Add.
// The "g" argument additionally has no effect on (the output of) Mul/Div/Inv.
// Both arguments affect Exp/Log.
func New(n, p, g uint) *GF {
	switch n {
	case 4, 8, 16, 32, 64, 128, 256:
		// OK
	default:
		panic(ErrFieldSize)
	}
	m := n - 1
	if p < n || p >= 2*n {
		panic(ErrPolyOutOfRange)
	}
	if g == 0 || g == 1 {
		panic(ErrNotGenerator)
	}
	if isReducible(p) {
		panic(ErrReduciblePoly)
	}
	params := params{n, p, g}

	mu.Lock()
	singleton, found := global[params]
	mu.Unlock()
	if found {
		return singleton
	}

	gf := &GF{params, m, make([]byte, n), make([]byte, 2*n-2)}

	// Use the generator to compute the exp/log tables.  We perform the
	// usual trick of doubling the exp table to simplify Mul.
	var x uint = 1
	for i := uint(0); i < m; i++ {
		if x == 1 && i != 0 {
			panic(ErrNotGenerator)
		}
		gf.exp[i] = byte(x)
		gf.exp[i+m] = byte(x)
		gf.log[x] = byte(i)
		x = mulSlow(x, g, p, n)
	}

	mu.Lock()
	singleton, found = global[params]
	if !found {
		singleton = gf
		global[params] = singleton
	}
	mu.Unlock()
	return singleton
}

// Equal compares two GFs for equality.
func Equal(x, y *GF) bool {
	if x == nil {
		x = Default
	}
	if y == nil {
		y = Default
	}
	return x.params == y.params
}

// Less provides a total ordering over GFs.
func Less(x, y *GF) bool {
	if x == nil {
		x = Default
	}
	if y == nil {
		y = Default
	}
	if x.Size() != y.Size() {
		return x.Size() < y.Size()
	}
	if x.Polynomial() != y.Polynomial() {
		return x.Polynomial() < y.Polynomial()
	}
	return x.Generator() < y.Generator()
}

func (gf *GF) Size() uint       { return uint(len(gf.log)) }
func (gf *GF) Polynomial() uint { return gf.p }
func (gf *GF) Generator() uint  { return gf.g }

// Add returns x+y == x-y == x^y in GF(2**k).
func (_ *GF) Add(x, y byte) byte { return x ^ y }

// Sub returns x+y == x-y == x^y in GF(2**k).
func (_ *GF) Sub(x, y byte) byte { return x ^ y }

// Neg returns -x == x in GF(2**k).
func (_ *GF) Neg(x byte) byte { return x }

// Mul returns x*y in GF(2**k).
func (gf *GF) Mul(x, y byte) byte {
	if x == 0 || y == 0 {
		return 0
	}
	if gf == nil {
		gf = Default
	}
	return gf.exp[uint(gf.log[x])+uint(gf.log[y])]
}

// Div returns x/y in GF(2**k).
func (gf *GF) Div(x, y byte) byte {
	if x == 0 {
		return 0
	}
	if y == 0 {
		panic(ErrDivByZero)
	}
	if gf == nil {
		gf = Default
	}
	return gf.exp[gf.m+uint(gf.log[x])-uint(gf.log[y])]
}

// Inv returns 1/x in GF(2**k).
func (gf *GF) Inv(x byte) byte {
	if x == 0 {
		panic(ErrDivByZero)
	}
	if gf == nil {
		gf = Default
	}
	return gf.exp[gf.m-uint(gf.log[x])]
}

// Exp returns g**x in GF(2**k).
func (gf *GF) Exp(x byte) byte {
	if gf == nil {
		gf = Default
	}
	return gf.exp[uint(x)%gf.m]
}

// Log returns log_g(x) in GF(2**k).
func (gf *GF) Log(x byte) byte {
	if x == 0 {
		panic(ErrLogZero)
	}
	if gf == nil {
		gf = Default
	}
	return gf.log[x]
}

func (gf *GF) GoString() string {
	if gf == nil {
		gf = Default
	}
	if gf == Poly84310_g3 {
		return "Poly84310_g3"
	}
	if gf == Poly84320_g2 {
		return "Poly84320_g2"
	}
	return fmt.Sprintf("New(%d, %#x, %#x)", gf.Size(), gf.p, gf.g)
}

func (gf *GF) String() string {
	if gf == nil {
		gf = Default
	}
	return fmt.Sprintf("GF(%d;p=%#x;g=%#x)", gf.Size(), gf.p, gf.g)
}

// mulSlow returns x*y mod p.
func mulSlow(x, y, p, n uint) uint {
	r := uint(0)
	for x > 0 {
		if (x & 1) != 0 {
			r ^= y
		}
		x >>= 1
		y <<= 1
		if (y & n) != 0 {
			y ^= p
		}
	}
	return r
}

// isReducible returns true iff it can find a smaller polynomial that evenly
// divides the given polynomial.
func isReducible(p uint) bool {
	n := uint(1) << ((degree(p) / 2) + 1)
	for divisor := uint(2); divisor < n; divisor++ {
		if polyDiv(p, divisor) == 0 {
			return true
		}
	}
	return false
}

// polyDiv divides two polynomials and returns the remainder.
func polyDiv(dividend, divisor uint) uint {
	for m, n := degree(dividend), degree(divisor); m >= n; m-- {
		if (dividend & (1 << (m - 1))) != 0 {
			dividend ^= divisor << (m - n)
		}
	}
	return dividend
}

// degree returns the degree of the polynomial.  In this representation, the
// degree of a polynomial is:
//	[p == 0] 0
//	[p >  0] (k+1) such that (1<<k) is the highest 1 bit
func degree(p uint) uint {
	var d uint
	for p > 0 {
		d++
		p >>= 1
	}
	return d
}
