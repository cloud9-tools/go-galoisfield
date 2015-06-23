# go-galoisfield
An implementation of `GF(2**m)` for Go

## What is a Galois field?

A Galois field, also known as a finite field, is a mathematical field with a
number of elements equal to a prime number to a positive integer power.  While
finite fields with a prime number of elements are familiar to most
programmers -- boolean arithmetic is an example -- taking that prime to powers
higher than 1 is less well-known.  Basically, an element of `GF(2**m)` can be
seen as a list of m bits, where addition is elementwise using boolean 
arithmetic (`a+b` is `a XOR b`), and the remaining rules of field arithmetic 
follow from linear algebra (vectors, or alternatively, polynomial coefficients).

* http://en.wikipedia.org/wiki/Finite_field
* http://www.cs.utsa.edu/~wagner/laws/FFM.html
* http://research.swtch.com/field

Short version: an element of `GF(2**8)` element may be represented as a byte
(0 ≤ n ≤ 255), but it's really a vector of 8 bits -- like a very primitive
MMX/SSE.  We then treat said vector as the coefficients of a polynomial, and
that allows us to define multiplication, giving us a full mathematical field.

Finite fields -- and `GF(2**8)` in particular -- get a tons of use in codes,
in both the "error-correcting code" and "cryptographic code" senses.
However, this implementation has NOT been hardened against timing attacks,
so it MUST NOT be used in cryptography.
