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
