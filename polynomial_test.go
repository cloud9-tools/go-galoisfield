package galoisfield

import (
	"math/rand"
	"testing"
)

func TestNewPolynomial(t *testing.T) {
	type testrow struct {
		input Polynomial
		str   string
		gostr string
		field *GF
		deg   uint
		coeff []byte
	}
	for idx, row := range []testrow{
		testrow{NewPolynomial(nil),
			"0",
			"NewPolynomial(Poly84320_g2)",
			Poly84320_g2, 0, nil},
		testrow{NewPolynomial(nil, 1),
			"1",
			"NewPolynomial(Poly84320_g2, 1)",
			Poly84320_g2, 0, []byte{1}},
		testrow{NewPolynomial(nil, 2),
			"2",
			"NewPolynomial(Poly84320_g2, 2)",
			Poly84320_g2, 0, []byte{2}},
		testrow{NewPolynomial(nil, 17),
			"17",
			"NewPolynomial(Poly84320_g2, 17)",
			Poly84320_g2, 0, []byte{17}},
		testrow{NewPolynomial(nil, 0, 2),
			"2x",
			"NewPolynomial(Poly84320_g2, 0, 2)",
			Poly84320_g2, 1, []byte{0, 2}},
		testrow{NewPolynomial(nil, 1, 2),
			"2x + 1",
			"NewPolynomial(Poly84320_g2, 1, 2)",
			Poly84320_g2, 1, []byte{1, 2}},
		testrow{NewPolynomial(nil, 1, 0, 1),
			"x^2 + 1",
			"NewPolynomial(Poly84320_g2, 1, 0, 1)",
			Poly84320_g2, 2, []byte{1, 0, 1}},
		testrow{NewPolynomial(nil, 0, 1, 1),
			"x^2 + x",
			"NewPolynomial(Poly84320_g2, 0, 1, 1)",
			Poly84320_g2, 2, []byte{0, 1, 1}},
		testrow{NewPolynomial(nil, 0, 1, 1, 0),
			"x^2 + x",
			"NewPolynomial(Poly84320_g2, 0, 1, 1)",
			Poly84320_g2, 2, []byte{0, 1, 1}},
		testrow{NewPolynomial(nil, 3, 1, 4),
			"4x^2 + x + 3",
			"NewPolynomial(Poly84320_g2, 3, 1, 4)",
			Poly84320_g2, 2, []byte{3, 1, 4}},
	} {
		str := row.input.String()
		if str != row.str {
			t.Errorf("[%2d] expected %q, got %q", idx, row.str, str)
		}
		gostr := row.input.GoString()
		if gostr != row.gostr {
			t.Errorf("[%2d] expected %q, got %q", idx, row.gostr, gostr)
		}
		field := row.input.Field()
		if field != row.field {
			t.Errorf("[%2d] expected %#v, got %#v", idx, row.field, field)
		}
		deg := row.input.Degree()
		if deg != row.deg {
			t.Errorf("[%2d] expected %d, got %d", idx, row.deg, deg)
		}
		coeff := row.input.Coefficients()
		if !equalBytes(coeff, row.coeff) {
			t.Errorf("[%2d] expected %v, got %v", idx, row.coeff, coeff)
		}
		for i, k := range row.coeff {
			actual := row.input.Coefficient(uint(i))
			if actual != k {
				t.Errorf("[%2d] expected %d, got %d", idx, k, actual)
			}
		}
		for i := 0; i < len(row.coeff); i++ {
			actual := row.input.Coefficient(uint(i + len(row.coeff)))
			if actual != 0 {
				t.Errorf("[%2d] expected 0, got %d", idx, actual)
			}
		}
	}
}

func TestPolynomial_Scale(t *testing.T) {
	type testrow struct {
		scalar   byte
		input    Polynomial
		expected Polynomial
	}
	for _, row := range []testrow{
		testrow{5,
			NewPolynomial(nil, 3, 0, 1),
			NewPolynomial(nil, 15, 0, 5)},
		testrow{1,
			NewPolynomial(nil, 3, 0, 1),
			NewPolynomial(nil, 3, 0, 1)},
		testrow{0,
			NewPolynomial(nil, 3, 0, 1),
			NewPolynomial(nil)},
	} {
		actual := row.input.Scale(row.scalar)
		if !actual.Equal(row.expected) {
			t.Errorf("expected %d*(%v)=(%v), got (%v)",
				row.scalar, row.input, row.expected, actual)
		}
	}
}

