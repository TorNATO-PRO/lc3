// Package opcodes contains opcodes that will be used by the
// LC3 processor. An opcode specifies the kind of task to perform
// over a set of parameters.
package opcodes

const (
	// OPBR specifies the "branch" opcode.
	OPBR = iota

	// OPADD specifies the "add" opcode.
	OPADD

	// OPLD specifies the "load" opcode.
	OPLD

	// OPST specifies the "store" opcode.
	OPST

	// OPJSR specifies the "jump" opcode.
	OPJSR

	// OPAND specifies the "bitwise and" opcode.
	OPAND

	// OPLDR specifies the "load" opcode.
	OPLDR

	// OPSTR specifies the "store" opcode.
	OPSTR

	// OPRTI specifies the "unused" opcode.
	OPRTI

	// OPNOT specifies the "bitwise not" opcode.
	OPNOT

	// OPLDI specifies the "load indirect" opcode.
	OPLDI

	// OPSTI specifies the "store indirect" opcode.
	OPSTI

	// OPJMP specifies the "jump" opcode.
	OPJMP

	// OPRES specifies the "reserved" opcode.
	OPRES

	// OPLEA specifies the "load effective address" opcode.
	OPLEA

	// OPTRAP specifies the "executes trap" opcode.
	OPTRAP
)
