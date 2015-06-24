package galoisfield

import (
	"math/rand"
	"testing"
)

func TestNew_gf256(t *testing.T) {
	for _, gf := range []*GF{Poly84310_g3, Poly84320_g2} {
		for i := byte(0); i < 255; i++ {
			x1 := gf.exp[i]
			y1 := gf.log[x1]
			if y1 != i {
				t.Errorf("expected log[exp[i]]=i, got i=%d exp[.]=%d log[.]=%d",
					i, x1, y1)
			}
			x2 := gf.exp[int(i)+255]
			y2 := gf.log[x2]
			if y2 != i {
				t.Errorf("expected log[exp[i+255]]=i, got i=%d .+255=%d exp[.]=%d log[.]=%d",
					i, int(i)+255, x2, y2)
			}
		}
		for i := 1; i < 256; i++ {
			x := gf.log[byte(i)]
			y := gf.exp[x]
			if y != byte(i) {
				t.Errorf("expected exp[log[i]]=i, got i=%d log[.]=%d exp[.]=%d",
					i, x, y)
			}
		}
	}

	// Test data taken from:
	//	"The Laws of Cryptography: The Finite Field GF(2^8)"
	//	Neal R. Wagner
	//	http://www.cs.utsa.edu/~wagner/laws/FFM.html
	gf := Poly84310_g3
	expectExp := []byte{
		0x01, 0x03, 0x05, 0x0f, 0x11, 0x33, 0x55, 0xff, 0x1a, 0x2e, 0x72, 0x96, 0xa1, 0xf8, 0x13, 0x35,
		0x5f, 0xe1, 0x38, 0x48, 0xd8, 0x73, 0x95, 0xa4, 0xf7, 0x02, 0x06, 0x0a, 0x1e, 0x22, 0x66, 0xaa,
		0xe5, 0x34, 0x5c, 0xe4, 0x37, 0x59, 0xeb, 0x26, 0x6a, 0xbe, 0xd9, 0x70, 0x90, 0xab, 0xe6, 0x31,
		0x53, 0xf5, 0x04, 0x0c, 0x14, 0x3c, 0x44, 0xcc, 0x4f, 0xd1, 0x68, 0xb8, 0xd3, 0x6e, 0xb2, 0xcd,
		0x4c, 0xd4, 0x67, 0xa9, 0xe0, 0x3b, 0x4d, 0xd7, 0x62, 0xa6, 0xf1, 0x08, 0x18, 0x28, 0x78, 0x88,
		0x83, 0x9e, 0xb9, 0xd0, 0x6b, 0xbd, 0xdc, 0x7f, 0x81, 0x98, 0xb3, 0xce, 0x49, 0xdb, 0x76, 0x9a,
		0xb5, 0xc4, 0x57, 0xf9, 0x10, 0x30, 0x50, 0xf0, 0x0b, 0x1d, 0x27, 0x69, 0xbb, 0xd6, 0x61, 0xa3,
		0xfe, 0x19, 0x2b, 0x7d, 0x87, 0x92, 0xad, 0xec, 0x2f, 0x71, 0x93, 0xae, 0xe9, 0x20, 0x60, 0xa0,
		0xfb, 0x16, 0x3a, 0x4e, 0xd2, 0x6d, 0xb7, 0xc2, 0x5d, 0xe7, 0x32, 0x56, 0xfa, 0x15, 0x3f, 0x41,
		0xc3, 0x5e, 0xe2, 0x3d, 0x47, 0xc9, 0x40, 0xc0, 0x5b, 0xed, 0x2c, 0x74, 0x9c, 0xbf, 0xda, 0x75,
		0x9f, 0xba, 0xd5, 0x64, 0xac, 0xef, 0x2a, 0x7e, 0x82, 0x9d, 0xbc, 0xdf, 0x7a, 0x8e, 0x89, 0x80,
		0x9b, 0xb6, 0xc1, 0x58, 0xe8, 0x23, 0x65, 0xaf, 0xea, 0x25, 0x6f, 0xb1, 0xc8, 0x43, 0xc5, 0x54,
		0xfc, 0x1f, 0x21, 0x63, 0xa5, 0xf4, 0x07, 0x09, 0x1b, 0x2d, 0x77, 0x99, 0xb0, 0xcb, 0x46, 0xca,
		0x45, 0xcf, 0x4a, 0xde, 0x79, 0x8b, 0x86, 0x91, 0xa8, 0xe3, 0x3e, 0x42, 0xc6, 0x51, 0xf3, 0x0e,
		0x12, 0x36, 0x5a, 0xee, 0x29, 0x7b, 0x8d, 0x8c, 0x8f, 0x8a, 0x85, 0x94, 0xa7, 0xf2, 0x0d, 0x17,
		0x39, 0x4b, 0xdd, 0x7c, 0x84, 0x97, 0xa2, 0xfd, 0x1c, 0x24, 0x6c, 0xb4, 0xc7, 0x52, 0xf6, 0x01,
	}
	for i, expect := range expectExp {
		actual := gf.Exp(byte(i))
		if actual != expect {
			t.Errorf("Exp(%#02x): expected %d, got %d", i, expect, actual)
		}
	}
	expectLog := []byte{
		0x00, 0x00, 0x19, 0x01, 0x32, 0x02, 0x1a, 0xc6, 0x4b, 0xc7, 0x1b, 0x68, 0x33, 0xee, 0xdf, 0x03,
		0x64, 0x04, 0xe0, 0x0e, 0x34, 0x8d, 0x81, 0xef, 0x4c, 0x71, 0x08, 0xc8, 0xf8, 0x69, 0x1c, 0xc1,
		0x7d, 0xc2, 0x1d, 0xb5, 0xf9, 0xb9, 0x27, 0x6a, 0x4d, 0xe4, 0xa6, 0x72, 0x9a, 0xc9, 0x09, 0x78,
		0x65, 0x2f, 0x8a, 0x05, 0x21, 0x0f, 0xe1, 0x24, 0x12, 0xf0, 0x82, 0x45, 0x35, 0x93, 0xda, 0x8e,
		0x96, 0x8f, 0xdb, 0xbd, 0x36, 0xd0, 0xce, 0x94, 0x13, 0x5c, 0xd2, 0xf1, 0x40, 0x46, 0x83, 0x38,
		0x66, 0xdd, 0xfd, 0x30, 0xbf, 0x06, 0x8b, 0x62, 0xb3, 0x25, 0xe2, 0x98, 0x22, 0x88, 0x91, 0x10,
		0x7e, 0x6e, 0x48, 0xc3, 0xa3, 0xb6, 0x1e, 0x42, 0x3a, 0x6b, 0x28, 0x54, 0xfa, 0x85, 0x3d, 0xba,
		0x2b, 0x79, 0x0a, 0x15, 0x9b, 0x9f, 0x5e, 0xca, 0x4e, 0xd4, 0xac, 0xe5, 0xf3, 0x73, 0xa7, 0x57,
		0xaf, 0x58, 0xa8, 0x50, 0xf4, 0xea, 0xd6, 0x74, 0x4f, 0xae, 0xe9, 0xd5, 0xe7, 0xe6, 0xad, 0xe8,
		0x2c, 0xd7, 0x75, 0x7a, 0xeb, 0x16, 0x0b, 0xf5, 0x59, 0xcb, 0x5f, 0xb0, 0x9c, 0xa9, 0x51, 0xa0,
		0x7f, 0x0c, 0xf6, 0x6f, 0x17, 0xc4, 0x49, 0xec, 0xd8, 0x43, 0x1f, 0x2d, 0xa4, 0x76, 0x7b, 0xb7,
		0xcc, 0xbb, 0x3e, 0x5a, 0xfb, 0x60, 0xb1, 0x86, 0x3b, 0x52, 0xa1, 0x6c, 0xaa, 0x55, 0x29, 0x9d,
		0x97, 0xb2, 0x87, 0x90, 0x61, 0xbe, 0xdc, 0xfc, 0xbc, 0x95, 0xcf, 0xcd, 0x37, 0x3f, 0x5b, 0xd1,
		0x53, 0x39, 0x84, 0x3c, 0x41, 0xa2, 0x6d, 0x47, 0x14, 0x2a, 0x9e, 0x5d, 0x56, 0xf2, 0xd3, 0xab,
		0x44, 0x11, 0x92, 0xd9, 0x23, 0x20, 0x2e, 0x89, 0xb4, 0x7c, 0xb8, 0x26, 0x77, 0x99, 0xe3, 0xa5,
		0x67, 0x4a, 0xed, 0xde, 0xc5, 0x31, 0xfe, 0x18, 0x0d, 0x63, 0x8c, 0x80, 0xc0, 0xf7, 0x70, 0x07,
	}
	for i, expect := range expectLog {
		if i == 0 {
			continue
		}
		actual := gf.Log(byte(i))
		if actual != expect {
			t.Errorf("Log(%#02x): expected %d, got %d", i, expect, actual)
		}
	}
}