func TestPolynomial_Compare(t *testing.T) {
	type testrow struct {
		a        Polynomial
		b        Polynomial
		expected int
	}
	for _, row := range []testrow{
		testrow{NewPolynomial(nil), NewPolynomial(nil), 0},
		testrow{NewPolynomial(nil, 5), NewPolynomial(nil, 5), 0},
		testrow{NewPolynomial(nil, 3, 5), NewPolynomial(nil, 3, 5), 0},
		testrow{NewPolynomial(nil), NewPolynomial(nil, 1), -1},
		testrow{NewPolynomial(nil, 0), NewPolynomial(nil, 1), -1},
		testrow{NewPolynomial(nil, 2, 1), NewPolynomial(nil, 1, 2), -1},
		testrow{NewPolynomial(Poly310_g2),
			NewPolynomial(Poly210_g2), 1},
	} {
		a, b, expected := row.a, row.b, row.expected
		actual := a.Compare(b)
		if actual != expected {
			t.Errorf(
				"expected %#v cmp %#v == %d, got %d",
				a, b, expected, actual)
		}
		checkCompareAxioms(
			t, a, b, actual,
			a.Less(b),
			b.Less(a),
			a.Equal(b),
			b.Equal(a))
	}
}

func TestPolynomial_Add(t *testing.T) {
	type testrow struct {
		a, b     Polynomial
		expected Polynomial
	}
	for _, row := range []testrow{
		testrow{NewPolynomial(nil, 1, 0, 0, 1),
			NewPolynomial(nil),
			NewPolynomial(nil, 1, 0, 0, 1)},
		testrow{NewPolynomial(nil, 1, 0, 0, 1),
			NewPolynomial(nil, 0, 1),
			NewPolynomial(nil, 1, 1, 0, 1)},
		testrow{NewPolynomial(nil, 1, 0, 0, 1),
			NewPolynomial(nil, 0, 0, 1, 1),
			NewPolynomial(nil, 1, 0, 1)},
	} {
		actual := row.a.Add(row.b)
		if !actual.Equal(row.expected) {
			t.Errorf(
				"expected (%v)+(%v)=(%v), got %v",
				row.a, row.b, row.expected, actual)
		}
	}
}

func TestPolynomial_axioms(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	for _, field := range []*GF{
		Default,
	} {
		zero := NewPolynomial(field)
		one := NewPolynomial(field, 1)
		for deg := 0; deg <= 8; deg++ {
			for trial := 0; trial < 256; trial++ {
				var a, b, c []byte
				zipf := rand.NewZipf(prng, 2.0, 1.0, uint64(deg))
				bdeg := deg - int(zipf.Uint64())
				cdeg := deg - int(zipf.Uint64())
				for j := 0; j <= deg; j++ {
					a = append(a, byte(prng.Intn(int(field.Size()))))
				}
				for j := 0; j <= bdeg; j++ {
					b = append(b, byte(prng.Intn(int(field.Size()))))
				}
				for j := 0; j <= cdeg; j++ {
					c = append(c, byte(prng.Intn(int(field.Size()))))
				}
				checkAddAxioms(
					t,
					NewPolynomial(field, a...),
					NewPolynomial(field, b...),
					NewPolynomial(field, c...),
					zero)
				checkMulAxioms(
					t,
					NewPolynomial(field, a...),
					NewPolynomial(field, b...),
					NewPolynomial(field, c...),
					zero, one)
			}
		}
	}
}

func TestPolynomial_Add_incompatible(t *testing.T) {
	e := panicValue(func() {
		_ = NewPolynomial(Poly210_g2).
			Add(NewPolynomial(Poly310_g2))
	})
	if e != ErrIncompatibleFields {
		t.Errorf("expected ErrIncompatibleFields, got %q", e.Error())
	}
}

func TestPolynomial_Mul_incompatible(t *testing.T) {
	e := panicValue(func() {
		_ = NewPolynomial(Poly210_g2).
			Mul(NewPolynomial(Poly310_g2))
	})
	if e != ErrIncompatibleFields {
		t.Errorf("expected ErrIncompatibleFields, got %q", e.Error())
	}
}

func checkCompareAxioms(t *testing.T, a, b Polynomial, cmp int, lt, gt, eq, qe bool) {
	if eq != qe {
		t.Errorf("equality not commutative for %#v and %#v", a, b)
	}
	if eq && lt {
		t.Errorf("equality and lessthan not disjoint for %#v and %#v", a, b)
	}
	if eq && gt {
		t.Errorf("equality and greaterthan not disjoint for %#v and %#v", a, b)
	}
	if cmp < 0 && !lt {
		t.Errorf("expected %#v < %#v, got >=", a, b)
	}
	if cmp == 0 && !eq {
		t.Errorf("expected %#v == %#v, got !=", a, b)
	}
	if cmp > 0 && !gt {
		t.Errorf("expected %#v > %#v, got <=", a, b)
	}
}

