// Package cflags contains condition flags information, it is
// used for providing information about the most recently
// executed condition.
package cflags

const (
	// FLPOS is the positive flag which indicates a positive sign.
	FLPOS = 1 << 0

	// FLZRO is the zero flag which indicates a zero sign.
	FLZRO = 1 << 1

	// FLNEG is the negative flag which indicates a negative sign.
	FLNEG = 1 << 2
)