func TestNew_singleton(t *testing.T) {
	a := New(4, 0x7, 2)
	b := New(4, 0x7, 2)
	if a != b {
		t.Errorf("expected singleton, got multiple instances of %#v", a)
	}
}

func TestNew_bad_field_size(t *testing.T) {
	e := panicValue(func() {
		New(17, 0, 0)
	})
	if e != ErrFieldSize {
		t.Errorf("expected panic(ErrFieldSize), got %q", e.Error())
	}
}

func TestNew_out_of_range(t *testing.T) {
	e := panicValue(func() {
		New(16, 15, 0)
	})
	if e != ErrPolyOutOfRange {
		t.Errorf("expected panic(ErrPolyOutOfRange), got %q", e.Error())
	}
	e = panicValue(func() {
		New(16, 32, 0)
	})
	if e != ErrPolyOutOfRange {
		t.Errorf("expected panic(ErrPolyOutOfRange), got %q", e.Error())
	}
}

func TestNew_reducible(t *testing.T) {
	e := panicValue(func() {
		New(64, 0x42, 0x2)
	})
	if e != ErrReduciblePoly {
		t.Errorf("expected panic(ErrReduciblePoly), got %q", e.Error())
	}
}

func TestNew_bad_generator(t *testing.T) {
	e := panicValue(func() {
		New(64, 0x43, 0x1)
	})
	if e != ErrNotGenerator {
		t.Errorf("expected panic(ErrNotGenerator), got %q", e.Error())
	}
	e = panicValue(func() {
		New(64, 0x43, 0x3)
	})
	if e != ErrNotGenerator {
		t.Errorf("expected panic(ErrNotGenerator), got %q", e.Error())
	}
}