func checkAddAxioms(t *testing.T, a, b, c, zero Polynomial) {
	// 0+0 = 0
	zz := zero.Add(zero)
	if !zero.Equal(zz) {
		t.Errorf(
			"additive 'identity' isn't for %#v and itself: "+
			"got 0+0=%#v", zero, zz)
	}

	// a+0 = 0+a = a
	az := a.Add(zero)
	za := zero.Add(a)
	if !az.Equal(za) {
		t.Errorf(
			"addition not commutative for a=%#v: "+
			"got a+0=%#v vs 0+a=%#v", a, az, za)
	}
	if !az.Equal(a) {
		t.Errorf(
			"additive 'identity' isn't for a=%#v: "+
			"got a+0=%#v", a, az)
	}
	// b+0 = 0+b = b
	bz := b.Add(zero)
	zb := zero.Add(b)
	if !bz.Equal(zb) {
		t.Errorf(
			"addition not commutative for b=%#v: "+
			"got b+0=%#v vs 0+b=%#v", b, bz, zb)
	}
	if !bz.Equal(b) {
		t.Errorf(
			"additive 'identity' isn't for b=%#v: "+
			"got b+0=%#v", b, bz)
	}
	// c+0 = 0+c = c
	cz := c.Add(zero)
	zc := zero.Add(c)
	if !cz.Equal(zc) {
		t.Errorf(
			"addition not commutative for c=%#v: "+
			"got c+0=%#v vs 0+c=%#v", c, cz, zc)
	}
	if !cz.Equal(c) {
		t.Errorf(
			"additive 'identity' isn't for c=%#v: "+
			"got c+0=%#v", c, cz)
	}

	// a+b = b+a
	ab := a.Add(b)
	ba := b.Add(a)
	if !ab.Equal(ba) {
		t.Errorf(
			"addition not commutative for a=%#v b=%#v: "+
			"got a+b=%#v vs b+a=%#v", a, b, ab, ba)
	}
	// a+c = c+a
	ac := a.Add(c)
	ca := c.Add(a)
	if !ac.Equal(ca) {
		t.Errorf(
			"addition not commutative for a=%#v c=%#v: "+
			"got a+c=%#v, c+a=%#v", a, c, ac, ca)
	}
	// b+c = c+b
	bc := b.Add(c)
	cb := c.Add(b)
	if !bc.Equal(cb) {
		t.Errorf(
			"addition not commutative for b=%#v c=%#v: "+
			"got b+c=%#v, c+b=%#v", b, c, bc, cb)
	}

	type item struct {
		name  string
		value Polynomial
	}

	// a+(b+c) = (a+b)+c etc.
	list := []item{
		item{"a+(b+c)", a.Add(b.Add(c))},
		item{"a+(c+b)", a.Add(c.Add(b))},
		item{"b+(a+c)", b.Add(a.Add(c))},
		item{"b+(c+a)", b.Add(c.Add(a))},
		item{"c+(a+b)", c.Add(a.Add(b))},
		item{"c+(b+a)", c.Add(b.Add(a))},
	}
	for x, xitem := range list {
		for y := x + 1; y < len(list); y++ {
			yitem := list[y]
			if !xitem.value.Equal(yitem.value) {
				t.Errorf(
					"addition not associative for "+
						"a=%#v b=%#v c=%#v: "+
						"got %s=%#v vs %s=%#v",
					a, b, c,
					xitem.name, xitem.value,
					yitem.name, yitem.value)
			}
		}
	}
}

