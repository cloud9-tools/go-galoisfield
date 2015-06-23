package galoispoly

import "testing"

func checkCompareAxioms(t *testing.T, a, b interface{}, cmp int, lt, gt, eq, qe bool) {
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

func checkAddAxioms(t *testing.T, a, b, c, z interface{}, add func(_, _ interface{}) interface{}, eq func(_, _ interface{}) bool) {
	// a+0 = 0+a = a
	az := add(a, z)
	za := add(z, a)
	if !eq(az, za) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", a, z, az, za)
	} else if !eq(a, az) {
		t.Errorf("additive 'identity' isn't for %#v and %#v; got x+0=0+x=%#v", a, z, az)
	}

	// b+0 = 0+b = b
	bz := add(b, z)
	zb := add(z, b)
	if !eq(bz, zb) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", b, z, bz, zb)
	} else if !eq(b, bz) {
		t.Errorf("additive 'identity' isn't for %#v and %#v; got x+0=0+x=%#v", b, z, bz)
	}

	// c+0 = 0+c = b
	cz := add(c, z)
	zc := add(z, c)
	if !eq(cz, zc) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", c, z, cz, zc)
	} else if !eq(c, cz) {
		t.Errorf("additive 'identity' isn't for %#v and %#v; got x+0=0+x=%#v", c, z, cz)
	}

	// 0+0 = 0
	zz := add(z, z)
	if !eq(z, zz) {
		t.Errorf("additive 'identity' isn't for %#v and itself; got 0+0=%#v", z, zz)
	}

	// a+b = b+a
	ab := add(a, b)
	ba := add(b, a)
	if !eq(ab, ba) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", a, b, ab, ba)
	}
	// a+c = c+a
	ac := add(a, c)
	ca := add(c, a)
	if !eq(ac, ca) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", a, c, ac, ca)
	}
	// b+c = c+b
	bc := add(b, c)
	cb := add(c, b)
	if !eq(bc, cb) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", b, c, bc, cb)
	}

	// (a+b)+c = c+(a+b)
	ab_c := add(ab, c)
	c_ab := add(c, ab)
	if !eq(ab_c, c_ab) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", ab, c, ab_c, c_ab)
	}
	// a+(b+c) = (b+c)+a
	a_bc := add(a, bc)
	bc_a := add(bc, a)
	if !eq(a_bc, bc_a) {
		t.Errorf("addition not commutative for %#v and %#v; got x+y=%#v, y+x=%#v", a, bc, a_bc, bc_a)
	}
	// (a+b)+c = a+(b+c)
	if !eq(ab_c, a_bc) {
		t.Errorf("addition not associative; got (a+b)+c=%#v, a+(b+c)=%#v", ab_c, a_bc)
	}
}