func TestGF_String(t *testing.T) {
	type testrow struct {
		field *GF
		gostr string
		str   string
	}
	for idx, row := range []testrow{
		testrow{nil, "nil", "<nil>"},
		testrow{Poly210_g2, "Poly210_g2", "GF(4;b^2+b+1;2)"},
		testrow{Poly310_g2, "Poly310_g2", "GF(8;b^3+b+1;2)"},
		testrow{Poly410_g2, "Poly410_g2", "GF(16;b^4+b+1;2)"},
		testrow{Poly520_g2, "Poly520_g2", "GF(32;b^5+b^2+1;2)"},
		testrow{Poly610_g2, "Poly610_g2", "GF(64;b^6+b+1;2)"},
		testrow{Poly610_g7, "Poly610_g7", "GF(64;b^6+b+1;7)"},
		testrow{Poly710_g2, "Poly710_g2", "GF(128;b^7+b+1;2)"},
		testrow{Poly84310_g3, "Poly84310_g3", "GF(256;b^8+b^4+b^3+b+1;3)"},
		testrow{Poly84320_g2, "Poly84320_g2", "GF(256;b^8+b^4+b^3+b^2+1;2)"},
		testrow{New(16, 0x19, 2), "New(16, 0x19, 2)", "GF(16;b^4+b^3+1;2)"},
	} {
		gostr := row.field.GoString()
		str := row.field.String()
		if gostr != row.gostr {
			t.Errorf("[%2d] expected %q, got %q", idx, row.gostr, gostr)
		}
		if str != row.str {
			t.Errorf("[%2d] expected %q, got %q", idx, row.str, str)
		}
	}
}