func checkMulAxioms(t *testing.T, a, b, c, zero, one Polynomial) {
	// 1*1 = 1
	oo := one.Mul(one)
	if !oo.Equal(one) {
		t.Errorf(
			"multiplicative 'identity' isn't for %#v and itself: "+
				"got 1*1=%#v", one, oo)
	}

	// a*0 = 0*a = 0
	az := a.Mul(zero)
	za := zero.Mul(a)
	if !az.Equal(zero) {
		t.Errorf(
			"multiplicative 'zero' isn't for a=%#v: "+
				"got a*0=%#v", a, az)
	}
	if !az.Equal(za) {
		t.Errorf(
			"multiplication by 0 not commutative for a=%#v: "+
				"got a*0=%#v vs 0*a=%#v", a, az, za)
	}
	// b*0 = 0*b = 0
	bz := b.Mul(zero)
	zb := zero.Mul(b)
	if !bz.Equal(zero) {
		t.Errorf(
			"multiplicative 'zero' isn't for b=%#v: "+
				"got b*0=%#v", b, bz)
	}
	if !bz.Equal(zb) {
		t.Errorf(
			"multiplication by 0 not commutative for b=%#v: "+
				"got b*0=%#v vs 0*b=%#v", b, bz, zb)
	}
	// c*0 = 0*c = 0
	cz := c.Mul(zero)
	zc := zero.Mul(c)
	if !cz.Equal(zero) {
		t.Errorf(
			"multiplicative 'zero' isn't for c=%#v: "+
				"got c*0=%#v", c, cz)
	}
	if !cz.Equal(zc) {
		t.Errorf(
			"multiplication by 0 not commutative for c=%#v: "+
				"got c*0=%#v vs 0*c=%#v", c, cz, zc)
	}

	// a*1 = 1*a = a
	ao := a.Mul(one)
	oa := one.Mul(a)
	if !ao.Equal(a) {
		t.Errorf(
			"multiplicative 'identity' isn't for a=%#v: "+
				"got x*0=0*x=%#v", a, ao)
	}
	if !ao.Equal(oa) {
		t.Errorf(
			"multiplication by 1 not commutative for a=%#v: "+
				"got a*1=%#v vs 1*a=%#v", a, ao, oa)
	}
	// b*1 = 1*b = b
	bo := b.Mul(one)
	ob := one.Mul(b)
	if !bo.Equal(b) {
		t.Errorf(
			"multiplicative 'identity' isn't for b=%#v: "+
				"got x*0=0*x=%#v", b, bo)
	}
	if !bo.Equal(ob) {
		t.Errorf(
			"multiplication by 1 not commutative for b=%#v: "+
				"got b*1=%#v vs 1*b=%#v", b, bo, ob)
	}
	// c*1 = 1*c = c
	co := c.Mul(one)
	oc := one.Mul(c)
	if !co.Equal(c) {
		t.Errorf(
			"multiplicative 'identity' isn't for c=%#v: "+
				"got x*0=0*x=%#v", c, co)
	}
	if !co.Equal(oc) {
		t.Errorf(
			"multiplication by 1 not commutative for c=%#v: "+
				"got c*1=%#v vs 1*c=%#v", c, co, oc)
	}

	// a*b = b*a
	ab := a.Mul(b)
	ba := b.Mul(a)
	if !ab.Equal(ba) {
		t.Errorf(
			"multiplication not commutative for a=%#v b=%#v: "+
				"got a*b=%#v vs b*a=%#v", a, b, ab, ba)
	}
	// a*c = c*a
	ac := a.Mul(c)
	ca := c.Mul(a)
	if !ac.Equal(ca) {
		t.Errorf(
			"multiplication not commutative for a=%#v c=%#v: "+
				"got a*c=%#v vs c*a=%#v", a, c, ac, ca)
	}
	// b*c = c*b
	bc := b.Mul(c)
	cb := c.Mul(b)
	if !bc.Equal(cb) {
		t.Errorf(
			"multiplication not commutative for b=%#v c=%#v: "+
				"got b*c=%#v vs c*b=%#v", b, c, bc, cb)
	}

	type item struct {
		name  string
		value Polynomial
	}

	// a*(b*c) = (a*b)*c etc.
	list := []item{
		item{"a*(b*c)", a.Mul(b.Mul(c))},
		item{"a*(c*b)", a.Mul(c.Mul(b))},
		item{"b*(a*c)", b.Mul(a.Mul(c))},
		item{"b*(c*a)", b.Mul(c.Mul(a))},
		item{"c*(a*b)", c.Mul(a.Mul(b))},
		item{"c*(b*a)", c.Mul(b.Mul(a))},
	}
	for x, xitem := range list {
		for y := x + 1; y < len(list); y++ {
			yitem := list[y]
			if !xitem.value.Equal(yitem.value) {
				t.Errorf(
					"multiplication not associative for "+
						"a=%#v b=%#v c=%#v: "+
						"got %s=%#v vs %s=%#v",
					a, b, c,
					xitem.name, xitem.value,
					yitem.name, yitem.value)
			}
		}
	}

	// a*(b+c) = a*b + a*c
	ax_bpc := a.Mul(b.Add(c))
	axbp_axc := a.Mul(b).Add(a.Mul(c))
	if !ax_bpc.Equal(axbp_axc) {
		t.Errorf(
			"multiplication not distributive for "+
				"a=%#v b=%#v c=%#v: "+
				"got a*(b+c)=%#v vs a*b+a*c=%#v",
			a, b, c, ax_bpc, axbp_axc)
	}
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
