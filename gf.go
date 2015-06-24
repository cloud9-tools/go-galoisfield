package galoisfield

import (
	"errors"
	"fmt"
	"strings"
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

type params struct {
	p uint16
	k byte
	g byte
}

// GF represents a particular permutation of GF(2**k) for some fixed k.
type GF struct {
	params
	m   uint
	log []byte
	exp []byte
}

var (
	mu     sync.Mutex
	global map[params]*GF = make(map[params]*GF)
)

// Some handy pre-chosen polynomial/generator combinations.
var (
	// GF(4) p=(x^2 + x + 1) g=2
	Poly210_g2 = New(4, 0x7, 2)

	// GF(8) p=(x^3 + x + 1) g=2
	Poly310_g2 = New(8, 0xb, 2)

	// GF(16) p=(x^4 + x + 1) g=2
	Poly410_g2 = New(16, 0x13, 2)

	// GF(32) p=(x^5 + x^2 + 1) g=2
	Poly520_g2 = New(32, 0x25, 2)

	// GF(64) p=(x^6 + x + 1) g=2
	Poly610_g2 = New(64, 0x43, 2)
	// GF(64) p=(x^6 + x + 1) g=7
	Poly610_g7 = New(64, 0x43, 7)

	// GF(128) p=(x^7 + x + 1) g=2
	Poly710_g2 = New(128, 0x83, 2)

	// GF(256), p (x^8 + x^4 + x^3 + x + 1), g 3
	Poly84310_g3 = New(256, 0x11b, 0x03)
	// GF(256), p (x^8 + x^4 + x^3 + x^2 + 1), g 2
	Poly84320_g2 = New(256, 0x11d, 0x02)

	// Some arbitrarily-chosen permutations of GF(n).
	DefaultGF4   = Poly210_g2
	DefaultGF8   = Poly310_g2
	DefaultGF16  = Poly410_g2
	DefaultGF32  = Poly520_g2
	DefaultGF64  = Poly610_g2
	DefaultGF128 = Poly710_g2
	DefaultGF256 = Poly84320_g2

	// Some arbitrarily-chosen permutation of GF(256).
	Default = DefaultGF256
)

type wki struct {
	field *GF
	name string
}

var wellknown = []wki{
	wki{nil, "nil"},
	wki{Poly210_g2, "Poly210_g2"},
	wki{Poly310_g2, "Poly310_g2"},
	wki{Poly410_g2, "Poly410_g2"},
	wki{Poly520_g2, "Poly520_g2"},
	wki{Poly610_g2, "Poly610_g2"},
	wki{Poly610_g7, "Poly610_g7"},
	wki{Poly710_g2, "Poly710_g2"},
	wki{Poly84310_g3, "Poly84310_g3"},
	wki{Poly84320_g2, "Poly84320_g2"},
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
//
//	x**k + ∑_i: p_i*x**i; i ∈ [0..(k-1)]
//
// where the coefficient p_i is ((p>>i)&1), i.e. the i-th bit counting from the
// LSB.  The k-th bit MUST be 1, and all higher bits MUST be 0.
// Thus, n ≤ p < 2n.
//
// The "g" argument determines the permutation of field elements.  The value g
// chosen must be a generator for the field, i.e. the sequence
//
//	g**0, g**1, g**2, ... g**(n-1)
//
// must be a complete list of all elements in the field.  The field is small
// enough that the easiest way to discover generators is trial-and-error.
//
// The "p" and "g" arguments both have no effect on Add.
// The "g" argument additionally has no effect on (the output of) Mul/Div/Inv.
// Both arguments affect Exp/Log.
func New(n, p uint, g byte) *GF {
	k, ok := log2table[n]
	if !ok {
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
	params := params{
		p: uint16(p),
		k: k,
		g: g,
	}

	mu.Lock()
	singleton, found := global[params]
	mu.Unlock()
	if found {
		return singleton
	}

	gf := &GF{
		params: params,
		m: m,
		log: make([]byte, n),
		exp: make([]byte, 2*n-2),
	}

	// Use the generator to compute the exp/log tables.  We perform the
	// usual trick of doubling the exp table to simplify Mul.
	var x byte = 1
	for i := uint(0); i < m; i++ {
		if x == 1 && i != 0 {
			panic(ErrNotGenerator)
		}
		gf.exp[i] = x
		gf.exp[i+m] = x
		gf.log[x] = byte(i)
		x = mulSlow(x, g, byte(p), k)
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

// Size returns the order of the Galois field, i.e. the number of elements.
func (gf *GF) Size() uint { return 1 << gf.k }

// Polynomial returns the polynomial used to generate the Galois field.
func (gf *GF) Polynomial() uint { return uint(gf.p) }

// Generator returns the exponent base used to generate the Galois field.
func (gf *GF) Generator() uint { return uint(gf.g) }

// Compare defines a total order for finite fields: -1 if a < b, 0 if a == b,
// or +1 if a > b.
func (a *GF) Compare(b *GF) int {
	switch {
	case a.k < b.k:
		return -1
	case a.k > b.k:
		return 1
	case a.p < b.p:
		return -1
	case a.p > b.p:
		return 1
	case a.g < b.g:
		return -1
	case a.g > b.g:
		return 1
	default:
		return 0
	}
}

// Equal returns true iff a == b.
func (a *GF) Equal(b *GF) bool {
	return a.Compare(b) == 0
}

// Less returns true iff a < b.
func (a *GF) Less(b *GF) bool {
	return a.Compare(b) < 0
}

// GoString returns a Go-syntax representation of this GF.
func (gf *GF) GoString() string {
	for _, wk := range wellknown {
		if gf == wk.field {
			return wk.name
		}
	}
	return fmt.Sprintf("New(%d, %#x, %d)", 1 << gf.k, gf.p, gf.g)
}

// String returns a human-readable representation of this GF.
func (gf *GF) String() string {
	if gf == nil {
		return "<nil>"
	}
	var poly []string
	for i := 8; i >= 0; i-- {
		if (gf.p & (1 << uint(i))) != 0 {
			var mono string
			if i == 0 {
				mono = "1"
			} else if i == 1 {
				mono = "b"
			} else {
				mono = fmt.Sprintf("b^%d", i)
			}
			poly = append(poly, mono)
		}
	}
	polystr := strings.Join(poly, "+")
	return fmt.Sprintf("GF(%d;%s;%d)", 1 << gf.k, polystr, gf.g)
}

// Add returns x+y == x-y == x^y in GF(2**k).
func (_ *GF) Add(x, y byte) byte { return x ^ y }

// Mul returns x*y in GF(2**k).
func (gf *GF) Mul(x, y byte) byte {
	if x == 0 || y == 0 {
		return 0
	}
	return gf.exp[uint(gf.log[x])+uint(gf.log[y])]
}

// Div returns x/y in GF(2**k).
func (gf *GF) Div(x, y byte) byte {
	if x == 0 || y == 0 {
		if y == 0 {
			panic(ErrDivByZero)
		}
		return 0
	}
	return gf.exp[gf.m+uint(gf.log[x])-uint(gf.log[y])]
}

// Inv returns 1/x in GF(2**k).
func (gf *GF) Inv(x byte) byte {
	if x == 0 {
		panic(ErrDivByZero)
	}
	return gf.exp[gf.m-uint(gf.log[x])]
}

// Exp returns g**x in GF(2**k).
func (gf *GF) Exp(x byte) byte {
	return gf.exp[uint(x)%gf.m]
}

// Log returns log_g(x) in GF(2**k).
func (gf *GF) Log(x byte) byte {
	if x == 0 {
		panic(ErrLogZero)
	}
	return gf.log[x]
}

// mulSlow returns x*y mod poly.
func mulSlow(x, y, poly, k byte) byte {
	var hibit byte = (1 << (k - 1))
	var p byte = 0
	for i := uint(0); i < uint(k); i++ {
		if (y & 1) != 0 {
			p ^= x
		}
		wasset := (x & hibit) != 0
		x <<= 1
		y >>= 1
		if wasset {
			x ^= poly
		}
	}
	return p
}

// isReducible returns true iff it can find a smaller polynomial that evenly
// divides the given polynomial.
func isReducible(p uint) bool {
	var n uint = 1 << ((degree(p) / 2) + 1)
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

var log2table = map[uint]byte{2: 1, 4: 2, 8: 3, 16: 4, 32: 5, 64: 6, 128: 7, 256: 8}