var fields = []*GF{
	Poly210_g2,
	Poly310_g2,
	Poly410_g2,
	Poly520_g2,
	Poly610_g2,
	Poly610_g7,
	Poly710_g2,
	Poly84310_g3,
	Poly84320_g2,
}

func TestGF_Add(t *testing.T) {
	var prng = rand.New(rand.NewSource(42))
	for _, field := range fields {
		for i := 0; i < 32; i++ {
			a := byte(prng.Intn(int(field.Size())))
			aa := field.Add(a, a)
			az := field.Add(a, 0)
			za := field.Add(0, a)
			if aa != 0 {
				t.Errorf("[%3d] expected a+a=0, got %d", a, aa)
			}
			if az != a {
				t.Errorf("[%3d] expected a+0=a, got %d", a, az)
			}
			if za != a {
				t.Errorf("[%3d] expected 0+a=a, got %d", a, za)
			}
			for j := 0; j < 16; j++ {
				b := byte(prng.Intn(int(field.Size())))
				ab := field.Add(a, b)
				ba := field.Add(b, a)
				if ab != ba {
					t.Errorf("[%3d,%3d] expected a+b=b+a, got %d != %d",
						a, b, ab, ba)
				}
				for k := 0; k < 16; k++ {
					c := byte(prng.Intn(int(field.Size())))
					type item struct {
						name  string
						value byte
					}
					list := []item{
						item{"a+(b+c)", field.Add(a, field.Add(b, c))},
						item{"a+(c+b)", field.Add(a, field.Add(c, b))},
						item{"b+(a+c)", field.Add(b, field.Add(a, c))},
						item{"b+(c+a)", field.Add(b, field.Add(c, a))},
						item{"c+(a+b)", field.Add(c, field.Add(a, b))},
						item{"c+(b+a)", field.Add(c, field.Add(b, a))},
					}
					for x := range list {
						for y := x + 1; y < len(list); y++ {
							xitem, yitem := list[x], list[y]
							if xitem.value != yitem.value {
								t.Errorf("[%3d,%3d,%3d] expected %s=%s, got %d != %d",
									a, b, c,
									xitem.name, yitem.name,
									xitem.value, yitem.value)
							}
						}
					}
				}
			}
		}
	}
}

func TestGF_Mul(t *testing.T) {
	var prng = rand.New(rand.NewSource(42))
	for _, field := range fields {
		for i := 0; i < 32; i++ {
			a := byte(prng.Intn(int(field.Size())))
			az := field.Mul(a, 0)
			za := field.Mul(0, a)
			ao := field.Mul(a, 1)
			oa := field.Mul(1, a)
			if az != 0 {
				t.Errorf("[%3d] expected a*0=0, got %d", a, az)
			}
			if za != 0 {
				t.Errorf("[%3d] expected 0*a=0, got %d", a, za)
			}
			if ao != a {
				t.Errorf("[%3d] expected a*1=a, got %d", a, ao)
			}
			if oa != a {
				t.Errorf("[%3d] expected 1*a=a, got %d", a, oa)
			}
			for j := 0; j < 16; j++ {
				b := byte(prng.Intn(int(field.Size())))
				ab := field.Mul(a, b)
				ba := field.Mul(b, a)
				if ab != ba {
					t.Errorf("[%3d,%3d] expected a*b=b*a, got %d != %d",
						a, b, ab, ba)
				}
				for k := 0; k < 16; k++ {
					c := byte(prng.Intn(int(field.Size())))
					type item struct {
						name  string
						value byte
					}
					list := []item{
						item{"a*(b*c)", field.Mul(a, field.Mul(b, c))},
						item{"a*(c*b)", field.Mul(a, field.Mul(c, b))},
						item{"b*(a*c)", field.Mul(b, field.Mul(a, c))},
						item{"b*(c*a)", field.Mul(b, field.Mul(c, a))},
						item{"c*(a*b)", field.Mul(c, field.Mul(a, b))},
						item{"c*(b*a)", field.Mul(c, field.Mul(b, a))},
					}
					for x := range list {
						for y := x + 1; y < len(list); y++ {
							xitem, yitem := list[x], list[y]
							if xitem.value != yitem.value {
								t.Errorf("[%3d,%3d,%3d] expected %s=%s, got %d != %d",
									a, b, c,
									xitem.name, yitem.name,
									xitem.value, yitem.value)
							}
						}
					}
				}
			}
		}
	}
}

func TestGF_Div(t *testing.T) {
	var a byte = 0x11
	var b byte = 0x14
	var axb byte = 0x49
	result := Default.Div(axb, b)
	if result != a {
		t.Errorf("expected %d/%d=%d, but got %d", axb, b, a, result)
	}
	result = Default.Div(axb, a)
	if result != b {
		t.Errorf("expected %d/%d=%d, but got %d", axb, a, b, result)
	}
	result = Default.Div(axb, 1)
	if result != axb {
		t.Errorf("expected %d/1=%[1]d, but got %d", axb, result)
	}
	result = Default.Div(0, b)
	if result != 0 {
		t.Errorf("expected 0*%d=0, but got %d", b, result)
	}
	var invb byte = 0xe0
	result = Default.Div(1, b)
	if result != invb {
		t.Errorf("expected 1/%d=%d, but got %d", b, invb, result)
	}
	result = Default.Inv(b)
	if result != invb {
		t.Errorf("expected 1/%d=%d, but got %d", b, invb, result)
	}
	result = Default.Mul(b, invb)
	if result != 1 {
		t.Errorf("expected %d*%d=1, but got %d", b, invb, result)
	}
}

func TestGF_Div_zero(t *testing.T) {
	e := panicValue(func() {
		Default.Div(1, 0)
	})
	if e != ErrDivByZero {
		t.Errorf("expected panic(ErrDivByZero), got %q", e.Error())
	}
	e = panicValue(func() {
		Default.Inv(0)
	})
	if e != ErrDivByZero {
		t.Errorf("expected panic(ErrDivByZero), got %q", e.Error())
	}
}

func TestGF_Log_zero(t *testing.T) {
	e := panicValue(func() {
		Default.Log(0)
	})
	if e != ErrLogZero {
		t.Errorf("expected panic(ErrLogZero), got %q", e.Error())
	}
}

func TestGF_Exp(t *testing.T) {
	var ggg byte = 8
	result := Default.Exp(3)
	if result != ggg {
		t.Errorf("expected 3^3=%d, got %d", ggg, result)
	}
	result = Default.Log(ggg)
	if result != 3 {
		t.Errorf("expected log_3(%d)=3, got %d", ggg, result)
	}
}

func TestGF_Compare(t *testing.T) {
	type testrow struct {
		left, right *GF
		cmp         int
	}
	for _, row := range []testrow{
		testrow{Default, Default, 0},
		testrow{Poly84310_g3, Poly84320_g2, -1},
		testrow{Poly84320_g2, Poly84310_g3, 1},
		testrow{Poly84310_g3, Poly610_g7, 1},
		testrow{Poly610_g7, Poly84310_g3, -1},
		testrow{Poly610_g2, Poly610_g7, -1},
		testrow{Poly610_g7, Poly610_g2, 1},
	} {
		cmp := row.left.Compare(row.right)
		if cmp != row.cmp {
			t.Errorf("expected %#v.Compare(%#v) == %d, got %d", row.left, row.right, row.cmp, cmp)
		}
		cmp2 := -row.right.Compare(row.left)
		if cmp2 != cmp {
			t.Errorf("expected a.Compare(b) %d == -b.Compare(a) %d for %#v and %#v", cmp, cmp2, row.left, row.right)
		}
		if cmp == 0 {
			if !row.left.Equal(row.right) {
				t.Errorf("expected [cmp=0] => [a.Equal(b)] for %#v and %#v", row.left, row.right)
			}
			if !row.right.Equal(row.left) {
				t.Errorf("expected [cmp=0] => [a.Equal(b)] for %#v and %#v", row.right, row.left)
			}
		} else if cmp < 0 {
			if !row.left.Less(row.right) {
				t.Errorf("expected [cmp<0] => [a.Less(b)] for %#v and %#v", row.left, row.right)
			}
		} else {
			if !row.right.Less(row.left) {
				t.Errorf("expected [cmp>0] => [b.Less(a)] for %#v and %#v", row.left, row.right)
			}
		}
	}
}

func TestGF_Polynomial(t *testing.T) {
	p := Poly84310_g3.Polynomial()
	if p != 0x11b {
		t.Errorf("Poly84310_g3.Polynomial(): expected 0x11b, got %#x", p)
	}
	g := Poly84310_g3.Generator()
	if g != 0x3 {
		t.Errorf("Poly84310_g3.Generator(): expected 0x3, got %#x", g)
	}
}

func TestGF8(t *testing.T) {
	// Test data taken from:
	//	http://math.stackexchange.com/questions/245621/arithmetic-operations-in-galois-field
	// Not sure what the original source was.
	addmatrix := [][]byte{
		[]byte{0, 1, 2, 3, 4, 5, 6, 7},
		[]byte{1, 0, 3, 2, 5, 4, 7, 6},
		[]byte{2, 3, 0, 1, 6, 7, 4, 5},
		[]byte{3, 2, 1, 0, 7, 6, 5, 4},
		[]byte{4, 5, 6, 7, 0, 1, 2, 3},
		[]byte{5, 4, 7, 6, 1, 0, 3, 2},
		[]byte{6, 7, 4, 5, 2, 3, 0, 1},
		[]byte{7, 6, 5, 4, 3, 2, 1, 0},
	}
	for i := byte(0); i < 8; i++ {
		for j := byte(0); j < 8; j++ {
			actual := Poly310_g2.Add(i, j)
			expect := addmatrix[i][j]
			if actual != expect {
				t.Errorf("expected %v+%v=%v, got %v", i, j, expect, actual)
			}
		}
	}
	mulmatrix := [][]byte{
		[]byte{0, 0, 0, 0, 0, 0, 0, 0},
		[]byte{0, 1, 2, 3, 4, 5, 6, 7},
		[]byte{0, 2, 4, 6, 3, 1, 7, 5},
		[]byte{0, 3, 6, 5, 7, 4, 1, 2},
		[]byte{0, 4, 3, 7, 6, 2, 5, 1},
		[]byte{0, 5, 1, 4, 2, 7, 3, 6},
		[]byte{0, 6, 7, 1, 5, 3, 2, 4},
		[]byte{0, 7, 5, 2, 1, 6, 4, 3},
	}
	for i := byte(0); i < 8; i++ {
		for j := byte(0); j < 8; j++ {
			actual := Poly310_g2.Mul(i, j)
			expect := mulmatrix[i][j]
			if actual != expect {
				t.Errorf("expected %v*%v=%v, got %v", i, j, expect, actual)
			}
		}
	}
}

func BenchmarkGF_Add_256(b *testing.B) {
	gf := Default
	var x byte = 17
	var y byte = 48
	for i := 0; i < b.N; i++ {
		_ = gf.Add(x, y)
	}
}

func BenchmarkGF_Mul_256(b *testing.B) {
	gf := Default
	var x byte = 1
	var y byte = 3
	for i := 0; i < b.N; i++ {
		_ = gf.Mul(x, y)
	}
}

func BenchmarkGF_Div_256(b *testing.B) {
	gf := Default
	var x byte = 1
	var y byte = 3
	for i := 0; i < b.N; i++ {
		_ = gf.Div(x, y)
	}
}

func BenchmarkGF_Exp_256(b *testing.B) {
	gf := Default
	var x byte = 1
	for i := 0; i < b.N; i++ {
		_ = gf.Exp(x)
	}
}

func BenchmarkGF_Log_256(b *testing.B) {
	gf := Default
	var x byte = 1
	for i := 0; i < b.N; i++ {
		_ = gf.Log(x)
	}
}

func panicValue(f func()) (value error) {
	defer func() {
		if e, ok := recover().(error); ok {
			value = e
		}
	}()
	f()
	return
}
